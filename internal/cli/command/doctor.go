package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/aruodore/aruo/internal/cli/iostreams"
	"github.com/aruodore/aruo/internal/clierror"
	"github.com/aruodore/aruo/internal/doctor"
	"github.com/spf13/cobra"
)

type doctorOptions struct {
	format       string
	minimumScore int
}

func newDoctor(streams iostreams.IOStreams, service *doctor.Service) *cobra.Command {
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
			report, err := service.Audit(command.Context(), target)
			if err != nil {
				return err
			}
			switch options.format {
			case "human":
				err = renderDoctorHuman(streams, report)
			case "json":
				encoder := json.NewEncoder(streams.Out)
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
	return command
}

func renderDoctorHuman(streams iostreams.IOStreams, report doctor.Report) error {
	if _, err := fmt.Fprintf(streams.Out, "Repository health: %d/%d (%s)\n%s\n\n", report.Score, report.MaxScore, report.Grade, report.Repository); err != nil {
		return err
	}
	for _, category := range report.Categories {
		if _, err := fmt.Fprintf(streams.Out, "  %-15s %2d/%2d\n", category.Category, category.Points, category.MaxPoints); err != nil {
			return err
		}
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
		_, err := fmt.Fprintln(streams.Out, "\nNo recommendations. All local v1 checks passed.")
		return err
	}
	if _, err := fmt.Fprintln(streams.Out, "\nRecommendations:"); err != nil {
		return err
	}
	for _, item := range deductions {
		if _, err := fmt.Fprintf(streams.Out, "  - %s\n    %s\n", item.recommendation.Message, item.recommendation.Action); err != nil {
			return err
		}
	}
	return nil
}
