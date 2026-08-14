// Package charm adapts the Charm v2 rendering stack to Aruo terminal ports.
package charm

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/aruodore/aruo-cli/internal/tux"
)

// Presenter renders capability-aware human output with semantic styles.
type Presenter struct {
	result       io.Writer
	diagnostic   io.Writer
	capabilities tux.Capabilities
	policy       tux.Policy
	theme        theme
}

// NewPresenter constructs a styled presenter over explicit result and diagnostic streams.
func NewPresenter(result, diagnostic io.Writer, capabilities tux.Capabilities, policy tux.Policy) *Presenter {
	return &Presenter{
		result:       result,
		diagnostic:   diagnostic,
		capabilities: capabilities,
		policy:       policy,
		theme:        newTheme(capabilities, policy),
	}
}

// Message writes one durable semantic message.
func (p *Presenter) Message(ctx context.Context, message tux.Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	label, marker, style := p.messageStyle(message.Kind)
	prefix := ""
	if marker != "" {
		prefix = style.Render(marker) + " "
	} else if label != "" {
		prefix = style.Render(label+":") + " "
	}
	if _, err := fmt.Fprintln(p.result, prefix+message.Text); err != nil {
		return fmt.Errorf("write message: %w", err)
	}
	return nil
}

// Table writes a width-aware table and omits lower-priority columns when necessary.
func (p *Presenter) Table(ctx context.Context, table tux.Table) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	columns := visibleColumns(table, p.capabilities.Width)
	if len(columns) == 0 {
		return nil
	}
	widths := columnWidths(table, columns)
	if err := p.writeRow(table, columns, widths, tableHeadings(table, columns), true); err != nil {
		return err
	}
	for _, row := range table.Rows {
		if err := p.writeRow(table, columns, widths, row, false); err != nil {
			return err
		}
	}
	return nil
}

// Diagnostic writes one actionable error without exposing renderer details.
func (p *Presenter) Diagnostic(ctx context.Context, diagnostic tux.Diagnostic) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(p.diagnostic, p.theme.danger.Render("Error: "+diagnostic.Summary)); err != nil {
		return fmt.Errorf("write diagnostic summary: %w", err)
	}
	for _, item := range []struct{ label, value string }{
		{"Detail", diagnostic.Detail},
		{"Effect", diagnostic.Effect},
		{"Next", diagnostic.Suggestion},
	} {
		if item.value == "" {
			continue
		}
		if _, err := fmt.Fprintf(p.diagnostic, "  %s  %s\n", p.theme.muted.Render(item.label), item.value); err != nil {
			return fmt.Errorf("write diagnostic detail: %w", err)
		}
	}
	if diagnostic.Code != "" {
		if _, err := fmt.Fprintf(p.diagnostic, "  %s  %s\n", p.theme.muted.Render("Reference"), diagnostic.Code); err != nil {
			return fmt.Errorf("write diagnostic reference: %w", err)
		}
	}
	return nil
}

func (p *Presenter) messageStyle(kind tux.MessageKind) (string, string, lipgloss.Style) {
	unicode := p.capabilities.Unicode && p.policy.Unicode != tux.FeatureNever && !p.policy.Accessible
	switch kind {
	case tux.MessageSuccess:
		if unicode {
			return "", "✓", p.theme.success
		}
		return "Success", "", p.theme.success
	case tux.MessageWarning:
		if unicode {
			return "", "!", p.theme.warning
		}
		return "Warning", "", p.theme.warning
	case tux.MessageNote:
		return "Note", "", p.theme.muted
	default:
		return "", "", p.theme.emphasis
	}
}

func (p *Presenter) writeRow(table tux.Table, columns, widths []int, row []string, heading bool) error {
	cells := make([]string, 0, len(columns))
	for index, source := range columns {
		value := ""
		if source < len(row) {
			value = row[source]
		}
		alignment := table.Columns[source].Alignment
		cell := pad(value, widths[index], alignment)
		if heading {
			cell = p.theme.accent.Render(strings.ToUpper(cell))
		}
		cells = append(cells, cell)
	}
	if _, err := fmt.Fprintln(p.result, strings.Join(cells, "  ")); err != nil {
		return fmt.Errorf("write table row: %w", err)
	}
	return nil
}

func visibleColumns(table tux.Table, width int) []int {
	columns := make([]int, 0, len(table.Columns))
	for index, column := range table.Columns {
		if !column.Sensitive {
			columns = append(columns, index)
		}
	}
	if width <= 0 {
		width = 80
	}
	for len(columns) > 1 && requiredWidth(columnWidths(table, columns)) > width {
		remove := 0
		for index := 1; index < len(columns); index++ {
			candidate := table.Columns[columns[index]]
			current := table.Columns[columns[remove]]
			if candidate.Priority > current.Priority || candidate.Priority == current.Priority && columns[index] > columns[remove] {
				remove = index
			}
		}
		columns = append(columns[:remove], columns[remove+1:]...)
	}
	sort.Ints(columns)
	return columns
}

func columnWidths(table tux.Table, columns []int) []int {
	widths := make([]int, len(columns))
	for index, source := range columns {
		widths[index] = lipgloss.Width(table.Columns[source].Heading)
		for _, row := range table.Rows {
			if source < len(row) {
				widths[index] = max(widths[index], lipgloss.Width(row[source]))
			}
		}
	}
	return widths
}

func requiredWidth(widths []int) int {
	width := 0
	for _, column := range widths {
		width += column
	}
	if len(widths) > 1 {
		width += (len(widths) - 1) * 2
	}
	return width
}

func tableHeadings(table tux.Table, columns []int) []string {
	row := make([]string, len(table.Columns))
	for _, index := range columns {
		row[index] = table.Columns[index].Heading
	}
	return row
}

func pad(value string, width int, alignment tux.Alignment) string {
	missing := max(0, width-lipgloss.Width(value))
	if alignment == tux.AlignRight {
		return strings.Repeat(" ", missing) + value
	}
	return value + strings.Repeat(" ", missing)
}
