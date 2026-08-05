package charm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/aruodore/aruo/internal/tux"
)

// Progress renders task events in one Bubble Tea-owned terminal region.
type Progress struct {
	program   *tea.Program
	done      chan error
	closeOnce sync.Once
	closed    chan struct{}
	closeErr  error
}

// NewProgress starts a live progress renderer with Aruo-owned lifecycle and signals.
func NewProgress(ctx context.Context, output io.Writer, capabilities tux.Capabilities) *Progress {
	model := newProgressModel(capabilities)
	program := tea.NewProgram(
		model,
		tea.WithContext(ctx),
		tea.WithInput(nil),
		tea.WithOutput(output),
		tea.WithFPS(10),
		tea.WithWindowSize(capabilities.Width, capabilities.Height),
		tea.WithoutSignalHandler(),
	)
	progress := &Progress{program: program, done: make(chan error, 1), closed: make(chan struct{})}
	go func() {
		_, err := program.Run()
		progress.done <- err
	}()
	return progress
}

// Emit sends one task transition to the live renderer.
func (p *Progress) Emit(ctx context.Context, event tux.TaskEvent) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case <-p.closed:
		return tux.ErrUnavailable
	default:
		p.program.Send(taskEventMsg(event))
		return nil
	}
}

// Close completes the live region and restores Bubble Tea terminal state.
func (p *Progress) Close() error {
	p.closeOnce.Do(func() {
		close(p.closed)
		p.program.Quit()
		p.closeErr = <-p.done
		if errors.Is(p.closeErr, tea.ErrProgramKilled) {
			p.closeErr = nil
		}
	})
	return p.closeErr
}

type taskEventMsg tux.TaskEvent

type taskState struct {
	id       string
	parentID string
	kind     tux.TaskEventKind
	label    string
	current  int64
	total    int64
}

type progressModel struct {
	spinner spinner.Model
	tasks   map[string]taskState
	order   []string
	unicode bool
}

func newProgressModel(capabilities tux.Capabilities) progressModel {
	indicator := spinner.Line
	if capabilities.Unicode {
		indicator = spinner.MiniDot
	}
	return progressModel{
		spinner: spinner.New(spinner.WithSpinner(indicator)),
		tasks:   make(map[string]taskState),
		unicode: capabilities.Unicode,
	}
}

func (m progressModel) Init() tea.Cmd {
	return nil
}

func (m progressModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case taskEventMsg:
		event := tux.TaskEvent(message)
		wasRunning := m.hasRunning()
		state, exists := m.tasks[event.TaskID]
		if !exists {
			m.order = append(m.order, event.TaskID)
			state.id = event.TaskID
		}
		state.parentID = event.ParentID
		state.kind = event.Kind
		state.label = event.Label
		state.current = event.Current
		state.total = event.Total
		m.tasks[event.TaskID] = state
		if !wasRunning && m.hasRunning() {
			return m, m.spinner.Tick
		}
	case spinner.TickMsg:
		if !m.hasRunning() {
			return m, nil
		}
		var command tea.Cmd
		m.spinner, command = m.spinner.Update(message)
		return m, command
	}
	return m, nil
}

func (m progressModel) View() tea.View {
	lines := make([]string, 0, len(m.order))
	for _, id := range m.order {
		task := m.tasks[id]
		depth := m.depth(task)
		lines = append(lines, strings.Repeat("  ", depth)+m.taskLine(task))
	}
	return tea.NewView(strings.Join(lines, "\n"))
}

func (m progressModel) hasRunning() bool {
	for _, task := range m.tasks {
		switch task.kind {
		case tux.TaskStarted, tux.TaskAdvanced, tux.TaskMessage, tux.TaskRetrying:
			return true
		default:
		}
	}
	return false
}

func (m progressModel) depth(task taskState) int {
	depth := 0
	seen := map[string]struct{}{task.id: {}}
	for parent := task.parentID; parent != ""; {
		if _, cycle := seen[parent]; cycle {
			break
		}
		seen[parent] = struct{}{}
		depth++
		ancestor, found := m.tasks[parent]
		if !found {
			break
		}
		parent = ancestor.parentID
	}
	return depth
}

func (m progressModel) taskLine(task taskState) string {
	indicator := m.spinner.View()
	switch task.kind {
	case tux.TaskCompleted:
		indicator = "done"
		if m.unicode {
			indicator = "✓"
		}
	case tux.TaskFailed:
		indicator = "failed"
		if m.unicode {
			indicator = "✗"
		}
	case tux.TaskCancelled:
		indicator = "cancelled"
	case tux.TaskRetrying:
		indicator = "retry"
	default:
	}
	line := indicator + " " + task.label
	if task.total > 0 {
		line += " " + progressBar(task.current, task.total, m.unicode)
	}
	return line
}

func progressBar(current, total int64, unicode bool) string {
	const width = 10
	current = max(0, min(current, total))
	filled := 0
	for filled < width && int64(filled+1)*total <= current*int64(width) {
		filled++
	}
	full, empty := "#", "-"
	if unicode {
		full, empty = "━", "─"
	}
	percentage := current * 100 / total
	return fmt.Sprintf("[%s%s] %d%%", strings.Repeat(full, filled), strings.Repeat(empty, width-filled), percentage)
}
