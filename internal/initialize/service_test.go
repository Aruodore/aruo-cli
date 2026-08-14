package initialize

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/aruodore/aruo-cli/internal/contractmeta"
)

var expectedInstalledContractPaths = []string{
	".aruo/contract.yaml",
	".aruo/managed.json",
	".aruo/rules/api.md",
	".aruo/rules/architecture.md",
	".aruo/rules/data.md",
	".aruo/rules/delivery.md",
	".aruo/rules/observability.md",
	".aruo/rules/security.md",
	".aruo/rules/testing.md",
	".aruo/stack.yaml",
	"AGENTS.md",
	"aruo.yaml",
}

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
	if len(plan.Conflicts) != 0 {
		t.Fatalf("plan = %#v, want a conflict-free complete contract", plan)
	}
	plannedPaths := make([]string, 0, len(plan.Changes))
	for _, change := range plan.Changes {
		plannedPaths = append(plannedPaths, change.Path)
	}
	sort.Strings(plannedPaths)
	if !equalStrings(plannedPaths, expectedInstalledContractPaths) {
		t.Fatalf("planned paths = %#v, want %#v", plannedPaths, expectedInstalledContractPaths)
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
	for _, name := range expectedInstalledContractPaths {
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
	if managed.ContractVersion != contractmeta.CurrentVersion || len(managed.Files) != 10 || managed.Files["AGENTS.md"] == "" || managed.Files["aruo.yaml"] != "" || managed.Files[".aruo/managed.json"] != "" {
		t.Fatalf("managed manifest = %#v, want managed AGENTS and application-owned aruo.yaml", managed)
	}
	for _, name := range expectedInstalledContractPaths {
		if name == "aruo.yaml" || name == ".aruo/managed.json" {
			continue
		}
		if managed.Files[name] == "" {
			t.Errorf("managed manifest is missing %s", name)
		}
	}
	for _, name := range []string{"AGENTS.md", "aruo.yaml", ".aruo/contract.yaml", ".aruo/rules/security.md"} {
		info, statErr := os.Stat(filepath.Join(repository, filepath.FromSlash(name)))
		if statErr != nil {
			t.Fatal(statErr)
		}
		if got := info.Mode().Perm(); got != 0o644 {
			t.Errorf("%s mode = %o, want 644", name, got)
		}
	}
	if info, statErr := os.Stat(filepath.Join(repository, ".aruo")); statErr != nil {
		t.Errorf("stat .aruo: %v", statErr)
	} else if info.Mode().Perm() != 0o755 {
		t.Errorf(".aruo mode = %o, want 755", info.Mode().Perm())
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

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
