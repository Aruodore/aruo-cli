package initialize

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestServicePlansAndAppliesContractWithoutChangingApplication(t *testing.T) {
	t.Parallel()
	repository := t.TempDir()
	packageContent := `{"dependencies":{"next":"latest","react":"latest"},"devDependencies":{"vite":"latest"}}`
	writeTestFile(t, repository, "package.json", packageContent)
	writeTestFile(t, repository, "package-lock.json", "{}")
	service := NewService()
	plan, err := service.Plan(context.Background(), Request{Target: repository})
	if err != nil {
		t.Fatal(err)
	}
	if plan.Stack.Ecosystem != "node" || plan.Stack.PackageManager != "npm" {
		t.Fatalf("stack = %#v, want node/npm", plan.Stack)
	}
	if got := plan.Stack.Frameworks; len(got) != 3 || got[0] != "next" || got[1] != "react" || got[2] != "vite" {
		t.Fatalf("frameworks = %#v, want sorted next/react/vite", got)
	}
	if len(plan.Conflicts) != 0 || len(plan.Changes) < 10 {
		t.Fatalf("plan = %#v, want a conflict-free complete contract", plan)
	}
	result, err := service.Apply(context.Background(), plan)
	if err != nil {
		t.Fatal(err)
	}
	if result.DryRun || len(result.Changes) != len(plan.Changes) {
		t.Fatalf("result = %#v, want applied plan", result)
	}
	unchanged, err := os.ReadFile(filepath.Join(repository, "package.json"))
	if err != nil || string(unchanged) != packageContent {
		t.Fatalf("application package changed: %q, %v", unchanged, err)
	}
	for _, name := range []string{"AGENTS.md", "aruo.yaml", ".aruo/contract.yaml", ".aruo/stack.yaml", ".aruo/managed.json", ".aruo/rules/security.md"} {
		if _, statErr := os.Stat(filepath.Join(repository, filepath.FromSlash(name))); statErr != nil {
			t.Errorf("expected %s: %v", name, statErr)
		}
	}
	managedContent, err := os.ReadFile(filepath.Join(repository, ".aruo/managed.json"))
	if err != nil {
		t.Fatal(err)
	}
	var managed struct {
		ContractVersion string            `json:"contractVersion"`
		Files           map[string]string `json:"files"`
	}
	if err := json.Unmarshal(managedContent, &managed); err != nil {
		t.Fatal(err)
	}
	if managed.ContractVersion != contractVersion || managed.Files["AGENTS.md"] == "" || managed.Files["aruo.yaml"] != "" {
		t.Fatalf("managed manifest = %#v, want managed AGENTS and application-owned aruo.yaml", managed)
	}
}

func TestServiceRefusesExistingContractPathWithoutWriting(t *testing.T) {
	t.Parallel()
	repository := t.TempDir()
	writeTestFile(t, repository, "AGENTS.md", "existing instructions")
	service := NewService()
	plan, err := service.Plan(context.Background(), Request{Target: repository})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Conflicts) != 1 || plan.Conflicts[0] != "AGENTS.md" {
		t.Fatalf("conflicts = %#v, want AGENTS.md", plan.Conflicts)
	}
	if _, applyErr := service.Apply(context.Background(), plan); applyErr == nil {
		t.Fatal("Apply() error = nil, want conflict error")
	}
	if _, err := os.Stat(filepath.Join(repository, ".aruo")); !os.IsNotExist(err) {
		t.Fatalf(".aruo exists after refused apply: %v", err)
	}
}

func TestServiceRejectsMalformedPackageManifest(t *testing.T) {
	t.Parallel()
	repository := t.TempDir()
	writeTestFile(t, repository, "package.json", "{")
	if _, err := NewService().Plan(context.Background(), Request{Target: repository}); err == nil {
		t.Fatal("Plan() error = nil, want malformed package.json error")
	}
}

func TestServiceReportsExistingManagedRootAsConflict(t *testing.T) {
	t.Parallel()
	repository := t.TempDir()
	writeTestFile(t, repository, ".aruo/local-note", "user content")
	plan, err := NewService().Plan(context.Background(), Request{Target: repository})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Conflicts) != 1 || plan.Conflicts[0] != ".aruo" {
		t.Fatalf("conflicts = %#v, want managed-root conflict", plan.Conflicts)
	}
}

func TestServiceDoesNotOverwriteFileCreatedAfterPlanning(t *testing.T) {
	t.Parallel()
	repository := t.TempDir()
	service := NewService()
	plan, err := service.Plan(context.Background(), Request{Target: repository})
	if err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, repository, "AGENTS.md", "created after planning")
	if _, applyErr := service.Apply(context.Background(), plan); applyErr == nil {
		t.Fatal("Apply() error = nil, want no-overwrite conflict")
	}
	content, err := os.ReadFile(filepath.Join(repository, "AGENTS.md"))
	if err != nil || string(content) != "created after planning" {
		t.Fatalf("racing user file changed: %q, %v", content, err)
	}
	if _, statErr := os.Stat(filepath.Join(repository, ".aruo")); !os.IsNotExist(statErr) {
		t.Fatalf("managed directory remains after rollback: %v", statErr)
	}
}

func writeTestFile(t *testing.T, root, name, content string) {
	t.Helper()
	target := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
