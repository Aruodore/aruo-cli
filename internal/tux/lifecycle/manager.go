// Package lifecycle owns process cancellation and terminal restoration.
package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

var (
	// ErrInterrupted is the cancellation cause for SIGINT and Ctrl+C.
	ErrInterrupted = errors.New("interrupted")
	// ErrTerminated is the cancellation cause for SIGTERM.
	ErrTerminated = errors.New("terminated")
)

// Options contains process-edge effects owned by the lifecycle manager.
type Options struct {
	Diagnostic io.Writer
	Restore    func() error
	ForceExit  func(int)
}

// Manager propagates signal causes and coordinates terminal restoration.
type Manager struct {
	ctx         context.Context
	cancel      context.CancelCauseFunc
	signals     chan os.Signal
	stop        func()
	closed      chan struct{}
	restore     func() error
	restoreErr  error
	forceExit   func(int)
	diagnostic  io.Writer
	exitCode    atomic.Int32
	closeOnce   sync.Once
	restoreOnce sync.Once
}

// New subscribes to process interrupt and termination signals.
func New(parent context.Context, options Options) *Manager {
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	return newManager(parent, signals, func() { signal.Stop(signals) }, options)
}

// Context returns the invocation context cancelled with a stable cause.
func (m *Manager) Context() context.Context {
	return m.ctx
}

// ExitCode returns the signal-derived process status or zero.
func (m *Manager) ExitCode() int {
	return int(m.exitCode.Load())
}

// Restore returns the terminal to its original state at most once.
func (m *Manager) Restore() error {
	m.restoreOnce.Do(func() {
		if m.restore != nil {
			m.restoreErr = m.restore()
		}
	})
	return m.restoreErr
}

// Close releases signal resources and restores terminal state.
func (m *Manager) Close() error {
	var err error
	m.closeOnce.Do(func() {
		m.stop()
		close(m.closed)
		err = m.Restore()
	})
	return err
}

func newManager(parent context.Context, signals chan os.Signal, stop func(), options Options) *Manager {
	ctx, cancel := context.WithCancelCause(parent)
	manager := &Manager{
		ctx:        ctx,
		cancel:     cancel,
		signals:    signals,
		stop:       stop,
		closed:     make(chan struct{}),
		restore:    options.Restore,
		forceExit:  options.ForceExit,
		diagnostic: options.Diagnostic,
	}
	if manager.stop == nil {
		manager.stop = func() {}
	}
	if manager.forceExit == nil {
		manager.forceExit = os.Exit
	}
	if manager.diagnostic == nil {
		manager.diagnostic = io.Discard
	}
	go manager.watch()
	return manager
}

func (m *Manager) watch() {
	interrupts := 0
	for {
		select {
		case <-m.closed:
			return
		case received := <-m.signals:
			switch received {
			case os.Interrupt:
				interrupts++
				if interrupts == 1 {
					m.exitCode.CompareAndSwap(0, 130)
					_ = m.Restore()
					_, _ = fmt.Fprintln(m.diagnostic, "Cancelling... Press Ctrl+C again to exit immediately.")
					m.cancel(ErrInterrupted)
					continue
				}
				_ = m.Restore()
				m.forceExit(130)
			case syscall.SIGTERM:
				m.exitCode.CompareAndSwap(0, 143)
				_ = m.Restore()
				m.cancel(ErrTerminated)
			}
		}
	}
}
