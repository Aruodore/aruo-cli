// Package doctor evaluates read-only repository health evidence.
package doctor

// Category identifies one stable scoring dimension.
type Category string

const (
	CategoryCompleteness  Category = "completeness"
	CategoryDocumentation Category = "documentation"
	CategoryCI            Category = "ci"
	CategoryTests         Category = "tests"
	CategoryLicense       Category = "license"
	CategorySecurity      Category = "security"
	CategoryGitHub        Category = "github"
)

// Recommendation describes one concrete improvement.
type Recommendation struct {
	Message string `json:"message"`
	Action  string `json:"action"`
}

// Assessment is the evidence and score returned by one check.
type Assessment struct {
	ID              string           `json:"id"`
	Category        Category         `json:"category"`
	Title           string           `json:"title"`
	Points          int              `json:"points"`
	MaxPoints       int              `json:"maxPoints"`
	Evidence        []string         `json:"evidence"`
	Recommendations []Recommendation `json:"recommendations,omitempty"`
}

// CategoryScore is an aggregate without hiding its denominator.
type CategoryScore struct {
	Category  Category `json:"category"`
	Points    int      `json:"points"`
	MaxPoints int      `json:"maxPoints"`
}

// Report is the versioned repository health result.
type Report struct {
	SchemaVersion string          `json:"schemaVersion"`
	Policy        string          `json:"policy"`
	Repository    string          `json:"repository"`
	Score         int             `json:"score"`
	MaxScore      int             `json:"maxScore"`
	Grade         string          `json:"grade"`
	Categories    []CategoryScore `json:"categories"`
	Assessments   []Assessment    `json:"assessments"`
}
