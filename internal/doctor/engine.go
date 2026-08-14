package doctor

import (
	"context"
	"errors"
	"fmt"
	"sort"
)

const PolicyV1 = "aruo.repository-health/v1"

// Check is the stable in-process contract mirrored by future plugin messages.
type Check interface {
	ID() string
	Category() Category
	MaxPoints() int
	Run(context.Context, Repository) (Assessment, error)
}

// Engine validates and aggregates registered checks.
type Engine struct{ checks []Check }

// NewEngine creates an engine and rejects ambiguous scoring contracts.
func NewEngine(checks ...Check) (*Engine, error) {
	seen := make(map[string]struct{}, len(checks))
	for _, check := range checks {
		if check == nil || check.ID() == "" || check.MaxPoints() <= 0 || !validCategory(check.Category()) {
			return nil, errors.New("invalid doctor check registration")
		}
		if _, exists := seen[check.ID()]; exists {
			return nil, fmt.Errorf("duplicate doctor check %q", check.ID())
		}
		seen[check.ID()] = struct{}{}
	}
	return &Engine{checks: append([]Check(nil), checks...)}, nil
}

// Run executes checks and returns a deterministic report.
func (e *Engine) Run(ctx context.Context, repositoryName string, repository Repository) (Report, error) {
	assessments := make([]Assessment, 0, len(e.checks))
	for _, check := range e.checks {
		if err := ctx.Err(); err != nil {
			return Report{}, err
		}
		assessment, err := check.Run(ctx, repository)
		if err != nil {
			return Report{}, fmt.Errorf("doctor check %s: %w", check.ID(), err)
		}
		if assessment.ID != check.ID() || assessment.Category != check.Category() ||
			assessment.MaxPoints != check.MaxPoints() || assessment.Points < 0 || assessment.Points > assessment.MaxPoints {
			return Report{}, fmt.Errorf("doctor check %s returned an invalid assessment", check.ID())
		}
		assessments = append(assessments, assessment)
	}
	sort.Slice(assessments, func(i, j int) bool { return assessments[i].ID < assessments[j].ID })

	categoryMap := map[Category]*CategoryScore{}
	report := Report{SchemaVersion: "3", Policy: PolicyV1, Repository: repositoryName, Assessments: assessments}
	for _, assessment := range assessments {
		report.Score += assessment.Points
		report.MaxScore += assessment.MaxPoints
		category := categoryMap[assessment.Category]
		if category == nil {
			category = &CategoryScore{Category: assessment.Category}
			categoryMap[assessment.Category] = category
		}
		category.Points += assessment.Points
		category.MaxPoints += assessment.MaxPoints
	}
	for _, category := range categoryOrder() {
		if score := categoryMap[category]; score != nil {
			report.Categories = append(report.Categories, *score)
		}
	}
	report.Grade = grade(report.Score, report.MaxScore)
	return report, nil
}

func grade(points, maximum int) string {
	if maximum == 0 {
		return "N/A"
	}
	percentage := points * 100 / maximum
	switch {
	case percentage >= 90:
		return "A"
	case percentage >= 80:
		return "B"
	case percentage >= 70:
		return "C"
	case percentage >= 60:
		return "D"
	default:
		return "F"
	}
}

func validCategory(category Category) bool {
	for _, candidate := range categoryOrder() {
		if category == candidate {
			return true
		}
	}
	return false
}

func categoryOrder() []Category {
	return []Category{CategoryCompleteness, CategoryDocumentation, CategoryCI, CategoryTests, CategoryLicense, CategorySecurity, CategoryGitHub}
}
