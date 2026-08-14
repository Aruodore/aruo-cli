package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/aruodore/aruo-cli/internal/clierror"
	"github.com/aruodore/aruo-cli/internal/doctor"
	"github.com/aruodore/aruo-cli/internal/tux"
	"github.com/spf13/cobra"
)

type doctorOptions struct {
	format       string
	minimumScore int
}

func newDoctor(factory sessionFactory, service *doctor.Service) *cobra.Command {
	options := doctorOptions{format: "human", minimumScore: 80}
	command := &cobra.Command{
		Use:   "doctor [repository]",
		Short: "Assess repository engineering health",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			target := "."
			if len(args) == 1 {
				target = args[0]
			}
			if options.minimumScore < 0 || options.minimumScore > 100 {
				return errors.New("--minimum-score must be between 0 and 100")
			}
			ctx := command.Context()
			report, err := service.Audit(ctx, target)
			if err != nil {
				return err
			}
			switch options.format {
			case "human":
				terminal, sessionErr := factory.build(ctx, false)
				if sessionErr != nil {
					return sessionErr
				}
				defer func() { _ = terminal.Close() }()
				err = renderDoctorHuman(ctx, terminal.Presenter(), report)
			case "json":
				encoder := json.NewEncoder(factory.streams.Out)
				encoder.SetIndent("", "  ")
				err = encoder.Encode(report)
			default:
				return fmt.Errorf("unsupported format %q; use human or json", options.format)
			}
			if err != nil {
				return err
			}
			if report.Contract.BlockingFindings > 0 {
				return &clierror.Error{Code: 3, Err: fmt.Errorf("repository has %d blocking contract finding(s)", report.Contract.BlockingFindings), Silent: true}
			}
			if report.Intent.BlockingFindings > 0 {
				return &clierror.Error{Code: 3, Err: fmt.Errorf("repository has %d blocking intent finding(s)", report.Intent.BlockingFindings), Silent: true}
			}
			if report.Score < options.minimumScore {
				return &clierror.Error{Code: 3, Err: fmt.Errorf("repository score %d is below minimum %d", report.Score, options.minimumScore), Silent: true}
			}
			return nil
		},
	}
	command.Flags().StringVar(&options.format, "format", "human", "output format: human or json")
	command.Flags().IntVar(&options.minimumScore, "minimum-score", 80, "minimum passing score (0-100)")
	_ = command.RegisterFlagCompletionFunc("format", func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return []string{"human", "json"}, cobra.ShellCompDirectiveNoFileComp
	})
	return command
}

func renderDoctorHuman(ctx context.Context, presenter tux.Presenter, report doctor.Report) error {
	summary := tux.Table{
		Columns: []tux.Column{
			{ID: "category", Heading: "Category"},
			{ID: "score", Heading: "Score", Alignment: tux.AlignRight},
		},
	}
	for _, category := range report.Categories {
		summary.Rows = append(summary.Rows, []string{string(category.Category), fmt.Sprintf("%d/%d", category.Points, category.MaxPoints)})
	}
	header := fmt.Sprintf("Repository health: %d/%d (%s)\n%s", report.Score, report.MaxScore, report.Grade, report.Repository)
	if err := presenter.Message(ctx, tux.Message{Kind: tux.MessageInfo, Text: header}); err != nil {
		return err
	}
	if err := presenter.Table(ctx, summary); err != nil {
		return err
	}

	type deduction struct {
		lost           int
		id             string
		recommendation doctor.Recommendation
	}
	var deductions []deduction
	for _, assessment := range report.Assessments {
		for _, recommendation := range assessment.Recommendations {
			deductions = append(deductions, deduction{assessment.MaxPoints - assessment.Points, assessment.ID, recommendation})
		}
	}
	sort.SliceStable(deductions, func(i, j int) bool {
		if deductions[i].lost != deductions[j].lost {
			return deductions[i].lost > deductions[j].lost
		}
		return deductions[i].id < deductions[j].id
	})
	if len(deductions) == 0 {
		if err := presenter.Message(ctx, tux.Message{Kind: tux.MessageSuccess, Text: "No recommendations. All local v1 checks passed."}); err != nil {
			return err
		}
	} else {
		lines := make([]string, 0, len(deductions)*2+1)
		lines = append(lines, "Recommendations:")
		for _, item := range deductions {
			lines = append(lines, "  - "+item.recommendation.Message, "    "+item.recommendation.Action)
		}
		if err := presenter.Message(ctx, tux.Message{Kind: tux.MessageInfo, Text: strings.Join(lines, "\n")}); err != nil {
			return err
		}
	}
	if err := renderDoctorContract(ctx, presenter, report.Contract); err != nil {
		return err
	}
	return renderDoctorIntent(ctx, presenter, report.Intent)
}

func renderDoctorContract(ctx context.Context, presenter tux.Presenter, report doctor.ContractReport) error {
	if !report.Present {
		return presenter.Message(ctx, tux.Message{Kind: tux.MessageInfo, Text: "AI engineering contract: not installed (run aruo init)."})
	}
	kind := tux.MessageSuccess
	if !report.Valid {
		kind = tux.MessageInfo
	}
	if err := presenter.Message(ctx, tux.Message{Kind: kind, Text: fmt.Sprintf("AI engineering contract: version %s, %d managed files, %d blocking findings", report.Version, len(report.Files), report.BlockingFindings)}); err != nil {
		return err
	}
	if len(report.Findings) == 0 {
		return nil
	}
	lines := []string{"Contract findings:"}
	for _, finding := range report.Findings {
		lines = append(lines, fmt.Sprintf("  - [%s] %s", finding.Severity, finding.Message), "    "+finding.Action)
	}
	return presenter.Message(ctx, tux.Message{Kind: tux.MessageInfo, Text: strings.Join(lines, "\n")})
}

func renderDoctorIntent(ctx context.Context, presenter tux.Presenter, report doctor.IntentReport) error {
	if !report.Present {
		return presenter.Message(ctx, tux.Message{Kind: tux.MessageInfo, Text: "Production intent: not declared (no aruo.yaml)."})
	}
	table := tux.Table{Columns: []tux.Column{
		{ID: "capability", Heading: "Capability"},
		{ID: "status", Heading: "Intent"},
		{ID: "evidence", Heading: "Evidence"},
	}}
	for _, capability := range report.Capabilities {
		evidence := string(capability.EvidenceStatus)
		if capability.Evidence != "" {
			evidence += ": " + capability.Evidence
		}
		table.Rows = append(table.Rows, []string{capability.Name, string(capability.Status), evidence})
	}
	header := fmt.Sprintf("Production intent: %d capabilities, %d blocking findings", len(report.Capabilities), report.BlockingFindings)
	kind := tux.MessageInfo
	if report.BlockingFindings == 0 && report.Valid {
		kind = tux.MessageSuccess
	}
	if err := presenter.Message(ctx, tux.Message{Kind: kind, Text: header}); err != nil {
		return err
	}
	if len(table.Rows) > 0 {
		if err := presenter.Table(ctx, table); err != nil {
			return err
		}
	}
	if len(report.Findings) == 0 {
		return nil
	}
	lines := []string{"Intent findings:"}
	for _, finding := range report.Findings {
		label := finding.Severity
		if finding.Capability != "" {
			label += " " + finding.Capability
		}
		lines = append(lines, fmt.Sprintf("  - [%s] %s", label, finding.Message), "    "+finding.Action)
	}
	return presenter.Message(ctx, tux.Message{Kind: tux.MessageInfo, Text: strings.Join(lines, "\n")})
}
