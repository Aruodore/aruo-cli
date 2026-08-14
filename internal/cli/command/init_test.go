package command

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aruodore/aruo-cli/internal/cli/iostreams"
	"github.com/aruodore/aruo-cli/internal/initialize"
)

func TestInitCommandDryRunJSONDoesNotWrite(t *testing.T) {
	t.Parallel()
	repository := t.TempDir()
	var out bytes.Buffer
	factory := initTestFactory(&out)
	command := newInit(factory, initialize.NewService())
	command.SetArgs([]string{repository, "--dry-run", "--format", "json"})
	command.SilenceErrors = true
	command.SilenceUsage = true
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	var result initialize.Result
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("decode output %q: %v", out.String(), err)
	}
	if !result.DryRun || len(result.Changes) == 0 {
		t.Fatalf("result = %#v, want non-empty dry-run plan", result)
	}
	if _, err := os.Stat(filepath.Join(repository, ".aruo")); !os.IsNotExist(err) {
		t.Fatalf("dry run wrote .aruo: %v", err)
	}
}

func TestInitCommandAppliesWithYes(t *testing.T) {
	t.Parallel()
	repository := t.TempDir()
	var out bytes.Buffer
	command := newInit(initTestFactory(&out), initialize.NewService())
	command.SetArgs([]string{repository, "--yes", "--format", "json"})
	command.SilenceErrors = true
	command.SilenceUsage = true
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(repository, ".aruo/managed.json")); err != nil {
		t.Fatalf("managed manifest missing: %v", err)
	}
}

func TestInitCommandReportsConflict(t *testing.T) {
	t.Parallel()
	repository := t.TempDir()
	if err := os.WriteFile(filepath.Join(repository, "aruo.yaml"), []byte("existing"), 0o600); err != nil {
		t.Fatal(err)
	}
	command := newInit(initTestFactory(&bytes.Buffer{}), initialize.NewService())
	command.SetArgs([]string{repository, "--yes"})
	command.SilenceErrors = true
	command.SilenceUsage = true
	err := command.Execute()
	if err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("Execute() error = %v, want overwrite refusal", err)
	}
}

func initTestFactory(out *bytes.Buffer) sessionFactory {
	return sessionFactory{
		streams:     iostreams.IOStreams{In: strings.NewReader(""), Out: out, ErrOut: &bytes.Buffer{}},
		environment: map[string]string{},
		probe:       fakeProbe{},
		global:      newGlobalOptions(),
	}
}
