package lifecycle

import (
	"context"
	"io"
	"os"
	"testing"
)

// BenchmarkInterruptCancellationLatency measures the delay between an
// interrupt arriving on the signal channel and the manager's context
// observing cancellation, evidence for the spec's 100ms p95 acknowledgement
// budget.
func BenchmarkInterruptCancellationLatency(b *testing.B) {
	for i := 0; i < b.N; i++ {
		signals := make(chan os.Signal, 2)
		manager := newManager(context.Background(), signals, nil, Options{
			Diagnostic: io.Discard,
			Restore:    func() error { return nil },
			ForceExit:  func(int) {},
		})
		signals <- os.Interrupt
		<-manager.Context().Done()
		_ = manager.Close()
	}
}
