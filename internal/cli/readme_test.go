package cli_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/aruodore/aruo-cli/internal/buildinfo"
	catalogbuiltin "github.com/aruodore/aruo-cli/internal/catalog/builtin"
	"github.com/aruodore/aruo-cli/internal/cli"
	"github.com/aruodore/aruo-cli/internal/cli/iostreams"
	"github.com/aruodore/aruo-cli/internal/create"
	"github.com/aruodore/aruo-cli/internal/doctor"
)

// These tests pin down the exact claims README.md makes about command
// output, generated file structure, flags, and exit codes, so an
// unintentional behavior change shows up as a failing test here rather
// than a silently stale README.

func newReadmeCreateDependencies(t *testing.T) cli.Dependencies {
	t.Helper()
	templateCatalog, err := catalogbuiltin.New()
	if err != nil {
		t.Fatal(err)
	}
	creator, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		t.Fatal(err)
	}
	return cli.Dependencies{
		Build:   buildinfo.Info{Version: "test"},
		Catalog: templateCatalog,
		Creator: creator,
		Logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func newReadmeDoctorService(t *testing.T) *doctor.Service {
	t.Helper()
	engine, err := doctor.NewEngine(doctor.BuiltinChecks()...)
	if err != nil {
		t.Fatal(err)
	}
	service, err := doctor.NewService(engine)
	if err != nil {
		t.Fatal(err)
	}
	return service
}

// TestReadmeQuickStartWorkflowMatchesDocumentedOutput runs the exact
// commands shown in README.md's "Example workflow" section and checks the
// output against what's quoted there verbatim.
func TestReadmeQuickStartWorkflowMatchesDocumentedOutput(t *testing.T) {
	t.Parallel()

	deps := newReadmeCreateDependencies(t)
	deps.Doctor = newReadmeDoctorService(t)

	destination := filepath.Join(t.TempDir(), "my-library")
	var createOut, createErr bytes.Buffer
	deps.Streams = iostreams.IOStreams{In: strings.NewReader(""), Out: &createOut, ErrOut: &createErr}
	code := cli.Run(context.Background(), []string{
		"create", destination,
		"--template", "go-library",
		"--no-input",
	}, deps)
	if code != 0 {
		t.Fatalf("create Run() code = %d; stderr = %q", code, createErr.String())
	}
	if !strings.Contains(createOut.String(), "Created go-library with 24 files") {
		t.Errorf("create stdout = %q, want the README's documented file count", createOut.String())
	}
	// README's Example workflow omits --module entirely, relying on it
	// defaulting to the project's own name ("my-library" here).
	goMod, err := os.ReadFile(filepath.Join(destination, "go.mod"))
	if err != nil {
		t.Fatalf("read generated go.mod: %v", err)
	}
	if !strings.Contains(string(goMod), "module my-library\n") {
		t.Errorf("go.mod = %q, want the README's documented \"ok  my-library\" go test line to match a real module default", goMod)
	}

	var doctorOut, doctorErr bytes.Buffer
	deps.Streams = iostreams.IOStreams{In: strings.NewReader(""), Out: &doctorOut, ErrOut: &doctorErr}
	code = cli.Run(context.Background(), []string{"doctor", destination}, deps)
	if code != 0 {
		t.Fatalf("doctor Run() code = %d; stderr = %q", code, doctorErr.String())
	}
	for _, want := range []string{
		"Repository health: 99/100 (A)",
		"completeness   20/20",
		"documentation  20/20",
		"ci             15/15",
		"tests          15/15",
		"license        10/10",
		"security       10/10",
		"github          9/10",
		"Missing CODEOWNERS",
	} {
		if !strings.Contains(doctorOut.String(), want) {
			t.Errorf("doctor stdout missing %q documented in README's Example workflow; stdout = %q", want, doctorOut.String())
		}
	}
}

// TestReadmeGeneratedProjectStructureMatchesDocumentedTree asserts the
// go-library file set matches the tree documented in README.md's
// "Generated project structure" section exactly.
func TestReadmeGeneratedProjectStructureMatchesDocumentedTree(t *testing.T) {
	t.Parallel()

	deps := newReadmeCreateDependencies(t)
	destination := filepath.Join(t.TempDir(), "my-library")
	var stdout, stderr bytes.Buffer
	deps.Streams = iostreams.IOStreams{In: strings.NewReader(""), Out: &stdout, ErrOut: &stderr}
	code := cli.Run(context.Background(), []string{
		"create", destination,
		"--template", "go-library",
		"--no-input",
	}, deps)
	if code != 0 {
		t.Fatalf("Run() code = %d; stderr = %q", code, stderr.String())
	}

	want := []string{
		".github/ISSUE_TEMPLATE/bug.yml",
		".github/ISSUE_TEMPLATE/feature.yml",
		".github/workflows/ci.yml",
		".github/workflows/pr-title.yml",
		".github/workflows/release.yml",
		".github/dependabot.yml",
		".github/pull_request_template.md",
		"docs/README.md",
		"aruo.yaml",
		"CHANGELOG.md",
		"CODE_OF_CONDUCT.md",
		"CONTRIBUTING.md",
		"go.mod",
		"LICENSE",
		"Makefile",
		"my_library.go",
		"my_library_test.go",
		"README.md",
		"release-please-config.json",
		".release-please-manifest.json",
		"ROADMAP.md",
		"SECURITY.md",
		".editorconfig",
		".gitignore",
	}
	sort.Strings(want)

	var got []string
	if err := filepath.WalkDir(destination, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(destination, path)
		if err != nil {
			return err
		}
		got = append(got, filepath.ToSlash(relative))
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	sort.Strings(got)

	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("generated file set does not match README.md's documented tree.\ngot:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

// TestReadmeVersionOutputFormat pins the exact "aruo version" wording the
// README quotes for a locally built (unreleased) binary.
func TestReadmeVersionOutputFormat(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	code := cli.Run(context.Background(), []string{"version"}, cli.Dependencies{
		Build:  buildinfo.Info{Version: "dev"},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Streams: iostreams.IOStreams{
			In: strings.NewReader(""), Out: &stdout, ErrOut: &stderr,
		},
	})
	if code != 0 {
		t.Fatalf("Run() code = %d; stderr = %q", code, stderr.String())
	}
	if stdout.String() != "aruo version dev\n" {
		t.Errorf("stdout = %q, want %q", stdout.String(), "aruo version dev\n")
	}
}

// TestReadmeCreateDotExampleWorks verifies the "aruo create . --name ..."
// example in README.md's Commands section: the destination "." always
// already exists, so this only works because Write adopts an existing
// empty directory instead of refusing any existing path outright.
func TestReadmeCreateDotExampleWorks(t *testing.T) {
	t.Parallel()

	deps := newReadmeCreateDependencies(t)
	destination := t.TempDir()
	var stdout, stderr bytes.Buffer
	deps.Streams = iostreams.IOStreams{In: strings.NewReader(""), Out: &stdout, ErrOut: &stderr}
	code := cli.Run(context.Background(), []string{
		"create", destination,
		"--name", "my-tool",
		"--template", "go-library",
		"--module", "github.com/you/my-tool",
		"--no-input",
	}, deps)
	if code != 0 {
		t.Fatalf("Run() code = %d; stderr = %q", code, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(destination, "go.mod")); err != nil {
		t.Errorf("go.mod stat error = %v, want create . to have populated the existing empty directory", err)
	}
}

// TestReadmeDocumentedCreateFlagsExist checks every create flag named in
// README.md's Commands section is a real, currently registered flag.
func TestReadmeDocumentedCreateFlagsExist(t *testing.T) {
	t.Parallel()

	deps := newReadmeCreateDependencies(t)
	var stdout, stderr bytes.Buffer
	deps.Streams = iostreams.IOStreams{In: strings.NewReader(""), Out: &stdout, ErrOut: &stderr}
	code := cli.Run(context.Background(), []string{"create", "--help"}, deps)
	if code != 0 {
		t.Fatalf("Run() code = %d; stderr = %q", code, stderr.String())
	}
	for _, flag := range []string{
		"--template", "--language", "--kind", "--module", "--description",
		"--author", "--license", "--set", "-y, --yes", "--no-input",
	} {
		if !strings.Contains(stdout.String(), flag) {
			t.Errorf("create --help missing documented flag %q; help = %q", flag, stdout.String())
		}
	}
}

// TestReadmeDocumentedDoctorFlagsExist checks every doctor flag named in
// README.md's Commands section is a real, currently registered flag.
func TestReadmeDocumentedDoctorFlagsExist(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	code := cli.Run(context.Background(), []string{"doctor", "--help"}, cli.Dependencies{
		Build:  buildinfo.Info{Version: "test"},
		Doctor: newReadmeDoctorService(t),
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Streams: iostreams.IOStreams{
			In: strings.NewReader(""), Out: &stdout, ErrOut: &stderr,
		},
	})
	if code != 0 {
		t.Fatalf("Run() code = %d; stderr = %q", code, stderr.String())
	}
	for _, flag := range []string{"--format", "--minimum-score"} {
		if !strings.Contains(stdout.String(), flag) {
			t.Errorf("doctor --help missing documented flag %q; help = %q", flag, stdout.String())
		}
	}
}

// TestReadmeDocumentedCommandsExist checks every subcommand named in
// README.md's Commands section is registered on the root command.
func TestReadmeDocumentedCommandsExist(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	deps := newReadmeCreateDependencies(t)
	deps.Doctor = newReadmeDoctorService(t)
	deps.Streams = iostreams.IOStreams{In: strings.NewReader(""), Out: &stdout, ErrOut: &stderr}
	code := cli.Run(context.Background(), []string{"--help"}, deps)
	if code != 0 {
		t.Fatalf("Run() code = %d; stderr = %q", code, stderr.String())
	}
	for _, command := range []string{"create", "doctor", "version", "completion", "help"} {
		if !strings.Contains(stdout.String(), command) {
			t.Errorf("aruo --help missing documented command %q; help = %q", command, stdout.String())
		}
	}
}
