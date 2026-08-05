// Command aruo is the Aruo command-line application.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/aruodore/aruo/internal/buildinfo"
	catalogbuiltin "github.com/aruodore/aruo/internal/catalog/builtin"
	"github.com/aruodore/aruo/internal/cli"
	"github.com/aruodore/aruo/internal/cli/iostreams"
	"github.com/aruodore/aruo/internal/create"
	"github.com/aruodore/aruo/internal/doctor"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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
		Build:   buildinfo.Current(),
		Catalog: templateCatalog,
		Creator: creator,
		Doctor:  doctorService,
		Logger:  slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Streams: iostreams.System(),
	}

	return cli.Run(ctx, os.Args[1:], dependencies)
}
