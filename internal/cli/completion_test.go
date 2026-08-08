package cli_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/aruodore/aruo-cli/internal/buildinfo"
	catalogbuiltin "github.com/aruodore/aruo-cli/internal/catalog/builtin"
	"github.com/aruodore/aruo-cli/internal/cli"
	"github.com/aruodore/aruo-cli/internal/cli/iostreams"
	"github.com/aruodore/aruo-cli/internal/create"
)

// unreadableStdin fails the test if completion ever reads from stdin;
// generated completion scripts and completion candidates must never prompt.
type unreadableStdin struct{ t *testing.T }

func (u unreadableStdin) Read([]byte) (int, error) {
	u.t.Helper()
	u.t.Fatal("completion read from stdin")
	return 0, errors.New("unreachable")
}

func TestCompletionScriptsGenerateForEachShell(t *testing.T) {
	t.Parallel()

	for _, shell := range []string{"bash", "zsh", "fish", "powershell"} {
		t.Run(shell, func(t *testing.T) {
			t.Parallel()

			var stdout, stderr bytes.Buffer
			code := cli.Run(context.Background(), []string{"completion", shell}, cli.Dependencies{
				Build:   buildinfo.Info{Version: "test"},
				Logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
				Streams: iostreams.IOStreams{In: unreadableStdin{t}, Out: &stdout, ErrOut: &stderr},
			})
			if code != 0 {
				t.Fatalf("Run() code = %d; stderr = %q", code, stderr.String())
			}
			if stdout.Len() == 0 {
				t.Fatal("completion script is empty")
			}
		})
	}
}

func TestDynamicCompletionListsLocalCatalogWithoutPrompting(t *testing.T) {
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
	code := cli.Run(context.Background(), []string{"__complete", "create", "--template", ""}, cli.Dependencies{
		Build:   buildinfo.Info{Version: "test"},
		Catalog: templateCatalog,
		Creator: creator,
		Logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
		Streams: iostreams.IOStreams{In: unreadableStdin{t}, Out: &stdout, ErrOut: &stderr},
	})
	if code != 0 {
		t.Fatalf("Run() code = %d; stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "go-library") {
		t.Fatalf("stdout = %q, want the local catalog's go-library entry", stdout.String())
	}
}
