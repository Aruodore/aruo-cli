package cli_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/aruodore/aruo-cli/internal/buildinfo"
	catalogbuiltin "github.com/aruodore/aruo-cli/internal/catalog/builtin"
	"github.com/aruodore/aruo-cli/internal/cli"
	"github.com/aruodore/aruo-cli/internal/cli/iostreams"
	"github.com/aruodore/aruo-cli/internal/create"
	"github.com/aruodore/aruo-cli/internal/doctor"
)

// These qualification tests exercise the full capability-detection ->
// policy-resolution -> adapter-selection pipeline for the scenarios the
// terminal UX specification calls out by name: NO_COLOR, CI, narrow
// terminals, and ANSI-free structured output. Signal cancellation and PTY
// prompt interaction are covered in cmd/aruo (real subprocess signals) and
// internal/tux/charm (Huh form behavior); redirected/dumb-terminal scenarios
// are covered by run_test.go and internal/tux/session.

func newDoctorDependencies(t *testing.T, environment map[string]string, probe fakeProbe) (cli.Dependencies, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	engine, err := doctor.NewEngine(doctor.BuiltinChecks()...)
	if err != nil {
		t.Fatal(err)
	}
	service, err := doctor.NewService(engine)
	if err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	stderr := ttyWriter{Buffer: &bytes.Buffer{}, fd: 2}
	return cli.Dependencies{
		Build:       buildinfo.Info{Version: "test"},
		Doctor:      service,
		Logger:      slog.New(slog.NewTextHandler(io.Discard, nil)),
		Streams:     iostreams.IOStreams{In: strings.NewReader(""), Out: &stdout, ErrOut: stderr},
		Environment: environment,
		Probe:       probe,
	}, &stdout, stderr.Buffer
}

func TestDoctorHumanRespectsNoColor(t *testing.T) {
	t.Parallel()

	dependencies, stdout, _ := newDoctorDependencies(t, map[string]string{
		"TERM":      "xterm-256color",
		"COLORTERM": "truecolor",
		"NO_COLOR":  "1",
	}, fakeProbe{terminals: map[int]bool{2: true}})
	code := cli.Run(context.Background(), []string{"doctor", t.TempDir()}, dependencies)
	if code != 3 {
		t.Fatalf("Run() code = %d, want 3", code)
	}
	if strings.Contains(stdout.String(), "\x1b[") {
		t.Fatalf("stdout contains ANSI escapes with NO_COLOR set: %q", stdout.String())
	}
}

func TestDoctorJSONNeverContainsANSI(t *testing.T) {
	t.Parallel()

	dependencies, stdout, stderr := newDoctorDependencies(t, map[string]string{
		"TERM":      "xterm-256color",
		"COLORTERM": "truecolor",
	}, fakeProbe{terminals: map[int]bool{2: true}})
	code := cli.Run(context.Background(), []string{"doctor", t.TempDir(), "--format", "json"}, dependencies)
	if code != 3 {
		t.Fatalf("Run() code = %d, want 3", code)
	}
	if strings.Contains(stdout.String(), "\x1b[") {
		t.Fatalf("structured stdout contains ANSI escapes: %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty for expected findings", stderr.String())
	}
}

func TestDoctorNarrowTerminalDropsLowerPriorityColumns(t *testing.T) {
	t.Parallel()

	dependencies, stdout, _ := newDoctorDependencies(t, map[string]string{"TERM": "xterm-256color"}, fakeProbe{
		terminals: map[int]bool{2: true},
		width:     12,
	})
	code := cli.Run(context.Background(), []string{"doctor", t.TempDir()}, dependencies)
	if code != 3 {
		t.Fatalf("Run() code = %d, want 3", code)
	}
	if !strings.Contains(stdout.String(), "completeness") {
		t.Fatalf("stdout = %q, want the higher-priority category column retained", stdout.String())
	}
	if strings.Contains(stdout.String(), "20/20") {
		t.Fatalf("stdout = %q, want the lower-priority score column dropped at width 12", stdout.String())
	}
}

func TestCreateDisablesPromptsInCI(t *testing.T) {
	t.Parallel()

	templateCatalog, err := catalogbuiltin.New()
	if err != nil {
		t.Fatal(err)
	}
	creator, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	stdin := ttyReader{Reader: strings.NewReader("must not be read"), fd: 0}
	code := cli.Run(context.Background(), []string{
		"create", "--template", "go-library",
	}, cli.Dependencies{
		Build:       buildinfo.Info{Version: "test"},
		Catalog:     templateCatalog,
		Creator:     creator,
		Logger:      slog.New(slog.NewTextHandler(io.Discard, nil)),
		Streams:     iostreams.IOStreams{In: stdin, Out: &stdout, ErrOut: ttyWriter{Buffer: &stderr, fd: 2}},
		Environment: map[string]string{"CI": "true"},
		Probe:       fakeProbe{terminals: map[int]bool{0: true, 2: true}},
	})
	if code == 0 {
		t.Fatalf("Run() code = 0, want non-zero under CI; stdout = %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Project name is required; provide name-or-path argument") {
		t.Fatalf("stderr = %q, want a fast failure instead of a CI prompt", stderr.String())
	}
}
