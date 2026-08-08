package templateengine_test

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"reflect"
	"strings"
	"sync"
	"testing"
	"testing/fstest"

	"github.com/aruodore/aruo-cli/internal/templateengine"
)

func TestRenderSubstitutionConditionRawModeAndOrder(t *testing.T) {
	t.Parallel()

	source := fstest.MapFS{
		"readme.tmpl": {Data: []byte("# {{ upper .Project.Name }}\n{{ if .Variables.Public }}public{{ else }}private{{ end }}\n")},
		"raw.bin":     {Data: []byte{0x00, '{', '{', 0xff}},
	}
	engine := newEngine(t, source, templateengine.Options{})
	blueprint := templateengine.Blueprint{
		ID:       "test/go",
		Language: "go",
		Files: []templateengine.FileSpec{
			{Source: "readme.tmpl", Destination: "docs/{{ kebab .Project.Name }}.md", Mode: 0o600, Template: true},
			{Source: "raw.bin", Destination: "assets/raw.bin"},
		},
	}

	plan, err := engine.Render(context.Background(), blueprint, validData())
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if plan.BlueprintID != blueprint.ID {
		t.Errorf("BlueprintID = %q, want %q", plan.BlueprintID, blueprint.ID)
	}
	if got := paths(plan); !reflect.DeepEqual(got, []string{"assets/raw.bin", "docs/acme-tools.md"}) {
		t.Errorf("paths = %v", got)
	}
	if !bytes.Equal(plan.Files[0].Content, source["raw.bin"].Data) {
		t.Errorf("raw content = %v", plan.Files[0].Content)
	}
	if plan.Files[0].Mode != 0o644 {
		t.Errorf("raw mode = %v, want 0644", plan.Files[0].Mode)
	}
	if got := string(plan.Files[1].Content); got != "# ACME TOOLS\npublic\n" {
		t.Errorf("rendered content = %q", got)
	}
	if plan.Files[1].Mode != 0o600 {
		t.Errorf("rendered mode = %v, want 0600", plan.Files[1].Mode)
	}
}

func TestRenderIsDeterministic(t *testing.T) {
	t.Parallel()

	engine := newEngine(t, fstest.MapFS{
		"a": {Data: []byte("{{ .Project.Name }}")},
		"b": {Data: []byte("{{ .Project.Module }}")},
	}, templateengine.Options{})
	blueprint := templateengine.Blueprint{ID: "test/go", Language: "go", Files: []templateengine.FileSpec{
		{Source: "b", Destination: "b", Template: true},
		{Source: "a", Destination: "a", Template: true},
	}}

	first, err := engine.Render(context.Background(), blueprint, validData())
	if err != nil {
		t.Fatal(err)
	}
	second, err := engine.Render(context.Background(), blueprint, validData())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Errorf("plans differ:\n%#v\n%#v", first, second)
	}
}

func TestRenderFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		blueprint templateengine.Blueprint
		data      templateengine.Data
		wantStage templateengine.Stage
		want      string
	}{
		{
			name:      "missing variable",
			blueprint: oneFile("value", "out", true),
			data:      validData(),
			wantStage: templateengine.StageExecute,
			want:      "Missing",
		},
		{
			name:      "template syntax",
			blueprint: oneFile("syntax", "out", true),
			data:      validData(),
			wantStage: templateengine.StageParse,
			want:      "unexpected EOF",
		},
		{
			name: "duplicate destination",
			blueprint: templateengine.Blueprint{ID: "test/go", Language: "go", Files: []templateengine.FileSpec{
				{Source: "plain", Destination: "same"},
				{Source: "other", Destination: "same"},
			}},
			data:      validData(),
			wantStage: templateengine.StageValidate,
			want:      "collides",
		},
		{
			name:      "destination traversal",
			blueprint: oneFile("plain", "../escape", false),
			data:      validData(),
			wantStage: templateengine.StageValidate,
			want:      "invalid destination",
		},
		{
			name:      "language mismatch",
			blueprint: oneFile("plain", "out", false),
			data: templateengine.Data{Project: templateengine.Project{
				Name: "Acme", Language: "python",
			}},
			wantStage: templateengine.StageValidate,
			want:      "does not match",
		},
		{
			name:      "unsafe callable variable",
			blueprint: oneFile("plain", "out", false),
			data: templateengine.Data{
				Project:   templateengine.Project{Name: "Acme", Language: "go"},
				Variables: map[string]any{"unsafe": func() string { return "no" }},
			},
			wantStage: templateengine.StageValidate,
			want:      "unsupported value type",
		},
	}

	source := fstest.MapFS{
		"value":  {Data: []byte("{{ .Variables.Missing }}")},
		"syntax": {Data: []byte("{{ if true }}")},
		"plain":  {Data: []byte("plain")},
		"other":  {Data: []byte("other")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			engine := newEngine(t, source, templateengine.Options{})
			_, err := engine.Render(context.Background(), test.blueprint, test.data)
			if err == nil {
				t.Fatal("Render() error = nil")
			}
			var renderErr *templateengine.Error
			if !errors.As(err, &renderErr) {
				t.Fatalf("error type = %T, want *templateengine.Error", err)
			}
			if renderErr.Stage != test.wantStage {
				t.Errorf("stage = %q, want %q", renderErr.Stage, test.wantStage)
			}
			if !strings.Contains(err.Error(), test.want) {
				t.Errorf("error = %q, want substring %q", err, test.want)
			}
		})
	}
}

func TestLimits(t *testing.T) {
	t.Parallel()

	t.Run("source", func(t *testing.T) {
		t.Parallel()
		engine := newEngine(t, fstest.MapFS{"large": {Data: []byte("12345")}}, templateengine.Options{
			MaxTemplateBytes: 4,
			MaxOutputBytes:   10,
		})
		_, err := engine.Render(context.Background(), oneFile("large", "out", false), validData())
		if !strings.Contains(errorString(err), "size limit exceeded") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("rendered output", func(t *testing.T) {
		t.Parallel()
		engine := newEngine(t, fstest.MapFS{"expand": {Data: []byte("{{ .Project.Name }}")}}, templateengine.Options{
			MaxTemplateBytes: 100,
			MaxOutputBytes:   4,
		})
		_, err := engine.Render(context.Background(), oneFile("expand", "out", true), validData())
		if !strings.Contains(errorString(err), "size limit exceeded") {
			t.Fatalf("error = %v", err)
		}
	})
}

func TestCancellationReturnsNoPlan(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	engine := newEngine(t, fstest.MapFS{"plain": {Data: []byte("plain")}}, templateengine.Options{})
	plan, err := engine.Render(ctx, oneFile("plain", "out", false), validData())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
	if len(plan.Files) != 0 {
		t.Errorf("partial files = %d, want 0", len(plan.Files))
	}
}

func TestEngineSupportsConcurrentRenders(t *testing.T) {
	t.Parallel()

	engine := newEngine(t, fstest.MapFS{"value": {Data: []byte("{{ .Project.Name }}")}}, templateengine.Options{})
	blueprint := oneFile("value", "out", true)
	var group sync.WaitGroup
	for range 16 {
		group.Add(1)
		go func() {
			defer group.Done()
			plan, err := engine.Render(context.Background(), blueprint, validData())
			if err != nil {
				t.Errorf("Render() error = %v", err)
				return
			}
			if got := string(plan.Files[0].Content); got != "Acme Tools" {
				t.Errorf("content = %q", got)
			}
		}()
	}
	group.Wait()
}

func TestNewValidation(t *testing.T) {
	t.Parallel()

	if _, err := templateengine.New(nil, templateengine.Options{}); err == nil {
		t.Error("New(nil) error = nil")
	}
	if _, err := templateengine.New(fstest.MapFS{}, templateengine.Options{MaxOutputBytes: -1}); err == nil {
		t.Error("New(negative limit) error = nil")
	}
}

func BenchmarkRenderRepository(b *testing.B) {
	source := make(fstest.MapFS, 40)
	files := make([]templateengine.FileSpec, 0, 40)
	for i := range 40 {
		name := "source/" + string(rune('a'+i/10)) + strings.Repeat("x", i%10) + ".tmpl"
		destination := "generated/" + string(rune('a'+i/10)) + strings.Repeat("x", i%10) + ".txt"
		source[name] = &fstest.MapFile{Data: []byte("{{ .Project.Name }} {{ .Project.Module }} {{ if .Variables.Public }}public{{ end }}\n")}
		files = append(files, templateengine.FileSpec{Source: name, Destination: destination, Template: true})
	}
	engine, err := templateengine.New(source, templateengine.Options{})
	if err != nil {
		b.Fatal(err)
	}
	blueprint := templateengine.Blueprint{ID: "benchmark/go", Language: "go", Files: files}
	data := validData()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		plan, renderErr := engine.Render(context.Background(), blueprint, data)
		if renderErr != nil {
			b.Fatal(renderErr)
		}
		if len(plan.Files) != len(files) {
			b.Fatalf("files = %d", len(plan.Files))
		}
	}
}

func newEngine(t testing.TB, source fs.FS, options templateengine.Options) *templateengine.Engine {
	t.Helper()
	engine, err := templateengine.New(source, options)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return engine
}

func validData() templateengine.Data {
	return templateengine.Data{
		Project: templateengine.Project{
			Name:        "Acme Tools",
			Module:      "example.com/acme/tools",
			Description: "Reliable tools.",
			Author:      "Acme",
			License:     "Apache-2.0",
			Language:    "go",
		},
		Variables: map[string]any{"Public": true, "IncludeInstall": true, "GoVersion": "1.26.0"},
	}
}

func oneFile(source, destination string, render bool) templateengine.Blueprint {
	return templateengine.Blueprint{
		ID:       "test/go",
		Language: "go",
		Files: []templateengine.FileSpec{{
			Source: source, Destination: destination, Template: render,
		}},
	}
}

func paths(plan templateengine.Plan) []string {
	result := make([]string, len(plan.Files))
	for index, file := range plan.Files {
		result[index] = file.Path
	}
	return result
}

func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
