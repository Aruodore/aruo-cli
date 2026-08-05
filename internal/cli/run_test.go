package cli_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/aruodore/aruo/internal/buildinfo"
	catalogbuiltin "github.com/aruodore/aruo/internal/catalog/builtin"
	"github.com/aruodore/aruo/internal/cli"
	"github.com/aruodore/aruo/internal/cli/iostreams"
	"github.com/aruodore/aruo/internal/create"
	"github.com/aruodore/aruo/internal/doctor"
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

type fakeProbe struct{ terminals map[int]bool }

func (p fakeProbe) IsTerminal(fd int) bool     { return p.terminals[fd] }
func (p fakeProbe) Size(int) (int, int, error) { return 80, 24, nil }

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
		"--template", "go-library",
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
	if !strings.Contains(stderr.String(), "Go module path is required; provide --module") {
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
		"example.com/guided",
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
	if !strings.Contains(stderr.String(), "Go module path") ||
		!strings.Contains(stderr.String(), "Example: github.com/your-name/my-library") ||
		!strings.Contains(stderr.String(), "Project summary:") ||
		!strings.Contains(stderr.String(), "Create this project?") ||
		!strings.Contains(stdout.String(), "Created go-library") {
		t.Errorf("stdout=%q stderr=%q", stdout.String(), stderr.String())
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
