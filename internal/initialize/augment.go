package initialize

import (
	"fmt"
	"sort"
	"testing/fstest"

	"github.com/aruodore/aruo-cli/internal/templateengine"
)

// AugmentPlan merges the managed Aruo contract into a rendered creation plan.
// The caller can then commit application files and governance files through
// one atomic writer operation.
func AugmentPlan(plan templateengine.Plan) (templateengine.Plan, error) {
	contents := make(fstest.MapFS, len(plan.Files))
	existing := make(map[string]templateengine.File, len(plan.Files))
	for _, file := range plan.Files {
		if _, duplicate := existing[file.Path]; duplicate {
			return templateengine.Plan{}, fmt.Errorf("creation plan contains duplicate path %q", file.Path)
		}
		existing[file.Path] = file
		contents[file.Path] = &fstest.MapFile{Data: file.Content, Mode: file.Mode}
	}
	stack, err := detectStack(contents)
	if err != nil {
		return templateengine.Plan{}, fmt.Errorf("detect rendered stack: %w", err)
	}
	overrides := map[string][]byte{}
	if intent, ok := existing["aruo.yaml"]; ok {
		overrides["aruo.yaml"] = intent.Content
	}
	guidance, hasGuidance := existing["AGENTS.md"]
	if hasGuidance {
		if _, conflict := existing["AGENTS.local.md"]; conflict {
			return templateengine.Plan{}, fmt.Errorf("creation plan contains both AGENTS.md guidance and AGENTS.local.md")
		}
		delete(existing, "AGENTS.md")
	}
	contract, err := renderContractWithOverrides(stack, overrides)
	if err != nil {
		return templateengine.Plan{}, err
	}

	result := templateengine.Plan{BlueprintID: plan.BlueprintID, Files: make([]templateengine.File, 0, len(plan.Files)+len(contract))}
	for _, file := range plan.Files {
		if hasGuidance && file.Path == "AGENTS.md" {
			file.Path = "AGENTS.local.md"
			file.Source = guidance.Source
		}
		result.Files = append(result.Files, file)
	}
	for _, file := range contract {
		if _, exists := existing[file.path]; exists {
			if file.path != "aruo.yaml" {
				return templateengine.Plan{}, fmt.Errorf("creation plan path %q conflicts with managed Aruo contract", file.path)
			}
			// aruo.yaml remains application-owned and is intentionally not
			// listed in managed.json.
			continue
		}
		result.Files = append(result.Files, templateengine.File{
			Path: file.path, Content: append([]byte(nil), file.content...), Mode: 0o644, Source: "aruo:contract",
		})
	}
	sort.Slice(result.Files, func(i, j int) bool { return result.Files[i].Path < result.Files[j].Path })
	return result, nil
}
