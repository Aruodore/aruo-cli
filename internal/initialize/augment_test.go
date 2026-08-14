package initialize

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/aruodore/aruo-cli/internal/templateengine"
)

func TestAugmentPlanPreservesIntentAndCombinesStackGuidance(t *testing.T) {
	t.Parallel()

	intent := "apiVersion: aruo.dev/v1alpha1\ntemplate: next\n"
	guidance := "Keep server code out of client modules.\n"
	plan, err := AugmentPlan(templateengine.Plan{BlueprintID: "next", Files: []templateengine.File{
		{Path: "package.json", Content: []byte(`{"dependencies":{"next":"latest","react":"latest"}}`)},
		{Path: "aruo.yaml", Content: []byte(intent)},
		{Path: "AGENTS.md", Content: []byte(guidance), Source: "next/AGENTS.md"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	files := planContents(plan)
	if string(files["aruo.yaml"]) != intent {
		t.Fatalf("aruo.yaml = %q, want application-owned intent preserved", files["aruo.yaml"])
	}
	agents := string(files["AGENTS.md"])
	if !strings.Contains(agents, "Read `.aruo/contract.yaml`") || !strings.Contains(agents, guidance) {
		t.Fatalf("AGENTS.md did not combine managed and stack guidance: %q", agents)
	}
	if !strings.Contains(string(files[".aruo/stack.yaml"]), "  - next\n") {
		t.Fatalf("stack.yaml = %q, want detected Next.js stack", files[".aruo/stack.yaml"])
	}
	architecture := string(files[".aruo/rules/architecture.md"])
	if !strings.Contains(architecture, "lowercase kebab-case") {
		t.Errorf("architecture rules do not require kebab-case filenames: %q", architecture)
	}

	var managed struct {
		Files map[string]string `json:"files"`
	}
	if err := json.Unmarshal(files[".aruo/managed.json"], &managed); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(files["AGENTS.md"])
	if got, want := managed.Files["AGENTS.md"], "sha256:"+hex.EncodeToString(sum[:]); got != want {
		t.Fatalf("managed AGENTS hash = %q, want %q", got, want)
	}
	if _, managedIntent := managed.Files["aruo.yaml"]; managedIntent {
		t.Fatal("application-owned aruo.yaml appears in managed hashes")
	}
}

func TestAugmentPlanRejectsDuplicateApplicationPaths(t *testing.T) {
	t.Parallel()

	_, err := AugmentPlan(templateengine.Plan{Files: []templateengine.File{{Path: "README.md"}, {Path: "README.md"}}})
	if err == nil || !strings.Contains(err.Error(), "duplicate path") {
		t.Fatalf("AugmentPlan() error = %v, want duplicate path rejection", err)
	}
}

func planContents(plan templateengine.Plan) map[string][]byte {
	result := make(map[string][]byte, len(plan.Files))
	for _, file := range plan.Files {
		result[file.Path] = file.Content
	}
	return result
}
