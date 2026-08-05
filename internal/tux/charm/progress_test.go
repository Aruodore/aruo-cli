package charm

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/aruodore/aruo/internal/tux"
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
