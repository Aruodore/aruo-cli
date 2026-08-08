package cli_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/aruodore/aruo-cli/internal/buildinfo"
	"github.com/aruodore/aruo-cli/internal/cli"
	"github.com/aruodore/aruo-cli/internal/cli/iostreams"
)

// BenchmarkRunVersion measures command-tree construction and execution
// overhead for the fast, dependency-free path. It excludes process startup
// (exec, dynamic linking, OS scheduling), which the release benchmark suite
// measures against the compiled binary; this isolates Aruo's own overhead.
func BenchmarkRunVersion(b *testing.B) {
	dependencies := cli.Dependencies{
		Build:  buildinfo.Info{Version: "bench"},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	for i := 0; i < b.N; i++ {
		var stdout, stderr bytes.Buffer
		dependencies.Streams = iostreams.IOStreams{In: strings.NewReader(""), Out: &stdout, ErrOut: &stderr}
		if code := cli.Run(context.Background(), []string{"version"}, dependencies); code != 0 {
			b.Fatalf("Run() code = %d", code)
		}
	}
}
