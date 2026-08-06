package doctor

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var pinnedAction = regexp.MustCompile(`(?m)^\s*-?\s*uses:\s*[^\s@]+@[0-9a-fA-F]{40}(?:\s|$)`)
var anyAction = regexp.MustCompile(`(?m)^\s*-?\s*uses:\s*[^\s]+`)

type categoryCheck struct {
	id       string
	category Category
	maximum  int
	run      func(Repository) (Assessment, error)
}

func (c categoryCheck) ID() string         { return c.id }
func (c categoryCheck) Category() Category { return c.category }
func (c categoryCheck) MaxPoints() int     { return c.maximum }
func (c categoryCheck) Run(_ context.Context, repository Repository) (Assessment, error) {
	assessment, err := c.run(repository)
	assessment.ID, assessment.Category, assessment.MaxPoints = c.id, c.category, c.maximum
	return assessment, err
}

// BuiltinChecks returns the complete v1 local repository policy.
func BuiltinChecks() []Check {
	return []Check{
		categoryCheck{"repository.completeness", CategoryCompleteness, 20, checkCompleteness},
		categoryCheck{"repository.documentation", CategoryDocumentation, 20, checkDocumentation},
		categoryCheck{"repository.ci", CategoryCI, 15, checkCI},
		categoryCheck{"repository.tests", CategoryTests, 15, checkTests},
		categoryCheck{"repository.license", CategoryLicense, 10, checkLicense},
		categoryCheck{"repository.security", CategorySecurity, 10, checkSecurity},
		categoryCheck{"repository.github", CategoryGitHub, 10, checkGitHub},
	}
}

func checkCompleteness(repository Repository) (Assessment, error) {
	assessment := Assessment{Title: "Repository completeness"}
	requirements := []struct {
		paths  []string
		points int
		label  string
		action string
	}{
		{[]string{"README.md", "README"}, 4, "README", "Add a root README explaining purpose, audience, installation, and first use."},
		{[]string{"CHANGELOG.md", "CHANGELOG"}, 3, "changelog", "Add a Keep a Changelog-compatible CHANGELOG.md."},
		{[]string{"ROADMAP.md"}, 2, "roadmap", "Document current and future themes in ROADMAP.md."},
		{[]string{"CONTRIBUTING.md", ".github/CONTRIBUTING.md"}, 3, "contribution guide", "Add exact setup, validation, and pull-request expectations."},
		{[]string{"CODE_OF_CONDUCT.md", ".github/CODE_OF_CONDUCT.md"}, 3, "code of conduct", "Add community behavior and enforcement guidance."},
		{[]string{"docs", "documentation"}, 2, "documentation directory", "Create a documentation entry point with user and contributor guidance."},
		{[]string{"Makefile", "Taskfile.yml", "justfile"}, 3, "developer task facade", "Add a documented, repeatable local check/test task."},
	}
	for _, requirement := range requirements {
		found, evidence, err := existsAny(repository, requirement.paths...)
		if err != nil {
			return assessment, err
		}
		if found {
			assessment.Points += requirement.points
			assessment.Evidence = append(assessment.Evidence, fmt.Sprintf("found %s at %s", requirement.label, evidence))
		} else {
			assessment.Recommendations = append(assessment.Recommendations, Recommendation{Message: "Missing " + requirement.label, Action: requirement.action})
		}
	}
	return assessment, nil
}

func checkDocumentation(repository Repository) (Assessment, error) {
	assessment := Assessment{Title: "Documentation quality"}
	readme, err := readFirst(repository, "README.md", "README")
	if err != nil {
		return assessment, err
	}
	if len(strings.TrimSpace(readme)) >= 200 {
		assessment.Points += 5
		assessment.Evidence = append(assessment.Evidence, "README contains substantive content")
	} else {
		assessment.Recommendations = append(assessment.Recommendations, Recommendation{"README is missing or too brief", "Explain what the project does, why it exists, and who should use it."})
	}
	lowerReadme := strings.ToLower(readme)
	if containsAny(lowerReadme, "install", "quick start", "getting started", "usage") {
		assessment.Points += 4
		assessment.Evidence = append(assessment.Evidence, "README includes onboarding or usage guidance")
	} else {
		assessment.Recommendations = append(assessment.Recommendations, Recommendation{"README lacks a runnable onboarding path", "Add installation and a copyable first-use example."})
	}
	if containsAny(lowerReadme, "license") && containsAny(lowerReadme, "contribut") {
		assessment.Points += 3
	} else {
		assessment.Recommendations = append(assessment.Recommendations, Recommendation{"README does not route readers to license and contribution guidance", "Link LICENSE and CONTRIBUTING.md from the README."})
	}
	docs, docsPath, err := readAny(repository, "docs/README.md", "docs/index.md", "docs/index.mdx")
	if err != nil {
		return assessment, err
	}
	if len(strings.TrimSpace(docs)) >= 100 {
		assessment.Points += 4
		assessment.Evidence = append(assessment.Evidence, "documentation entry point found at "+docsPath)
	} else {
		assessment.Recommendations = append(assessment.Recommendations, Recommendation{"Documentation entry point is absent or thin", "Add docs/README.md with guides, reference, and troubleshooting routes."})
	}
	contributing, err := readFirst(repository, "CONTRIBUTING.md", ".github/CONTRIBUTING.md")
	if err != nil {
		return assessment, err
	}
	if containsAny(strings.ToLower(contributing), "test", "check", "verify") && containsAny(strings.ToLower(contributing), "pull request", "review") {
		assessment.Points += 4
		assessment.Evidence = append(assessment.Evidence, "contribution guide includes validation and review expectations")
	} else {
		assessment.Recommendations = append(assessment.Recommendations, Recommendation{"Contribution guide is not actionable", "Document setup, the exact local validation command, and review expectations."})
	}
	return assessment, nil
}

func checkCI(repository Repository) (Assessment, error) {
	assessment := Assessment{Title: "Continuous integration"}
	workflows, matches, err := workflowText(repository)
	if err != nil {
		return assessment, err
	}
	if len(matches) == 0 {
		assessment.Recommendations = append(assessment.Recommendations, Recommendation{"No GitHub Actions workflows found", "Add CI under .github/workflows to validate pull requests and main."})
		return assessment, nil
	}
	assessment.Points += 6
	assessment.Evidence = append(assessment.Evidence, fmt.Sprintf("found %d workflow file(s)", len(matches)))
	if strings.Contains(workflows, "pull_request") {
		assessment.Points += 3
	} else {
		assessment.Recommendations = append(assessment.Recommendations, Recommendation{"CI does not appear to run for pull requests", "Add a pull_request trigger to validation workflows."})
	}
	actions := anyAction.FindAllString(workflows, -1)
	pinned := pinnedAction.FindAllString(workflows, -1)
	if len(actions) > 0 && len(actions) == len(pinned) {
		assessment.Points += 3
		assessment.Evidence = append(assessment.Evidence, "all detected action dependencies use full commit pins")
	} else {
		assessment.Recommendations = append(assessment.Recommendations, Recommendation{"One or more action dependencies are not pinned to full commits", "Replace action tags with reviewed 40-character commit SHAs and let Dependabot update them."})
	}
	if strings.Contains(workflows, "permissions:") && containsAny(workflows, "contents: read", "read-all") {
		assessment.Points += 3
	} else {
		assessment.Recommendations = append(assessment.Recommendations, Recommendation{"Workflow token permissions are not explicitly read-only by default", "Declare top-level permissions such as contents: read and elevate only individual jobs."})
	}
	return assessment, nil
}

func checkTests(repository Repository) (Assessment, error) {
	assessment := Assessment{Title: "Testing"}
	tests, err := repository.TestFiles()
	if err != nil {
		return assessment, err
	}
	if len(tests) > 0 {
		assessment.Points += 8
		assessment.Evidence = append(assessment.Evidence, fmt.Sprintf("found %d likely native test file(s)", len(tests)))
	} else {
		assessment.Recommendations = append(assessment.Recommendations, Recommendation{"No language-native test files detected", "Add tests using the ecosystem's standard naming and runner."})
	}
	workflows, _, err := workflowText(repository)
	if err != nil {
		return assessment, err
	}
	if containsAny(strings.ToLower(workflows), "go test", "pytest", "unittest", "cargo test", "npm test", "pnpm test", "bun test", "deno test", "node --test") {
		assessment.Points += 5
		assessment.Evidence = append(assessment.Evidence, "CI appears to invoke a native test runner")
	} else {
		assessment.Recommendations = append(assessment.Recommendations, Recommendation{"CI does not invoke a recognized native test runner", "Run the project's native test command in pull-request CI."})
	}
	testingDoc, err := readFirst(repository, "TESTING.md", "CONTRIBUTING.md", "README.md")
	if err != nil {
		return assessment, err
	}
	if containsAny(strings.ToLower(testingDoc), "test", "make check") {
		assessment.Points += 2
	} else {
		assessment.Recommendations = append(assessment.Recommendations, Recommendation{"Test execution is undocumented", "Document the exact local test command and required environment."})
	}
	return assessment, nil
}

func checkLicense(repository Repository) (Assessment, error) {
	assessment := Assessment{Title: "License"}
	license, licensePath, err := readAny(repository, "LICENSE", "LICENSE.md", "LICENSE.txt", "LICENCE", "COPYING")
	if err != nil {
		return assessment, err
	}
	if license == "" {
		assessment.Recommendations = append(assessment.Recommendations, Recommendation{"No root license file detected", "Choose a license explicitly and add its complete text in LICENSE."})
		return assessment, nil
	}
	assessment.Points += 5
	assessment.Evidence = append(assessment.Evidence, "license file found at "+licensePath)
	lower := strings.ToLower(license)
	if strings.Contains(lower, "mit license") || strings.Contains(lower, "apache license") || strings.Contains(lower, "mozilla public license") || strings.Contains(lower, "gnu general public license") {
		assessment.Points += 5
		assessment.Evidence = append(assessment.Evidence, "license text matches a recognized license family")
	} else {
		assessment.Recommendations = append(assessment.Recommendations, Recommendation{"License file is not locally recognizable", "Use complete standard license text and declare its SPDX identifier in project metadata."})
	}
	return assessment, nil
}

func checkSecurity(repository Repository) (Assessment, error) {
	assessment := Assessment{Title: "Security"}
	security, securityPath, err := readAny(repository, "SECURITY.md", ".github/SECURITY.md", "docs/SECURITY.md")
	if err != nil {
		return assessment, err
	}
	if security != "" {
		assessment.Points += 4
		assessment.Evidence = append(assessment.Evidence, "security policy found at "+securityPath)
	} else {
		assessment.Recommendations = append(assessment.Recommendations, Recommendation{"Security policy is missing", "Add supported versions and private vulnerability reporting instructions."})
	}
	lower := strings.ToLower(security)
	if containsAny(lower, "private", "email", "security advis", "do not open a public") && containsAny(lower, "supported", "version") {
		assessment.Points += 3
	} else {
		assessment.Recommendations = append(assessment.Recommendations, Recommendation{"Security policy lacks private reporting or support scope", "Name a private reporting channel, acknowledgement expectation, and supported versions."})
	}
	dependency, err := repository.Exists(".github/dependabot.yml")
	if err != nil {
		return assessment, err
	}
	if dependency {
		assessment.Points += 3
		assessment.Evidence = append(assessment.Evidence, "Dependabot configuration is present")
	} else {
		assessment.Recommendations = append(assessment.Recommendations, Recommendation{"Automated dependency updates are not configured", "Configure Dependabot or document an equivalent dependency-update process."})
	}
	return assessment, nil
}

func checkGitHub(repository Repository) (Assessment, error) {
	assessment := Assessment{Title: "GitHub configuration"}
	items := []struct {
		paths  []string
		points int
		label  string
		action string
	}{
		{[]string{".github/ISSUE_TEMPLATE"}, 3, "issue templates", "Add validated issue forms for bugs and proposals."},
		{[]string{".github/pull_request_template.md", ".github/PULL_REQUEST_TEMPLATE.md"}, 2, "pull request template", "Add summary, rationale, verification, and risk prompts."},
		{[]string{".github/dependabot.yml"}, 2, "dependency automation", "Configure dependency and Actions updates."},
		{[]string{".github/workflows/release.yml", ".github/workflows/release.yaml"}, 2, "release automation", "Add a reviewed semantic release workflow."},
		{[]string{".github/CODEOWNERS", "CODEOWNERS", "docs/CODEOWNERS"}, 1, "CODEOWNERS", "Declare real maintainers for review routing; do not add placeholders."},
	}
	for _, item := range items {
		found, evidence, err := existsAny(repository, item.paths...)
		if err != nil {
			return assessment, err
		}
		if found {
			assessment.Points += item.points
			assessment.Evidence = append(assessment.Evidence, "found "+item.label+" at "+evidence)
		} else {
			assessment.Recommendations = append(assessment.Recommendations, Recommendation{"Missing " + item.label, item.action})
		}
	}
	return assessment, nil
}

func existsAny(repository Repository, names ...string) (bool, string, error) {
	for _, name := range names {
		exists, err := repository.Exists(name)
		if err != nil {
			return false, "", err
		}
		if exists {
			return true, name, nil
		}
	}
	return false, "", nil
}

func readFirst(repository Repository, names ...string) (string, error) {
	content, _, err := readAny(repository, names...)
	return content, err
}

func readAny(repository Repository, names ...string) (string, string, error) {
	for _, name := range names {
		exists, err := repository.Exists(name)
		if err != nil {
			return "", "", err
		}
		if !exists {
			continue
		}
		content, err := repository.ReadText(name)
		return content, name, err
	}
	return "", "", nil
}

func workflowText(repository Repository) (string, []string, error) {
	yaml, err := repository.Glob(".github/workflows/*.yml")
	if err != nil {
		return "", nil, err
	}
	yamlLong, err := repository.Glob(".github/workflows/*.yaml")
	if err != nil {
		return "", nil, err
	}
	matches := append([]string(nil), yaml...)
	matches = append(matches, yamlLong...)
	var content strings.Builder
	for _, name := range matches {
		value, err := repository.ReadText(name)
		if err != nil {
			return "", nil, err
		}
		content.WriteString(value)
		content.WriteByte('\n')
	}
	return content.String(), matches, nil
}

func containsAny(value string, candidates ...string) bool {
	for _, candidate := range candidates {
		if strings.Contains(value, candidate) {
			return true
		}
	}
	return false
}
