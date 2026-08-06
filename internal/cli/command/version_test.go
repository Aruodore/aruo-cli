package command

import (
	"bytes"
	"testing"

	"github.com/aruodore/aruo/internal/buildinfo"
)

func TestVersionCommandPrintsVersion(t *testing.T) {
	t.Parallel()

	command := newVersion(buildinfo.Info{Version: "1.2.3"})
	var out bytes.Buffer
	command.SetOut(&out)
	command.SetArgs(nil)
	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	const want = "aruo version 1.2.3\n"
	if got := out.String(); got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestVersionCommandRejectsArgs(t *testing.T) {
	t.Parallel()

	command := newVersion(buildinfo.Info{Version: "dev"})
	command.SetArgs([]string{"unexpected"})
	command.SilenceErrors = true
	command.SilenceUsage = true
	if err := command.Execute(); err == nil {
		t.Fatal("Execute() error = nil, want error for unexpected positional argument")
	}
}
