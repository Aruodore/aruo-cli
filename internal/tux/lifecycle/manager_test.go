package lifecycle

import (
	"bytes"
	"context"
	"errors"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

func TestFirstInterruptCancelsAndRestores(t *testing.T) {
	t.Parallel()

	signals := make(chan os.Signal, 2)
	var diagnostic bytes.Buffer
	var restores atomic.Int32
	manager := newManager(context.Background(), signals, nil, Options{
		Diagnostic: &diagnostic,
		Restore: func() error {
			restores.Add(1)
			return nil
		},
		ForceExit: func(int) { t.Error("ForceExit called on first interrupt") },
	})
	t.Cleanup(func() { _ = manager.Close() })

	signals <- os.Interrupt
	awaitCancellation(t, manager.Context())
	if !errors.Is(context.Cause(manager.Context()), ErrInterrupted) {
		t.Fatalf("cause = %v", context.Cause(manager.Context()))
	}
	if manager.ExitCode() != 130 {
		t.Fatalf("ExitCode() = %d", manager.ExitCode())
	}
	if restores.Load() != 1 {
		t.Fatalf("restores = %d", restores.Load())
	}
	if diagnostic.String() != "Cancelling... Press Ctrl+C again to exit immediately.\n" {
		t.Fatalf("diagnostic = %q", diagnostic.String())
	}
}

func TestSecondInterruptForcesExit(t *testing.T) {
	t.Parallel()

	signals := make(chan os.Signal, 2)
	forced := make(chan int, 1)
	manager := newManager(context.Background(), signals, nil, Options{
		ForceExit: func(code int) { forced <- code },
	})
	t.Cleanup(func() { _ = manager.Close() })

	signals <- os.Interrupt
	awaitCancellation(t, manager.Context())
	signals <- os.Interrupt
	select {
	case code := <-forced:
		if code != 130 {
			t.Fatalf("forced code = %d", code)
		}
	case <-time.After(time.Second):
		t.Fatal("second interrupt did not force exit")
	}
}

func TestTerminationCancelsWithoutPrompt(t *testing.T) {
	t.Parallel()

	signals := make(chan os.Signal, 1)
	var diagnostic bytes.Buffer
	manager := newManager(context.Background(), signals, nil, Options{Diagnostic: &diagnostic})
	t.Cleanup(func() { _ = manager.Close() })

	signals <- syscall.SIGTERM
	awaitCancellation(t, manager.Context())
	if !errors.Is(context.Cause(manager.Context()), ErrTerminated) {
		t.Fatalf("cause = %v", context.Cause(manager.Context()))
	}
	if manager.ExitCode() != 143 {
		t.Fatalf("ExitCode() = %d", manager.ExitCode())
	}
	if diagnostic.Len() != 0 {
		t.Fatalf("diagnostic = %q", diagnostic.String())
	}
}

func TestRestoreIsIdempotent(t *testing.T) {
	t.Parallel()

	var restores atomic.Int32
	manager := newManager(context.Background(), make(chan os.Signal), nil, Options{
		Restore: func() error {
			restores.Add(1)
			return nil
		},
	})
	if err := manager.Restore(); err != nil {
		t.Fatal(err)
	}
	if err := manager.Close(); err != nil {
		t.Fatal(err)
	}
	if restores.Load() != 1 {
		t.Fatalf("restores = %d", restores.Load())
	}
}

func awaitCancellation(t *testing.T, ctx context.Context) {
	t.Helper()
	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled")
	}
}
