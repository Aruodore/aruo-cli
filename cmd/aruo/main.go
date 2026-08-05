// Command aruo is the Aruo command-line application.
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aruodore/aruo/internal/buildinfo"
	catalogbuiltin "github.com/aruodore/aruo/internal/catalog/builtin"
	"github.com/aruodore/aruo/internal/cli"
	"github.com/aruodore/aruo/internal/cli/iostreams"
	"github.com/aruodore/aruo/internal/create"
	"github.com/aruodore/aruo/internal/doctor"
	"github.com/aruodore/aruo/internal/tux/lifecycle"
	"github.com/aruodore/aruo/internal/tux/term"
)

// terminalResetSequence shows the cursor, exits the alternate screen, and
// disables bracketed paste. Writing it when nothing engaged those modes is a
// harmless no-op; this is the best-effort restoration a forced second
// interrupt performs before os.Exit bypasses every deferred adapter Close.
const terminalResetSequence = "\x1b[?25h\x1b[?1049l\x1b[?2004l\x1b[0m"

func main() {
	os.Exit(run())
}

func run() int {
	streams := iostreams.System()
	environment := term.SystemEnvironment()

	manager := lifecycle.New(context.Background(), lifecycle.Options{
		Diagnostic: streams.ErrOut,
		Restore:    terminalRestorer(streams, environment),
	})
	defer func() { _ = manager.Close() }()

	if mode := os.Getenv("ARUO_TEST_SIGNAL_MODE"); mode != "" {
		return blockForSignalTest(manager, mode)
	}

	templateCatalog, err := catalogbuiltin.New()
	if err != nil {
		_, _ = os.Stderr.WriteString("Error: initialize template catalog: " + err.Error() + "\n")
		return 1
	}
	creator, err := create.NewService(templateCatalog, create.OSWriter{})
	if err != nil {
		_, _ = os.Stderr.WriteString("Error: initialize create service: " + err.Error() + "\n")
		return 1
	}
	doctorEngine, err := doctor.NewEngine(doctor.BuiltinChecks()...)
	if err != nil {
		_, _ = os.Stderr.WriteString("Error: initialize doctor engine: " + err.Error() + "\n")
		return 1
	}
	doctorService, err := doctor.NewService(doctorEngine)
	if err != nil {
		_, _ = os.Stderr.WriteString("Error: initialize doctor service: " + err.Error() + "\n")
		return 1
	}
	dependencies := cli.Dependencies{
		Build:       buildinfo.Current(),
		Catalog:     templateCatalog,
		Creator:     creator,
		Doctor:      doctorService,
		Logger:      slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Streams:     streams,
		Environment: environment,
	}

	code := cli.Run(manager.Context(), os.Args[1:], dependencies)
	if signalCode := manager.ExitCode(); signalCode != 0 {
		return signalCode
	}
	return code
}

// blockForSignalTest exercises the compiled binary's real signal handling in
// subprocess fixtures, which cannot otherwise observe os.Exit or actual OS
// signal delivery. "cooperative" returns as soon as the root context is
// cancelled, matching well-behaved commands. "hang" ignores cancellation, the
// scenario the second forced interrupt exists to terminate.
func blockForSignalTest(manager *lifecycle.Manager, mode string) int {
	_, _ = os.Stderr.WriteString("ready\n")
	switch mode {
	case "cooperative":
		<-manager.Context().Done()
		return manager.ExitCode()
	case "hang":
		select {}
	default:
		_, _ = os.Stderr.WriteString("Error: unknown ARUO_TEST_SIGNAL_MODE " + mode + "\n")
		return 2
	}
}

// terminalRestorer resolves whether stderr is a terminal once and returns a
// restoration hook the lifecycle manager can call at every cancellation
// point, including the second interrupt that bypasses deferred cleanup.
func terminalRestorer(streams iostreams.IOStreams, environment map[string]string) func() error {
	capabilities := term.NewDetector().Detect(streams.In, streams.Out, streams.ErrOut, environment)
	if !capabilities.ErrorTTY {
		return func() error { return nil }
	}
	return func() error {
		_, err := streams.ErrOut.Write([]byte(terminalResetSequence))
		return err
	}
}
