package templateengine

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"path"
	"reflect"
	"sort"
	"strings"
	"text/template"
	"unicode"
)

const (
	defaultMaxTemplateBytes = 1 << 20
	defaultMaxOutputBytes   = 4 << 20
	maxDestinationBytes     = 4096
	maxValueDepth           = 32
)

var errSizeLimit = errors.New("size limit exceeded")

// Options configures resource limits. Zero values select conservative defaults.
type Options struct {
	MaxTemplateBytes int64
	MaxOutputBytes   int64
}

// Engine renders templates from one read-only bundle filesystem.
type Engine struct {
	source           fs.FS
	maxTemplateBytes int64
	maxOutputBytes   int64
}

// New creates an engine over source.
func New(source fs.FS, options Options) (*Engine, error) {
	if source == nil {
		return nil, errors.New("template source filesystem is required")
	}
	if options.MaxTemplateBytes < 0 || options.MaxOutputBytes < 0 {
		return nil, errors.New("template limits cannot be negative")
	}
	if options.MaxTemplateBytes == 0 {
		options.MaxTemplateBytes = defaultMaxTemplateBytes
	}
	if options.MaxOutputBytes == 0 {
		options.MaxOutputBytes = defaultMaxOutputBytes
	}

	return &Engine{
		source:           source,
		maxTemplateBytes: options.MaxTemplateBytes,
		maxOutputBytes:   options.MaxOutputBytes,
	}, nil
}

// Render validates and renders blueprint into a destination-sorted file plan.
func (e *Engine) Render(ctx context.Context, blueprint Blueprint, data Data) (Plan, error) {
	if err := validateInputs(blueprint, data); err != nil {
		return Plan{}, renderError(StageValidate, blueprint.ID, "", "", err)
	}

	files := make([]File, 0, len(blueprint.Files))
	destinations := make(map[string]string, len(blueprint.Files))
	for _, spec := range blueprint.Files {
		if err := ctx.Err(); err != nil {
			return Plan{}, renderError(StageExecute, blueprint.ID, spec.Source, "", err)
		}

		destination, err := e.renderDestination(spec.Destination, data)
		if err != nil {
			return Plan{}, renderError(StageExecute, blueprint.ID, spec.Source, spec.Destination, err)
		}
		if validationErr := validateDestination(destination); validationErr != nil {
			return Plan{}, renderError(StageValidate, blueprint.ID, spec.Source, destination, validationErr)
		}
		if previous, exists := destinations[destination]; exists {
			return Plan{}, renderError(
				StageValidate,
				blueprint.ID,
				spec.Source,
				destination,
				fmt.Errorf("destination collides with source %q", previous),
			)
		}
		destinations[destination] = spec.Source

		source, err := e.readSource(spec.Source)
		if err != nil {
			return Plan{}, renderError(StageRead, blueprint.ID, spec.Source, destination, err)
		}

		content := source
		if spec.Template {
			content, err = e.renderContent(spec.Source, source, data)
			if err != nil {
				return Plan{}, renderError(stageForTemplateError(err), blueprint.ID, spec.Source, destination, err)
			}
		} else if int64(len(content)) > e.maxOutputBytes {
			return Plan{}, renderError(StageExecute, blueprint.ID, spec.Source, destination, errSizeLimit)
		}

		mode := spec.Mode
		if mode == 0 {
			mode = 0o644
		}
		files = append(files, File{Path: destination, Content: content, Mode: mode, Source: spec.Source})
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return Plan{BlueprintID: blueprint.ID, Files: files}, nil
}

func (e *Engine) renderDestination(source string, data Data) (string, error) {
	if len(source) == 0 || len(source) > maxDestinationBytes {
		return "", errors.New("destination template must be between 1 and 4096 bytes")
	}
	tmpl, err := newTemplate("destination").Parse(source)
	if err != nil {
		return "", err
	}
	writer := &limitedBuffer{remaining: maxDestinationBytes}
	if err := tmpl.Execute(writer, data); err != nil {
		return "", err
	}
	return writer.String(), nil
}

func (e *Engine) renderContent(name string, source []byte, data Data) ([]byte, error) {
	tmpl, err := newTemplate(name).Parse(string(source))
	if err != nil {
		return nil, &templateStageError{stage: StageParse, err: err}
	}
	writer := &limitedBuffer{remaining: e.maxOutputBytes}
	if err := tmpl.Execute(writer, data); err != nil {
		return nil, &templateStageError{stage: StageExecute, err: err}
	}
	return bytes.Clone(writer.Bytes()), nil
}

func (e *Engine) readSource(name string) ([]byte, error) {
	file, err := e.source.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(io.LimitReader(file, e.maxTemplateBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(content)) > e.maxTemplateBytes {
		return nil, errSizeLimit
	}
	return content, nil
}

func newTemplate(name string) *template.Template {
	return template.New(name).Option("missingkey=error").Funcs(template.FuncMap{
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"kebab": kebab,
	})
}

func kebab(value string) string {
	var result strings.Builder
	separator := false
	for _, char := range strings.TrimSpace(value) {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			if separator && result.Len() > 0 {
				result.WriteByte('-')
			}
			result.WriteRune(unicode.ToLower(char))
			separator = false
			continue
		}
		separator = true
	}
	return result.String()
}

func validateInputs(blueprint Blueprint, data Data) error {
	if blueprint.ID == "" {
		return errors.New("blueprint ID is required")
	}
	if blueprint.Language == "" {
		return errors.New("blueprint language is required")
	}
	if data.Project.Name == "" {
		return errors.New("project name is required")
	}
	if data.Project.Language != blueprint.Language {
		return fmt.Errorf("project language %q does not match blueprint language %q", data.Project.Language, blueprint.Language)
	}
	if err := validateValue(reflect.ValueOf(data.Variables), 0); err != nil {
		return fmt.Errorf("variables: %w", err)
	}
	for _, file := range blueprint.Files {
		if file.Source == "." || !fs.ValidPath(file.Source) || strings.Contains(file.Source, `\`) {
			return fmt.Errorf("invalid source path %q", file.Source)
		}
		if file.Mode&^fs.ModePerm != 0 {
			return fmt.Errorf("invalid portable mode %v for source %q", file.Mode, file.Source)
		}
	}
	return nil
}

func validateValue(value reflect.Value, depth int) error {
	if depth > maxValueDepth {
		return errors.New("maximum nesting depth exceeded")
	}
	if !value.IsValid() {
		return nil
	}
	if value.Kind() == reflect.Interface {
		if value.IsNil() {
			return nil
		}
		return validateValue(value.Elem(), depth+1)
	}
	switch value.Kind() {
	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return nil
	case reflect.Float32, reflect.Float64:
		if math.IsNaN(value.Float()) || math.IsInf(value.Float(), 0) {
			return errors.New("non-finite numbers are not supported")
		}
		return nil
	case reflect.Slice, reflect.Array:
		for i := range value.Len() {
			if err := validateValue(value.Index(i), depth+1); err != nil {
				return err
			}
		}
		return nil
	case reflect.Map:
		if value.Type().Key().Kind() != reflect.String {
			return errors.New("map keys must be strings")
		}
		iterator := value.MapRange()
		for iterator.Next() {
			if err := validateValue(iterator.Value(), depth+1); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unsupported value type %s", value.Type())
	}
}

func validateDestination(destination string) error {
	if destination == "." || !fs.ValidPath(destination) {
		return fmt.Errorf("invalid destination path %q", destination)
	}
	if strings.Contains(destination, `\`) || path.Clean(destination) != destination {
		return fmt.Errorf("destination must be a clean slash-separated relative path: %q", destination)
	}
	return nil
}

func renderError(stage Stage, blueprintID, source, destination string, err error) error {
	return &Error{Stage: stage, BlueprintID: blueprintID, Source: source, Destination: destination, Err: err}
}

type templateStageError struct {
	stage Stage
	err   error
}

func (e *templateStageError) Error() string { return e.err.Error() }
func (e *templateStageError) Unwrap() error { return e.err }

func stageForTemplateError(err error) Stage {
	var staged *templateStageError
	if errors.As(err, &staged) {
		return staged.stage
	}
	return StageExecute
}

type limitedBuffer struct {
	bytes.Buffer
	remaining int64
}

func (b *limitedBuffer) Write(content []byte) (int, error) {
	if int64(len(content)) > b.remaining {
		return 0, errSizeLimit
	}
	written, err := b.Buffer.Write(content)
	b.remaining -= int64(written)
	return written, err
}
