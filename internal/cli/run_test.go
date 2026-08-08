package cli_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aruodore/aruo-cli/internal/buildinfo"
	catalogbuiltin "github.com/aruodore/aruo-cli/internal/catalog/builtin"
	"github.com/aruodore/aruo-cli/internal/cli"
	"github.com/aruodore/aruo-cli/internal/cli/iostreams"
	"github.com/aruodore/aruo-cli/internal/create"
	"github.com/aruodore/aruo-cli/internal/doctor"
)

// ttyReader and ttyWriter fake a real terminal descriptor so tests can
// exercise interactive-eligible sessions without an actual PTY; a paired
// fakeProbe reports which faked file descriptors are terminals.
type ttyReader struct {
	*strings.Reader
	fd uintptr
}

func (r ttyReader) Fd() uintptr { return r.fd }

type ttyWriter struct {
	*bytes.Buffer
	fd uintptr
}

func (w ttyWriter) Fd() uintptr { return w.fd }

// width defaults to 80 when unset so existing fixtures keep their prior
// behavior; qualification tests override it to exercise narrow rendering.
type fakeProbe struct {
	terminals map[int]bool
	width     int
}

func (p fakeProbe) IsTerminal(fd int) bool { return p.terminals[fd] }
func (p fakeProbe) Size(int) (int, int, error) {
	width := p.width
	if width == 0 {
		width = 80
	}
	return width, 24, nil
}

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		args        []string
		wantCode    int
		wantOut     string
		wantErrPart string
	}{
		{name: "root shows help", wantCode: 0, wantOut: "Usage:\n  aruo [flags]"},
		{name: "help command", args: []string{"help"}, wantCode: 0, wantOut: "Available Commands:"},
		{name: "version command", args: []string{"version"}, wantCode: 0, wantOut: "aruo version 1.2.3\n"},
		{name: "unknown command", args: []string{"unknown"}, wantCode: 1, wantErrPart: "unknown command"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var stdout bytes.Buffer
			var stderr bytes.Buffer
			dependencies := cli.Dependencies{
				Build:  buildinfo.Info{Version: "1.2.3", Commit: "abc123", Date: "2026-08-04"},
				Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
				Streams: iostreams.IOStreams{
					In:     strings.NewReader(""),
					Out:    &stdout,
					ErrOut: &stderr,
				},
			}

			code := cli.Run(context.Background(), test.args, dependencies)
			if code != test.wantCode {
				t.Fatalf("Run() code = %d, want %d", code, test.wantCode)
			}
			if !strings.Contains(stdout.String(), test.wantOut) {
				t.Errorf("stdout = %q, want substring %q", stdout.String(), test.wantOut)
			}
			if !strings.Contains(stderr.String(), test.wantErrPart) {
				t.Errorf("stderr = %q, want substring %q", stderr.String(), test.wantErrPart)
			}
		})
	}
}

func TestRunCreateNonInteractive(t *testing.T) {
	t.Parallel()

	templateCatalog, err := catalogbuiltin.New()
	if err != nil {
		t.Fatal(err)
	}
	creator, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	destination := t.TempDir() + "/created"
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := cli.Run(context.Background(), []string{
		"create", destination,
		"--template", "go-library",
		"--module", "example.com/created",
		"--description", "A generated library.",
		"--author", "Example Authors",
		"--non-interactive",
	}, cli.Dependencies{
		Build:   buildinfo.Info{Version: "test"},
		Catalog: templateCatalog,
		Creator: creator,
		Logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
		Streams: iostreams.IOStreams{In: strings.NewReader("must not be read"), Out: &stdout, ErrOut: &stderr},
	})
	if code != 0 {
		t.Fatalf("Run() code = %d; stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Created go-library with 24 files") {
		t.Errorf("stdout = %q", stdout.String())
	}
}

func TestRunCreateNoInputSucceedsWithoutOptionalFields(t *testing.T) {
	t.Parallel()

	templateCatalog, err := catalogbuiltin.New()
	if err != nil {
		t.Fatal(err)
	}
	creator, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	destination := t.TempDir() + "/created"
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := cli.Run(context.Background(), []string{
		"create", destination,
		"--template", "go-library",
		"--module", "example.com/created",
		"--no-input",
	}, cli.Dependencies{
		Build:   buildinfo.Info{Version: "test"},
		Catalog: templateCatalog,
		Creator: creator,
		Logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
		Streams: iostreams.IOStreams{In: strings.NewReader("must not be read"), Out: &stdout, ErrOut: &stderr},
	})
	if code != 0 {
		t.Fatalf("Run() code = %d, want 0 since --description/--author are optional; stderr = %q", code, stderr.String())
	}
	// Description defaults to the catalog entry's own description, not a
	// placeholder that needs manual follow-up.
	readme, err := os.ReadFile(destination + "/README.md")
	if err != nil {
		t.Fatalf("read generated README.md: %v", err)
	}
	wantDescription := "A production-ready Go library with zero external dependencies, tests, CI, governance, security, and documentation"
	if !strings.Contains(string(readme), wantDescription) {
		t.Errorf("README.md = %q, want the go-library entry's own description", readme)
	}
	// Author defaults to `git config user.name` when available; either way,
	// no visible TODO/placeholder text should land in the generated LICENSE.
	// Asserting a specific name here would make this test depend on the
	// test runner's git configuration, which CI environments don't reliably set.
	license, err := os.ReadFile(destination + "/LICENSE")
	if err != nil {
		t.Fatalf("read generated LICENSE: %v", err)
	}
	if strings.Contains(string(license), "TODO") {
		t.Errorf("LICENSE = %q, want no placeholder TODO text", license)
	}
}

func TestRunCreateNoInputDefaultsModuleToProjectName(t *testing.T) {
	t.Parallel()

	templateCatalog, err := catalogbuiltin.New()
	if err != nil {
		t.Fatal(err)
	}
	creator, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	destination := t.TempDir() + "/my-library"
	var stdout, stderr bytes.Buffer
	code := cli.Run(context.Background(), []string{
		"create", destination,
		"--template", "go-library",
		"--no-input",
	}, cli.Dependencies{
		Build:   buildinfo.Info{Version: "test"},
		Catalog: templateCatalog,
		Creator: creator,
		Logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
		Streams: iostreams.IOStreams{In: strings.NewReader("must not be read"), Out: &stdout, ErrOut: &stderr},
	})
	if code != 0 {
		t.Fatalf("Run() code = %d, want 0 since --module now defaults to the project name; stderr = %q", code, stderr.String())
	}
	goMod, err := os.ReadFile(destination + "/go.mod")
	if err != nil {
		t.Fatalf("read generated go.mod: %v", err)
	}
	if !strings.Contains(string(goMod), "module my-library\n") {
		t.Errorf("go.mod = %q, want the module to default to the destination's base name", goMod)
	}
}

func TestRunCreateNoInputFailsOnMissingValue(t *testing.T) {
	t.Parallel()

	templateCatalog, err := catalogbuiltin.New()
	if err != nil {
		t.Fatal(err)
	}
	creator, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	destination := t.TempDir() + "/blocked"
	var stdout, stderr bytes.Buffer
	stdin := ttyReader{Reader: strings.NewReader("must not be read"), fd: 0}
	code := cli.Run(context.Background(), []string{
		"create", destination,
		// Ambiguous on purpose: --kind/--language match both ts-library
		// and vue-library, and no --template picks between them.
		"--kind", "library", "--language", "typescript",
		"--no-input",
	}, cli.Dependencies{
		Build:   buildinfo.Info{Version: "test"},
		Catalog: templateCatalog,
		Creator: creator,
		Logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
		Streams: iostreams.IOStreams{In: stdin, Out: &stdout, ErrOut: ttyWriter{Buffer: &stderr, fd: 2}},
		Probe:   fakeProbe{terminals: map[int]bool{0: true, 2: true}},
	})
	if code == 0 {
		t.Fatalf("Run() code = 0, want non-zero; stdout = %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "--template is required when multiple templates match") {
		t.Errorf("stderr = %q", stderr.String())
	}
}

func TestRunCreateInteractive(t *testing.T) {
	t.Parallel()

	templateCatalog, err := catalogbuiltin.New()
	if err != nil {
		t.Fatal(err)
	}
	creator, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	destination := t.TempDir() + "/guided"
	input := strings.Join([]string{
		"2", // Kind: Library. Options are Application then Library.
		"1", // Template: Go library. Within the library kind, the catalog
		// sorts by ID: go-library, js-library, python-library,
		// ts-library, vue-library, putting go-library at position 1.
		"back", // On the description screen (the module screen was
		// removed entirely -- it always defaults to the project name
		// now): go back to the template screen instead, proving
		// Guide's back-navigation works end to end through create's
		// real wiring, not just in adapter-level tests.
		"", // Bare Enter on the revisited template screen keeps the
		// prior answer (go-library) as its default.
		"A guided library.",
		"Guided Authors",
		"yes",
	}, "\n") + "\n"
	var stdout bytes.Buffer
	stdin := ttyReader{Reader: strings.NewReader(input), fd: 0}
	stderr := ttyWriter{Buffer: &bytes.Buffer{}, fd: 2}
	code := cli.Run(context.Background(), []string{"create", destination, "--accessible"}, cli.Dependencies{
		Build:   buildinfo.Info{Version: "test"},
		Catalog: templateCatalog,
		Creator: creator,
		Logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
		Streams: iostreams.IOStreams{In: stdin, Out: &stdout, ErrOut: stderr},
		Probe:   fakeProbe{terminals: map[int]bool{0: true, 2: true}},
	})
	if code != 0 {
		t.Fatalf("Run() code = %d; stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "What are you building?") ||
		!strings.Contains(stderr.String(), "1. Application") ||
		!strings.Contains(stderr.String(), "2. Library") ||
		!strings.Contains(stderr.String(), "Go module path: guided") ||
		!strings.Contains(stderr.String(), "Project summary:") ||
		!strings.Contains(stderr.String(), "Create this project?") ||
		!strings.Contains(stdout.String(), "Created go-library") {
		t.Errorf("stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	if got := strings.Count(stderr.String(), "Template\n"); got != 2 {
		t.Errorf("stderr shows the Template screen %d times, want 2 (once forward, once after going back): %q", got, stderr.String())
	}
	if !strings.Contains(stderr.String(), "1. Go library (default)") {
		t.Errorf("stderr = %q, want the revisited template screen to show the prior answer as its default", stderr.String())
	}
	if strings.Contains(stderr.String(), "npm package name") || strings.Contains(stderr.String(), "Go module path\n") {
		t.Errorf("stderr = %q, want no module prompt screen at all -- it always defaults to the project name now", stderr.String())
	}
}

// TestRunCreateInteractiveDefaultsModuleToProjectName proves --module is
// never prompted for at all, even interactively: it silently becomes the
// project name the user already typed.
func TestRunCreateInteractiveDefaultsModuleToProjectName(t *testing.T) {
	// Not t.Parallel(): uses t.Chdir, which the testing package forbids
	// combining with parallel subtests.
	templateCatalog, err := catalogbuiltin.New()
	if err != nil {
		t.Fatal(err)
	}
	creator, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	input := strings.Join([]string{
		"my-guided-tool", // Project name (no destination positional arg).
		"2",              // Kind: Library.
		"1",              // Template: Go library.
		"",               // Description: bare Enter, keeps the catalog default.
		"",               // Author: bare Enter, keeps the git-detected default.
		"yes",
	}, "\n") + "\n"
	var stdout bytes.Buffer
	stdin := ttyReader{Reader: strings.NewReader(input), fd: 0}
	stderr := ttyWriter{Buffer: &bytes.Buffer{}, fd: 2}
	workdir := t.TempDir()
	t.Chdir(workdir)
	code := cli.Run(context.Background(), []string{"create", "--accessible"}, cli.Dependencies{
		Build:   buildinfo.Info{Version: "test"},
		Catalog: templateCatalog,
		Creator: creator,
		Logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
		Streams: iostreams.IOStreams{In: stdin, Out: &stdout, ErrOut: stderr},
		Probe:   fakeProbe{terminals: map[int]bool{0: true, 2: true}},
	})
	if code != 0 {
		t.Fatalf("Run() code = %d; stderr = %q", code, stderr.String())
	}
	if strings.Contains(stderr.String(), "npm package name") || strings.Contains(stderr.String(), "Go module path\n") {
		t.Errorf("stderr = %q, want no module prompt screen at all", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Go module path: my-guided-tool") {
		t.Errorf("stderr = %q, want the confirm summary to show the module defaulted to the project name", stderr.String())
	}
	goMod, err := os.ReadFile(filepath.Join(workdir, "my-guided-tool", "go.mod"))
	if err != nil {
		t.Fatalf("read generated go.mod: %v", err)
	}
	if !strings.Contains(string(goMod), "module my-guided-tool\n") {
		t.Errorf("go.mod = %q, want the module to default to the project name", goMod)
	}
}

func TestRunDoctorJSONAndFindingsExitCode(t *testing.T) {
	t.Parallel()
	engine, err := doctor.NewEngine(doctor.BuiltinChecks()...)
	if err != nil {
		t.Fatal(err)
	}
	service, err := doctor.NewService(engine)
	if err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := cli.Run(context.Background(), []string{"doctor", t.TempDir(), "--format", "json"}, cli.Dependencies{
		Build:  buildinfo.Info{Version: "test"},
		Doctor: service,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Streams: iostreams.IOStreams{
			In: strings.NewReader(""), Out: &stdout, ErrOut: &stderr,
		},
	})
	if code != 3 {
		t.Fatalf("Run() code = %d, want 3; stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"policy": "aruo.repository-health/v1"`) || !strings.Contains(stdout.String(), `"score": 0`) {
		t.Errorf("stdout = %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("stderr = %q, want empty for expected findings", stderr.String())
	}
}

func TestRunDoctorHumanFormat(t *testing.T) {
	t.Parallel()
	engine, err := doctor.NewEngine(doctor.BuiltinChecks()...)
	if err != nil {
		t.Fatal(err)
	}
	service, err := doctor.NewService(engine)
	if err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := cli.Run(context.Background(), []string{"doctor", t.TempDir()}, cli.Dependencies{
		Build:  buildinfo.Info{Version: "test"},
		Doctor: service,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Streams: iostreams.IOStreams{
			In: strings.NewReader(""), Out: &stdout, ErrOut: &stderr,
		},
	})
	if code != 3 {
		t.Fatalf("Run() code = %d, want 3; stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Repository health:") || !strings.Contains(stdout.String(), "Recommendations:") {
		t.Errorf("stdout = %q", stdout.String())
	}
}
