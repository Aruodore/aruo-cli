package command

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/aruodore/aruo-cli/internal/cli/iostreams"
	"github.com/aruodore/aruo-cli/internal/clierror"
	"github.com/aruodore/aruo-cli/internal/doctor"
)

// zeroCheckService builds a doctor.Service with no registered checks, so
// Audit against any real directory deterministically scores 0/0 (grade N/A)
// without depending on repository content.
func zeroCheckService(t *testing.T) *doctor.Service {
	t.Helper()
	engine, err := doctor.NewEngine()
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}
	service, err := doctor.NewService(engine)
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}
	return service
}

func TestDoctorCommandRejectsMinimumScoreBelowZero(t *testing.T) {
	t.Parallel()

	command := newDoctor(sessionFactory{}, zeroCheckService(t))
	command.SetArgs([]string{".", "--minimum-score", "-1"})
	command.SilenceErrors = true
	command.SilenceUsage = true
	if err := command.Execute(); err == nil {
		t.Fatal("Execute() error = nil, want error for --minimum-score below 0")
	}
}

func TestDoctorCommandRejectsMinimumScoreAboveHundred(t *testing.T) {
	t.Parallel()

	command := newDoctor(sessionFactory{}, zeroCheckService(t))
	command.SetArgs([]string{".", "--minimum-score", "150"})
	command.SilenceErrors = true
	command.SilenceUsage = true
	if err := command.Execute(); err == nil {
		t.Fatal("Execute() error = nil, want error for --minimum-score above 100")
	}
}

func TestDoctorCommandRejectsUnsupportedFormat(t *testing.T) {
	t.Parallel()

	command := newDoctor(sessionFactory{}, zeroCheckService(t))
	command.SetArgs([]string{".", "--format", "yaml"})
	command.SilenceErrors = true
	command.SilenceUsage = true
	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want error for unsupported --format value")
	}
	const want = `unsupported format "yaml"; use human or json`
	if err.Error() != want {
		t.Fatalf("Execute() error = %q, want %q", err.Error(), want)
	}
}

func TestDoctorCommandJSONFormatReturnsFindingsExitCodeBelowMinimum(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	factory := sessionFactory{streams: iostreams.IOStreams{Out: &out}}
	command := newDoctor(factory, zeroCheckService(t))
	// Default --minimum-score is 80; a zero-check engine always scores 0.
	command.SetArgs([]string{".", "--format", "json"})
	command.SilenceErrors = true
	command.SilenceUsage = true

	err := command.Execute()
	var clier *clierror.Error
	if !errors.As(err, &clier) {
		t.Fatalf("Execute() error = %v, want *clierror.Error", err)
	}
	if clier.ExitCode() != 3 {
		t.Errorf("ExitCode() = %d, want 3", clier.ExitCode())
	}
	if !clier.SuppressMessage() {
		t.Error("SuppressMessage() = false, want true since the report was already printed")
	}
	if out.Len() == 0 {
		t.Error("JSON report was not written to the output stream before the threshold error")
	}
}

func TestDoctorCommandJSONFormatSucceedsAtMinimumScoreZero(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	factory := sessionFactory{streams: iostreams.IOStreams{Out: &out}}
	command := newDoctor(factory, zeroCheckService(t))
	command.SetArgs([]string{".", "--format", "json", "--minimum-score", "0"})
	command.SilenceErrors = true
	command.SilenceUsage = true

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !bytes.Contains(out.Bytes(), []byte(`"schemaVersion"`)) {
		t.Errorf("output = %q, want a JSON report containing schemaVersion", out.String())
	}
}

func TestDoctorCommandReturnsFindingsExitCodeForRequiredIntent(t *testing.T) {
	t.Parallel()
	directory := t.TempDir()
	manifest := "apiVersion: aruo.dev/v1alpha1\nintent:\n  capabilities:\n    authentication: { status: REQUIRED, reason: identity-not-selected }\n"
	if err := os.WriteFile(filepath.Join(directory, "aruo.yaml"), []byte(manifest), 0o600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	command := newDoctor(sessionFactory{streams: iostreams.IOStreams{Out: &out}}, zeroCheckService(t))
	command.SetArgs([]string{directory, "--format", "json", "--minimum-score", "0"})
	command.SilenceErrors = true
	command.SilenceUsage = true
	err := command.Execute()
	var clier *clierror.Error
	if !errors.As(err, &clier) || clier.ExitCode() != 3 {
		t.Fatalf("Execute() error = %v, want findings exit code 3", err)
	}
	if !bytes.Contains(out.Bytes(), []byte(`"blockingFindings": 1`)) {
		t.Fatalf("output = %q, want structured blocking intent finding", out.String())
	}
}
