package charm

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/aruodore/aruo-cli/internal/tux"
)

func TestProgressModelRendersNestedTasks(t *testing.T) {
	t.Parallel()

	model := newProgressModel(tux.Capabilities{Unicode: true})
	model = updateProgress(t, model, tux.TaskEvent{TaskID: "create", Kind: tux.TaskStarted, Label: "Create repository"})
	model = updateProgress(t, model, tux.TaskEvent{TaskID: "render", ParentID: "create", Kind: tux.TaskAdvanced, Label: "Render files", Current: 5, Total: 10})
	view := model.View().Content
	if !strings.Contains(view, "Create repository") || !strings.Contains(view, "  ") || !strings.Contains(view, "50%") {
		t.Fatalf("view = %q", view)
	}
}

func TestProgressModelStopsAnimatingAfterTerminalState(t *testing.T) {
	t.Parallel()

	model := newProgressModel(tux.Capabilities{})
	model = updateProgress(t, model, tux.TaskEvent{TaskID: "audit", Kind: tux.TaskStarted, Label: "Audit"})
	if !model.hasRunning() {
		t.Fatal("hasRunning() = false after start")
	}
	model = updateProgress(t, model, tux.TaskEvent{TaskID: "audit", Kind: tux.TaskCompleted, Label: "Audit"})
	if model.hasRunning() {
		t.Fatal("hasRunning() = true after completion")
	}
	if got := model.View().Content; got != "done Audit" {
		t.Fatalf("view = %q", got)
	}
}

func TestProgressModelBreaksParentCycles(t *testing.T) {
	t.Parallel()

	model := newProgressModel(tux.Capabilities{})
	model = updateProgress(t, model, tux.TaskEvent{TaskID: "a", ParentID: "b", Kind: tux.TaskStarted, Label: "A"})
	model = updateProgress(t, model, tux.TaskEvent{TaskID: "b", ParentID: "a", Kind: tux.TaskStarted, Label: "B"})
	if got := model.View().Content; !strings.Contains(got, "A") || !strings.Contains(got, "B") {
		t.Fatalf("view = %q", got)
	}
}

func TestProgressBarClampsValues(t *testing.T) {
	t.Parallel()

	if got := progressBar(15, 10, false); got != "[##########] 100%" {
		t.Fatalf("progressBar() = %q", got)
	}
}

func TestProgressEmitAndCloseRoundTrip(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	progress := NewProgress(context.Background(), &output, tux.Capabilities{Width: 80, Height: 24})
	if err := progress.Emit(context.Background(), tux.TaskEvent{TaskID: "task", Kind: tux.TaskStarted, Label: "Working"}); err != nil {
		t.Fatalf("Emit() = %v", err)
	}
	if err := progress.Close(); err != nil {
		t.Fatalf("Close() = %v", err)
	}
}

func TestProgressCloseIsIdempotent(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	progress := NewProgress(context.Background(), &output, tux.Capabilities{})
	if err := progress.Close(); err != nil {
		t.Fatalf("Close() #1 = %v", err)
	}
	if err := progress.Close(); err != nil {
		t.Fatalf("Close() #2 = %v", err)
	}
}

func TestProgressEmitAfterCloseReturnsUnavailable(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	progress := NewProgress(context.Background(), &output, tux.Capabilities{})
	if err := progress.Close(); err != nil {
		t.Fatalf("Close() = %v", err)
	}
	err := progress.Emit(context.Background(), tux.TaskEvent{TaskID: "task", Kind: tux.TaskStarted})
	if !errors.Is(err, tux.ErrUnavailable) {
		t.Fatalf("Emit() after Close() = %v, want ErrUnavailable", err)
	}
}

func TestProgressEmitRespectsCancelledContext(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	progress := NewProgress(context.Background(), &output, tux.Capabilities{})
	t.Cleanup(func() { _ = progress.Close() })

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := progress.Emit(ctx, tux.TaskEvent{TaskID: "task", Kind: tux.TaskStarted}); !errors.Is(err, context.Canceled) {
		t.Fatalf("Emit() with cancelled context = %v, want context.Canceled", err)
	}
}

// TestProgressConcurrentEmitDuringCloseNeverBlocks guards the lifecycle
// RWMutex in Progress: Emit and the start of Close's shutdown must be
// mutually exclusive, or a racing Send can be issued after the Bubble Tea
// event loop has stopped reading, blocking the calling goroutine forever.
func TestProgressConcurrentEmitDuringCloseNeverBlocks(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	progress := NewProgress(context.Background(), &output, tux.Capabilities{})

	const emitters = 50
	var wg sync.WaitGroup
	wg.Add(emitters)
	for i := 0; i < emitters; i++ {
		go func(n int) {
			defer wg.Done()
			_ = progress.Emit(context.Background(), tux.TaskEvent{TaskID: "task", Kind: tux.TaskAdvanced, Current: int64(n), Total: emitters})
		}(i)
	}

	if err := progress.Close(); err != nil {
		t.Fatalf("Close() = %v", err)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Emit calls racing Close() did not return; likely blocked on a stopped event loop")
	}
}

func updateProgress(t *testing.T, model progressModel, event tux.TaskEvent) progressModel {
	t.Helper()
	updated, _ := model.Update(taskEventMsg(event))
	result, ok := updated.(progressModel)
	if !ok {
		t.Fatalf("Update() type = %T", updated)
	}
	return result
}

var _ tea.Model = progressModel{}
