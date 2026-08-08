package charm

import (
	"reflect"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/aruodore/aruo-cli/internal/tux"
)

var (
	keyDown     = tea.Key{Code: tea.KeyDown}
	keyUp       = tea.Key{Code: tea.KeyUp}
	keyEnter    = tea.Key{Code: tea.KeyEnter}
	keyShiftTab = tea.Key{Code: tea.KeyTab, Mod: tea.ModShift}
)

// drainGuideForm feeds cmd (and any further tea.Cmd it triggers,
// recursively) into the form until nothing is pending, simulating enough
// of Bubble Tea's real event loop -- including huh's async
// TitleFunc/OptionsFunc recompute cycle and Form.Init's tea.Sequence of
// per-group Init cmds -- to observe the result of one simulated keystroke
// without a real terminal program or PTY. tea.BatchMsg and the internal,
// unexported sequenceMsg are both just named []tea.Cmd under the hood, so
// they're unwrapped generically via reflection rather than needing the
// unexported type by name.
//
// Some cmds (cursor blink, spinner ticks) block on a real timer/channel by
// design and are irrelevant to the state this test cares about; each cmd
// gets a short window to resolve, and is discarded rather than awaited
// forever if it doesn't.
func drainGuideForm(t *testing.T, model huh.Model, cmd tea.Cmd) huh.Model {
	t.Helper()
	pending := []tea.Cmd{cmd}
	for iterations := 0; len(pending) > 0; iterations++ {
		if iterations > 200 {
			t.Fatal("drainGuideForm: too many pending commands, possible infinite loop")
		}
		next := pending[0]
		pending = pending[1:]
		if next == nil {
			continue
		}
		msg, ok := callWithTimeout(next)
		if !ok || msg == nil {
			continue
		}
		if more, ok := asCmdSlice(msg); ok {
			pending = append(pending, more...)
			continue
		}
		var newCmd tea.Cmd
		model, newCmd = model.Update(msg)
		if newCmd != nil {
			pending = append(pending, newCmd)
		}
	}
	return model
}

// callWithTimeout runs cmd and reports whether it resolved within a short
// window. A cmd that blocks longer than that is assumed to be a real-time
// animation command (cursor blink, spinner tick) rather than one that
// carries state this test needs, and is abandoned -- its goroutine may
// leak past the deadline, which is acceptable for a short-lived test.
func callWithTimeout(cmd tea.Cmd) (tea.Msg, bool) {
	result := make(chan tea.Msg, 1)
	go func() { result <- cmd() }()
	select {
	case msg := <-result:
		return msg, true
	case <-time.After(20 * time.Millisecond):
		return nil, false
	}
}

// asCmdSlice reports whether msg is a slice of tea.Cmd (tea.BatchMsg,
// bubbletea's unexported sequenceMsg, or anything shaped like them).
func asCmdSlice(msg tea.Msg) ([]tea.Cmd, bool) {
	value := reflect.ValueOf(msg)
	if value.Kind() != reflect.Slice {
		return nil, false
	}
	cmds := make([]tea.Cmd, 0, value.Len())
	for i := 0; i < value.Len(); i++ {
		cmd, ok := value.Index(i).Interface().(tea.Cmd)
		if !ok {
			return nil, false
		}
		cmds = append(cmds, cmd)
	}
	return cmds, true
}

func sendKey(t *testing.T, model huh.Model, key tea.Key) huh.Model {
	t.Helper()
	model, cmd := model.Update(tea.KeyPressMsg(key))
	return drainGuideForm(t, model, cmd)
}

func focusedKey(t *testing.T, model huh.Model) string {
	t.Helper()
	form, ok := model.(*huh.Form)
	if !ok {
		t.Fatalf("model = %T, want *huh.Form", model)
	}
	focused := form.GetFocusedField()
	if focused == nil {
		t.Fatal("GetFocusedField() = nil")
	}
	return focused.GetKey()
}

func containsOptionID(options []tux.Option, value any) bool {
	id, ok := value.(tux.OptionID)
	if !ok {
		return false
	}
	for _, option := range options {
		if option.ID == id {
			return true
		}
	}
	return false
}

// TestGuideFormOptionsFuncNarrowsToKind exercises the exact reactive shape
// create's flow depends on: a later Select step's option list narrows to
// an earlier Select step's live answer, and re-narrows correctly after the
// user goes back and changes that earlier answer.
func TestGuideFormOptionsFuncNarrowsToKind(t *testing.T) {
	t.Parallel()

	appOptions := []tux.Option{{ID: "next-app"}, {ID: "nuxt-app"}, {ID: "react-app"}}
	libraryOptions := []tux.Option{{ID: "go-library"}, {ID: "js-library"}, {ID: "python-library"}, {ID: "ts-library"}, {ID: "vue-library"}}

	steps := []tux.Step{
		{
			ID:   "kind",
			Kind: tux.StepSelect,
			Select: func(tux.Answers) tux.SelectRequest {
				return tux.SelectRequest{
					Label: "What are you building?",
					Options: []tux.Option{
						{ID: "app", Label: "Application"},
						{ID: "library", Label: "Library"},
					},
				}
			},
		},
		{
			ID:   "template",
			Kind: tux.StepSelect,
			Select: func(answers tux.Answers) tux.SelectRequest {
				kind, _ := answers["kind"].(tux.OptionID)
				options := libraryOptions
				if kind == "app" {
					options = appOptions
				}
				return tux.SelectRequest{Label: "Template", Options: options}
			},
		},
	}

	prompter := NewPrompter(nil, nil, tux.Capabilities{Color: tux.ColorTrueColor}, tux.Policy{})
	form, fields := prompter.buildGuideForm(steps)

	var model huh.Model = form
	model = drainGuideForm(t, model, form.Init())

	// kind defaults to "app" (first option); move down to "library" and
	// advance into the template step.
	model = sendKey(t, model, keyDown)
	model = sendKey(t, model, keyEnter)

	if !containsOptionID(libraryOptions, fields[1].value()) {
		t.Fatalf("template value = %v, want a library-kind option after choosing kind=library", fields[1].value())
	}

	// Back to kind, switch to "app", forward again -- proving OptionsFunc
	// recomputes on re-entry rather than staying stuck on the old list.
	model = sendKey(t, model, keyShiftTab)
	if got := focusedKey(t, model); got != "kind" {
		t.Fatalf("focused field = %q, want kind", got)
	}
	model = sendKey(t, model, keyUp)
	sendKey(t, model, keyEnter)

	if !containsOptionID(appOptions, fields[1].value()) {
		t.Fatalf("template value = %v, want an app-kind option after switching kind back to app", fields[1].value())
	}
}

// TestGuideFormSkipsHiddenStepBothDirections proves a Step whose Skip is
// true never becomes reachable, in either navigation direction.
func TestGuideFormSkipsHiddenStepBothDirections(t *testing.T) {
	t.Parallel()

	makeStep := func(id string, skip func() bool) tux.Step {
		return tux.Step{
			ID:   id,
			Kind: tux.StepSelect,
			Select: func(tux.Answers) tux.SelectRequest {
				return tux.SelectRequest{Label: id, Options: []tux.Option{{ID: tux.OptionID(id + "-only")}}}
			},
			Skip: skip,
		}
	}
	steps := []tux.Step{
		makeStep("a", nil),
		makeStep("b", func() bool { return true }),
		makeStep("c", nil),
	}

	prompter := NewPrompter(nil, nil, tux.Capabilities{}, tux.Policy{})
	form, _ := prompter.buildGuideForm(steps)
	var model huh.Model = form
	model = drainGuideForm(t, model, form.Init())

	model = sendKey(t, model, keyEnter)
	if got := focusedKey(t, model); got != "c" {
		t.Fatalf("focused field = %q, want c (b is hidden and must be skipped forward)", got)
	}

	model = sendKey(t, model, keyShiftTab)
	if got := focusedKey(t, model); got != "a" {
		t.Fatalf("focused field = %q, want a (b is hidden and must be skipped backward)", got)
	}
}

// TestGuideFormCollectAnswersAppliesInputDefaultOnEmptySubmission proves
// the Guide path preserves the earlier session's fix: a default never
// occupies the editable buffer, and is only substituted after a genuinely
// empty submission, once the form completes.
func TestGuideFormCollectAnswersAppliesInputDefaultOnEmptySubmission(t *testing.T) {
	t.Parallel()

	defaultAuthor := "Aruodore"
	steps := []tux.Step{
		{
			ID:   "author",
			Kind: tux.StepInput,
			Input: func(tux.Answers) tux.InputRequest {
				return tux.InputRequest{Label: "Author", Optional: true, Default: &defaultAuthor}
			},
		},
	}
	prompter := NewPrompter(nil, nil, tux.Capabilities{}, tux.Policy{})
	form, fields := prompter.buildGuideForm(steps)
	var model huh.Model = form
	model = drainGuideForm(t, model, form.Init())
	sendKey(t, model, keyEnter)

	answers := collectGuideAnswers(steps, fields)
	if answers["author"] != defaultAuthor {
		t.Fatalf("answers[author] = %v, want the default %q substituted after an empty submission", answers["author"], defaultAuthor)
	}
}
