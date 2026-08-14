package initialize

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/aruodore/aruo-cli/internal/contractmeta"
)

//go:embed contract
var contractFiles embed.FS

// Service plans and installs the embedded contract without touching application code.
type Service struct{}

// NewService constructs the stateless initializer.
func NewService() *Service { return &Service{} }

// Plan inspects target and returns every proposed write without mutating it.
func (s *Service) Plan(ctx context.Context, request Request) (Plan, error) {
	if err := ctx.Err(); err != nil {
		return Plan{}, err
	}
	absolute, err := filepath.Abs(request.Target)
	if err != nil {
		return Plan{}, fmt.Errorf("resolve repository: %w", err)
	}
	info, err := os.Stat(absolute)
	if err != nil || !info.IsDir() {
		return Plan{}, fmt.Errorf("repository is not an accessible directory: %s", absolute)
	}
	root, err := os.OpenRoot(absolute)
	if err != nil {
		return Plan{}, fmt.Errorf("open repository: %w", err)
	}
	defer root.Close()
	stack, err := detectStack(root.FS())
	if err != nil {
		return Plan{}, fmt.Errorf("detect stack: %w", err)
	}
	files, err := renderContract(stack)
	if err != nil {
		return Plan{}, err
	}
	plan := Plan{Repository: absolute, Stack: stack, files: files}
	managedRootExists := exists(root.FS(), ".aruo")
	if managedRootExists {
		plan.Conflicts = append(plan.Conflicts, ".aruo")
	}
	for _, file := range files {
		if managedRootExists && strings.HasPrefix(file.path, ".aruo/") {
			continue
		}
		if exists(root.FS(), file.path) {
			plan.Conflicts = append(plan.Conflicts, file.path)
			continue
		}
		plan.Changes = append(plan.Changes, FileChange{Path: file.path, Action: "create"})
	}
	return plan, nil
}

// Apply commits a previously computed conflict-free plan through same-parent staging.
func (s *Service) Apply(ctx context.Context, plan Plan) (Result, error) {
	if len(plan.Conflicts) > 0 {
		return Result{}, fmt.Errorf("initialization conflicts with existing paths: %s", strings.Join(plan.Conflicts, ", "))
	}
	if len(plan.files) == 0 {
		return Result{}, errors.New("initialization plan has no files")
	}
	staging, err := os.MkdirTemp(plan.Repository, ".aruo-init-*")
	if err != nil {
		return Result{}, fmt.Errorf("create initialization staging directory: %w", err)
	}
	defer os.RemoveAll(staging)
	for _, file := range plan.files {
		if contextErr := ctx.Err(); contextErr != nil {
			return Result{}, contextErr
		}
		target := filepath.Join(staging, filepath.FromSlash(file.path))
		if mkdirErr := os.MkdirAll(filepath.Dir(target), 0o755); mkdirErr != nil {
			return Result{}, fmt.Errorf("stage %s: %w", file.path, mkdirErr)
		}
		if writeErr := os.WriteFile(target, file.content, 0o600); writeErr != nil {
			return Result{}, fmt.Errorf("stage %s: %w", file.path, writeErr)
		}
		// Contract files contain no secrets and must remain readable by repository
		// collaborators and CI users after the staged commit.
		if chmodErr := os.Chmod(target, 0o644); chmodErr != nil { //nolint:gosec // Repository governance files are intentionally world-readable.
			return Result{}, fmt.Errorf("set repository-readable mode on %s: %w", file.path, chmodErr)
		}
	}
	committed := make([]string, 0, 3)
	rollback := func() {
		for _, name := range committed {
			_ = os.RemoveAll(filepath.Join(plan.Repository, name))
		}
	}
	managedRoot := filepath.Join(plan.Repository, ".aruo")
	if mkdirErr := os.Mkdir(managedRoot, 0o755); mkdirErr != nil {
		return Result{}, fmt.Errorf("commit .aruo without overwriting: %w", mkdirErr)
	}
	committed = append(committed, ".aruo")
	managedEntries, err := os.ReadDir(filepath.Join(staging, ".aruo"))
	if err != nil {
		rollback()
		return Result{}, fmt.Errorf("inspect staged contract: %w", err)
	}
	for _, entry := range managedEntries {
		if err := os.Rename(filepath.Join(staging, ".aruo", entry.Name()), filepath.Join(managedRoot, entry.Name())); err != nil {
			rollback()
			return Result{}, fmt.Errorf("commit .aruo/%s: %w", entry.Name(), err)
		}
	}
	for _, topLevel := range []string{"AGENTS.md", "aruo.yaml"} {
		if err := ctx.Err(); err != nil {
			rollback()
			return Result{}, err
		}
		// Link is intentionally used instead of rename: it fails when the target
		// exists, preserving user content even if the repository changed after Plan.
		if err := os.Link(filepath.Join(staging, topLevel), filepath.Join(plan.Repository, topLevel)); err != nil {
			rollback()
			return Result{}, fmt.Errorf("commit %s without overwriting: %w", topLevel, err)
		}
		committed = append(committed, topLevel)
	}
	return Result{Repository: plan.Repository, Stack: plan.Stack, Changes: plan.Changes}, nil
}

func renderContract(stack Stack) ([]plannedFile, error) {
	return renderContractWithOverrides(stack, nil)
}

func renderContractWithOverrides(stack Stack, overrides map[string][]byte) ([]plannedFile, error) {
	entries, err := fs.ReadDir(contractFiles, "contract")
	if err != nil {
		return nil, err
	}
	files := make([]plannedFile, 0, len(entries)+3)
	for _, entry := range entries {
		if entry.IsDir() {
			children, readErr := fs.ReadDir(contractFiles, "contract/"+entry.Name())
			if readErr != nil {
				return nil, readErr
			}
			for _, child := range children {
				content, readErr := contractFiles.ReadFile("contract/" + entry.Name() + "/" + child.Name())
				if readErr != nil {
					return nil, readErr
				}
				files = append(files, plannedFile{path: ".aruo/" + entry.Name() + "/" + child.Name(), content: content})
			}
			continue
		}
		content, readErr := contractFiles.ReadFile("contract/" + entry.Name())
		if readErr != nil {
			return nil, readErr
		}
		destination := entry.Name()
		if entry.Name() == "contract.yaml" {
			destination = ".aruo/contract.yaml"
		}
		if override, ok := overrides[destination]; ok {
			content = override
		}
		files = append(files, plannedFile{path: destination, content: content})
	}
	stackContent := renderStack(stack)
	files = append(files, plannedFile{path: ".aruo/stack.yaml", content: stackContent})
	sort.Slice(files, func(i, j int) bool { return files[i].path < files[j].path })
	managed := struct {
		ContractVersion string            `json:"contractVersion"`
		Files           map[string]string `json:"files"`
	}{ContractVersion: contractmeta.CurrentVersion, Files: map[string]string{}}
	for _, file := range files {
		if file.path == "aruo.yaml" {
			continue
		}
		sum := sha256.Sum256(file.content)
		managed.Files[file.path] = "sha256:" + hex.EncodeToString(sum[:])
	}
	required, _ := contractmeta.RequiredFiles(contractmeta.CurrentVersion)
	if inventoryErr := validateManagedInventory(managed.Files, required); inventoryErr != nil {
		return nil, inventoryErr
	}
	managedContent, err := json.MarshalIndent(managed, "", "  ")
	if err != nil {
		return nil, err
	}
	managedContent = append(managedContent, '\n')
	files = append(files, plannedFile{path: ".aruo/managed.json", content: managedContent})
	sort.Slice(files, func(i, j int) bool { return files[i].path < files[j].path })
	return files, nil
}

func validateManagedInventory(files map[string]string, required []string) error {
	expected := make(map[string]struct{}, len(required))
	for _, path := range required {
		expected[path] = struct{}{}
		if files[path] == "" {
			return fmt.Errorf("embedded contract version %s is missing managed file %q", contractmeta.CurrentVersion, path)
		}
	}
	for path := range files {
		if _, present := expected[path]; !present {
			return fmt.Errorf("embedded contract version %s contains unexpected managed file %q", contractmeta.CurrentVersion, path)
		}
	}
	return nil
}

func renderStack(stack Stack) []byte {
	var builder strings.Builder
	builder.WriteString("apiVersion: aruo.dev/v1alpha1\necosystem: ")
	builder.WriteString(stack.Ecosystem)
	builder.WriteString("\npackageManager: ")
	builder.WriteString(stack.PackageManager)
	builder.WriteString("\nframeworks:")
	if len(stack.Frameworks) == 0 {
		builder.WriteString(" []\n")
	} else {
		builder.WriteByte('\n')
		for _, framework := range stack.Frameworks {
			builder.WriteString("  - ")
			builder.WriteString(framework)
			builder.WriteByte('\n')
		}
	}
	return []byte(builder.String())
}
