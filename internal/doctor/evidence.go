package doctor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"strings"
)

type evidenceVerification struct {
	recognized bool
	verified   bool
	action     string
}

type packageManifest struct {
	Scripts map[string]string `json:"scripts"`
}

func verifySemanticEvidence(repository Repository, evidence string) (evidenceVerification, error) {
	switch evidence {
	case "npm-run-check":
		return verifyPackageScript(repository, "check", []string{"format", "lint", "test", "build"})
	case "vitest":
		verification, err := verifyPackageScript(repository, "test", []string{"vitest"})
		if err != nil || !verification.verified {
			return verification, err
		}
		tests, err := repository.TestFiles()
		if err != nil {
			return evidenceVerification{}, err
		}
		verification.verified = len(tests) > 0
		verification.action = "Add at least one discoverable test file for the declared test runner."
		return verification, nil
	case "vite-build":
		return verifyPackageScript(repository, "build", []string{"vite build"})
	case "next-build":
		return verifyPackageScript(repository, "build", []string{"next build"})
	case "nuxt-build":
		return verifyPackageScript(repository, "build", []string{"nuxt build"})
	case "nuxt-typecheck":
		return verifyPackageScript(repository, "typecheck", []string{"nuxt typecheck"})
	case "strict-typescript":
		return verifyStrictTypeScript(repository)
	case "ci-high-severity-gate":
		return verifyWorkflowText(repository, []string{"npm audit", "--audit-level=high"}, "Run a high-severity dependency audit in a committed CI workflow.")
	case "committed-drizzle-sql":
		return verifyGlob(repository, "server/db/migrations/*.sql", "Commit at least one SQL migration under server/db/migrations.")
	case "live-and-ready":
		return verifyPathSet(repository, []string{"server/api/health/live.get.ts", "server/api/health/ready.get.ts"}, "Add both Nuxt liveness and readiness handlers.")
	case "live-and-ready-routes":
		return verifyHealthRoutePair(repository)
	default:
		return evidenceVerification{}, nil
	}
}

func verifyPackageScript(repository Repository, script string, required []string) (evidenceVerification, error) {
	manifest, err := readPackageManifest(repository)
	if err != nil {
		return evidenceVerification{}, err
	}
	command := strings.ToLower(manifest.Scripts[script])
	verification := evidenceVerification{recognized: true, verified: command != "", action: fmt.Sprintf("Define package.json script %q with the declared checks.", script)}
	for _, fragment := range required {
		if !strings.Contains(command, fragment) {
			verification.verified = false
			break
		}
	}
	return verification, nil
}

func readPackageManifest(repository Repository) (packageManifest, error) {
	content, err := repository.ReadText("package.json")
	if errors.Is(err, fs.ErrNotExist) {
		return packageManifest{}, nil
	}
	if err != nil {
		return packageManifest{}, err
	}
	var manifest packageManifest
	if err := json.Unmarshal([]byte(content), &manifest); err != nil {
		return packageManifest{}, fmt.Errorf("parse package.json for intent evidence: %w", err)
	}
	return manifest, nil
}

func verifyStrictTypeScript(repository Repository) (evidenceVerification, error) {
	verification := evidenceVerification{recognized: true, action: "Enable compilerOptions.strict or extend a recognized strict TypeScript base configuration."}
	matches, err := repository.Glob("tsconfig*.json")
	if err != nil {
		return evidenceVerification{}, err
	}
	for _, match := range matches {
		content, readErr := repository.ReadText(match)
		if readErr != nil {
			return evidenceVerification{}, readErr
		}
		compact := strings.NewReplacer(" ", "", "\n", "", "\r", "", "\t", "").Replace(content)
		if strings.Contains(compact, `"strict":true`) || strings.Contains(compact, `"extends":"@vue/tsconfig/`) {
			verification.verified = true
			break
		}
	}
	return verification, nil
}

func verifyWorkflowText(repository Repository, required []string, action string) (evidenceVerification, error) {
	var matches []string
	for _, pattern := range []string{".github/workflows/*.yml", ".github/workflows/*.yaml"} {
		patternMatches, err := repository.Glob(pattern)
		if err != nil {
			return evidenceVerification{}, err
		}
		matches = append(matches, patternMatches...)
	}
	verification := evidenceVerification{recognized: true, action: action}
	for _, match := range matches {
		content, readErr := repository.ReadText(match)
		if readErr != nil {
			return evidenceVerification{}, readErr
		}
		content = strings.ToLower(content)
		found := true
		for _, fragment := range required {
			found = found && strings.Contains(content, fragment)
		}
		if found {
			verification.verified = true
			break
		}
	}
	return verification, nil
}

func verifyGlob(repository Repository, pattern, action string) (evidenceVerification, error) {
	matches, err := repository.Glob(pattern)
	if err != nil {
		return evidenceVerification{}, err
	}
	return evidenceVerification{recognized: true, verified: len(matches) > 0, action: action}, nil
}

func verifyPathSet(repository Repository, names []string, action string) (evidenceVerification, error) {
	verification := evidenceVerification{recognized: true, verified: true, action: action}
	for _, name := range names {
		exists, err := repository.Exists(name)
		if err != nil {
			return evidenceVerification{}, err
		}
		verification.verified = verification.verified && exists
	}
	return verification, nil
}

func verifyHealthRoutePair(repository Repository) (evidenceVerification, error) {
	pairs := [][]string{
		{"app/api/health/live/route.ts", "app/api/health/ready/route.ts"},
		{"src/app/api/health/live/route.ts", "src/app/api/health/ready/route.ts"},
	}
	for _, pair := range pairs {
		verification, err := verifyPathSet(repository, pair, "Add distinct liveness and readiness routes.")
		if err != nil {
			return evidenceVerification{}, err
		}
		if verification.verified {
			return verification, nil
		}
	}
	content, err := repository.ReadText("server/index.ts")
	if err == nil {
		return evidenceVerification{
			recognized: true,
			verified:   strings.Contains(content, "/health/live") && strings.Contains(content, "/health/ready"),
			action:     "Add distinct liveness and readiness routes.",
		}, nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return evidenceVerification{}, err
	}
	return evidenceVerification{recognized: true, action: "Add distinct liveness and readiness routes."}, nil
}
