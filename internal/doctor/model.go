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
	Intent        IntentReport    `json:"intent"`
}

// CapabilityStatus is the repository's explicit production-intent vocabulary.
type CapabilityStatus string

const (
	CapabilitySolved   CapabilityStatus = "SOLVED"
	CapabilityRequired CapabilityStatus = "REQUIRED"
	CapabilityOptional CapabilityStatus = "OPTIONAL"
	CapabilityDeferred CapabilityStatus = "DEFERRED"
	CapabilityUnknown  CapabilityStatus = "UNKNOWN"
)

// EvidenceStatus describes what Doctor can prove from the local repository.
type EvidenceStatus string

const (
	EvidenceVerified      EvidenceStatus = "VERIFIED"
	EvidenceDeclared      EvidenceStatus = "DECLARED"
	EvidenceMissing       EvidenceStatus = "MISSING"
	EvidenceNotApplicable EvidenceStatus = "NOT_APPLICABLE"
)

// IntentCapability is one normalized capability declaration.
type IntentCapability struct {
	Name           string           `json:"name"`
	Status         CapabilityStatus `json:"status"`
	Evidence       string           `json:"evidence,omitempty"`
	Reason         string           `json:"reason,omitempty"`
	EvidenceStatus EvidenceStatus   `json:"evidenceStatus"`
}

// IntentFinding reports a malformed claim or visible unresolved responsibility.
type IntentFinding struct {
	Capability string `json:"capability,omitempty"`
	Severity   string `json:"severity"`
	Message    string `json:"message"`
	Action     string `json:"action"`
}

// IntentReport audits aruo.yaml without changing the repository-health/v1 score.
type IntentReport struct {
	Present          bool               `json:"present"`
	Valid            bool               `json:"valid"`
	APIVersion       string             `json:"apiVersion,omitempty"`
	TemplateID       string             `json:"templateId,omitempty"`
	Profile          string             `json:"profile,omitempty"`
	Capabilities     []IntentCapability `json:"capabilities,omitempty"`
	Findings         []IntentFinding    `json:"findings,omitempty"`
	BlockingFindings int                `json:"blockingFindings"`
}
