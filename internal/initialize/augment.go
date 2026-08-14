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
	if guidance, ok := existing["AGENTS.md"]; ok {
		base, readErr := contractFiles.ReadFile("contract/AGENTS.md")
		if readErr != nil {
			return templateengine.Plan{}, fmt.Errorf("read base agent contract: %w", readErr)
		}
		combined := append(append(append([]byte(nil), base...), []byte("\n## Stack-specific guidance\n\n")...), guidance.Content...)
		overrides["AGENTS.md"] = combined
	}
	contract, err := renderContractWithOverrides(stack, overrides)
	if err != nil {
		return templateengine.Plan{}, err
	}

	result := templateengine.Plan{BlueprintID: plan.BlueprintID, Files: append([]templateengine.File(nil), plan.Files...)}
	for _, file := range contract {
		if current, exists := existing[file.path]; exists {
			if file.path != "AGENTS.md" && file.path != "aruo.yaml" {
				return templateengine.Plan{}, fmt.Errorf("creation plan path %q conflicts with managed Aruo contract", file.path)
			}
			if file.path == "AGENTS.md" {
				for index := range result.Files {
					if result.Files[index].Path == file.path {
						result.Files[index].Content = append([]byte(nil), file.content...)
						result.Files[index].Source = "aruo:contract+" + current.Source
					}
				}
			}
			// aruo.yaml remains application-owned and is intentionally not
			// listed in managed.json.
			continue
		}
		result.Files = append(result.Files, templateengine.File{
			Path: file.path, Content: append([]byte(nil), file.content...), Mode: 0o600, Source: "aruo:contract",
		})
	}
	sort.Slice(result.Files, func(i, j int) bool { return result.Files[i].Path < result.Files[j].Path })
	return result, nil
}
