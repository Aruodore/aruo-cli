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
		return presenter.Message(ctx, tux.Message{Kind: tux.MessageSuccess, Text: "No recommendations. All local v1 checks passed."})
	}
	lines := make([]string, 0, len(deductions)*2+1)
	lines = append(lines, "Recommendations:")
	for _, item := range deductions {
		lines = append(lines, "  - "+item.recommendation.Message, "    "+item.recommendation.Action)
	}
	return presenter.Message(ctx, tux.Message{Kind: tux.MessageInfo, Text: strings.Join(lines, "\n")})
}
