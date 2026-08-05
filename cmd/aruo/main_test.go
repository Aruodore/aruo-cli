//go:build !windows

package main_test

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

var binaryPath string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "aruo-signal-test")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	binaryPath = filepath.Join(dir, "aruo")
	build := exec.Command("go", "build", "-o", binaryPath, ".")
	if out, buildErr := build.CombinedOutput(); buildErr != nil {
		fmt.Fprintf(os.Stderr, "build aruo: %v\n%s", buildErr, out)
		os.RemoveAll(dir)
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

func TestFirstInterruptCancelsCooperatively(t *testing.T) {
	t.Parallel()

	cmd, stderr := startSignalFixture(t, "cooperative")
	waitForReady(t, stderr)

	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatalf("Signal(SIGINT) = %v", err)
	}

	waitForLine(t, stderr, "Cancelling... Press Ctrl+C again to exit immediately.")
	assertExitCode(t, cmd, 130)
}

func TestSecondInterruptForcesExit(t *testing.T) {
	t.Parallel()

	cmd, stderr := startSignalFixture(t, "hang")
	waitForReady(t, stderr)

	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatalf("first Signal(SIGINT) = %v", err)
	}
	waitForLine(t, stderr, "Cancelling... Press Ctrl+C again to exit immediately.")

	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatalf("second Signal(SIGINT) = %v", err)
	}
	assertExitCode(t, cmd, 130)
}

func TestSIGTERMCancelsCooperatively(t *testing.T) {
	t.Parallel()

	cmd, stderr := startSignalFixture(t, "cooperative")
	waitForReady(t, stderr)

	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("Signal(SIGTERM) = %v", err)
	}
	assertExitCode(t, cmd, 143)
}

func startSignalFixture(t *testing.T, mode string) (*exec.Cmd, *bufio.Scanner) {
	t.Helper()
	cmd := exec.Command(binaryPath)
	cmd.Env = append(os.Environ(), "ARUO_TEST_SIGNAL_MODE="+mode)
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("StderrPipe() = %v", err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("Start() = %v", err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	})
	return cmd, bufio.NewScanner(stderrPipe)
}

func waitForReady(t *testing.T, stderr *bufio.Scanner) {
	t.Helper()
	waitForLine(t, stderr, "ready")
}

func waitForLine(t *testing.T, stderr *bufio.Scanner, want string) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for stderr.Scan() {
			if strings.Contains(stderr.Text(), want) {
				return
			}
		}
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for line %q", want)
	}
}

func assertExitCode(t *testing.T, cmd *exec.Cmd, want int) {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		var exitErr *exec.ExitError
		switch {
		case err == nil && want == 0:
		case errors.As(err, &exitErr) && exitErr.ExitCode() == want:
		default:
			t.Fatalf("exit error = %v, want code %d", err, want)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("process did not exit within 5s, want code %d", want)
	}
}
