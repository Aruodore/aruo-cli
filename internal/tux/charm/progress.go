package charm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/aruodore/aruo-cli/internal/tux"
)

// Progress renders task events in one Bubble Tea-owned terminal region.
type Progress struct {
	program   *tea.Program
	done      chan error
	closeOnce sync.Once
	closeErr  error

	// lifecycle guards Emit against Close: Emit holds it for read for the
	// duration of one Send, and Close takes it for write before starting
	// shutdown. This guarantees Send is never called concurrently with or
	// after Quit, so it can never block on a program whose event loop has
	// already stopped reading.
	lifecycle sync.RWMutex
	closed    bool
}

// NewProgress starts a live progress renderer with Aruo-owned lifecycle and signals.
func NewProgress(ctx context.Context, output io.Writer, capabilities tux.Capabilities) *Progress {
	model := newProgressModel(capabilities)
	program := tea.NewProgram(
		model,
		tea.WithContext(ctx),
		tea.WithInput(nil),
		tea.WithOutput(output),
		tea.WithFPS(10),
		tea.WithWindowSize(capabilities.Width, capabilities.Height),
		tea.WithoutSignalHandler(),
	)
	progress := &Progress{program: program, done: make(chan error, 1)}
	go func() {
		_, err := program.Run()
		progress.done <- err
	}()
	return progress
}

// Emit sends one task transition to the live renderer.
func (p *Progress) Emit(ctx context.Context, event tux.TaskEvent) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	p.lifecycle.RLock()
	defer p.lifecycle.RUnlock()
	if p.closed {
		return tux.ErrUnavailable
	}
	p.program.Send(taskEventMsg(event))
	return nil
}

// Close completes the live region and restores Bubble Tea terminal state.
func (p *Progress) Close() error {
	p.closeOnce.Do(func() {
		p.lifecycle.Lock()
		p.closed = true
		p.lifecycle.Unlock()
		p.program.Quit()
		p.closeErr = <-p.done
		if errors.Is(p.closeErr, tea.ErrProgramKilled) {
			p.closeErr = nil
		}
	})
	return p.closeErr
}

type taskEventMsg tux.TaskEvent

type taskState struct {
	id       string
	parentID string
	kind     tux.TaskEventKind
	label    string
	current  int64
	total    int64
}

type progressModel struct {
	spinner spinner.Model
	tasks   map[string]taskState
	order   []string
	unicode bool
	theme   theme
}

func newProgressModel(capabilities tux.Capabilities) progressModel {
	indicator := spinner.Line
	if capabilities.Unicode {
		indicator = spinner.MiniDot
	}
	visual := newTheme(capabilities, tux.Policy{})
	model := spinner.New(spinner.WithSpinner(indicator))
	model.Style = visual.accent
	return progressModel{
		spinner: model,
		tasks:   make(map[string]taskState),
		unicode: capabilities.Unicode,
		theme:   visual,
	}
}

func (m progressModel) Init() tea.Cmd {
	return nil
}

func (m progressModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case taskEventMsg:
		event := tux.TaskEvent(message)
		wasRunning := m.hasRunning()
		state, exists := m.tasks[event.TaskID]
		if !exists {
			m.order = append(m.order, event.TaskID)
			state.id = event.TaskID
		}
		state.parentID = event.ParentID
		state.kind = event.Kind
		state.label = event.Label
		state.current = event.Current
		state.total = event.Total
		m.tasks[event.TaskID] = state
		if !wasRunning && m.hasRunning() {
			return m, m.spinner.Tick
		}
	case spinner.TickMsg:
		if !m.hasRunning() {
			return m, nil
		}
		var command tea.Cmd
		m.spinner, command = m.spinner.Update(message)
		return m, command
	}
	return m, nil
}

func (m progressModel) View() tea.View {
	lines := make([]string, 0, len(m.order))
	for _, id := range m.order {
		task := m.tasks[id]
		depth := m.depth(task)
		lines = append(lines, strings.Repeat("  ", depth)+m.taskLine(task))
	}
	return tea.NewView(strings.Join(lines, "\n"))
}

func (m progressModel) hasRunning() bool {
	for _, task := range m.tasks {
		switch task.kind {
		case tux.TaskStarted, tux.TaskAdvanced, tux.TaskMessage, tux.TaskRetrying:
			return true
		default:
		}
	}
	return false
}

func (m progressModel) depth(task taskState) int {
	depth := 0
	seen := map[string]struct{}{task.id: {}}
	for parent := task.parentID; parent != ""; {
		if _, cycle := seen[parent]; cycle {
			break
		}
		seen[parent] = struct{}{}
		depth++
		ancestor, found := m.tasks[parent]
		if !found {
			break
		}
		parent = ancestor.parentID
	}
	return depth
}

func (m progressModel) taskLine(task taskState) string {
	indicator := m.spinner.View()
	switch task.kind {
	case tux.TaskCompleted:
		indicator = m.theme.success.Render("done")
		if m.unicode {
			indicator = m.theme.success.Render("✓")
		}
	case tux.TaskFailed:
		indicator = m.theme.danger.Render("failed")
		if m.unicode {
			indicator = m.theme.danger.Render("✗")
		}
	case tux.TaskCancelled:
		indicator = m.theme.muted.Render("cancelled")
	case tux.TaskRetrying:
		indicator = m.theme.warning.Render("retry")
	default:
	}
	line := indicator + " " + task.label
	if task.total > 0 {
		line += " " + progressBar(task.current, task.total, m.unicode)
	}
	return line
}

func progressBar(current, total int64, unicode bool) string {
	const width = 10
	current = max(0, min(current, total))
	filled := 0
	for filled < width && int64(filled+1)*total <= current*int64(width) {
		filled++
	}
	full, empty := "#", "-"
	if unicode {
		full, empty = "━", "─"
	}
	percentage := current * 100 / total
	return fmt.Sprintf("[%s%s] %d%%", strings.Repeat(full, filled), strings.Repeat(empty, width-filled), percentage)
}
