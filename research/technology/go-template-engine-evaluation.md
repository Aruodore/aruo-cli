# Go Template Engine Evaluation

Status: completed  
Reviewed: 2026-08-04

## Question

Which rendering and filesystem primitives give Aruo deterministic repository generation without creating an unsafe general-purpose execution environment?

## Findings

The standard library provides the required core. `text/template` supplies substitution and conditions, and its `missingkey=error` option prevents silent missing map values. It permits custom functions but warns that templates executing against unsafe objects are not a security sandbox. Aruo therefore provides only plain typed metadata, recursively validated JSON-like variables, and a minimal pure function set.

`embed.FS` implements `io/fs.FS`, is immutable and safe for concurrent reads. The same `fs.FS` contract is implemented by test filesystems and directory-backed sources. `io/fs` also standardizes unrooted slash paths across operating systems. Official documentation warns that `os.DirFS` and `fs.Sub` do not prevent symlink escape; containment is an artifact-ingestion responsibility.

Sprig and gomplate demonstrate the convenience and risk of broad function/data-source catalogs. Their surface includes capabilities Aruo intentionally excludes from rendering. Jinja-, Django-, Liquid-, Handlebars-, and high-throughput engines offer alternate syntax or request-rendering speed, but add dependencies and compatibility obligations without improving repository planning.

Afero remains useful for applications requiring a mutable virtual filesystem. Aruo's renderer needs read-only sources and returns a plan instead of mutating a destination, so standard interfaces are smaller and sufficient.

## Recommendation

Use `text/template`, `embed`, and `io/fs`. Keep the renderer pure and internal. Do not use Sprig, gomplate, Afero, or an alternate template language in the core. Treat plugin templates as versioned data bundles and never accept plugin-defined in-process template functions.

## Primary sources

- [Go `text/template` documentation](https://pkg.go.dev/text/template)
- [Go `embed` documentation](https://pkg.go.dev/embed)
- [Go `io/fs` documentation](https://pkg.go.dev/io/fs)
- [Go `os.DirFS` and `Root.FS` documentation](https://pkg.go.dev/os)
- [Sprig repository](https://github.com/Masterminds/sprig)
- [gomplate repository](https://github.com/hairyhenderson/gomplate)
- [Pongo2 repository](https://github.com/flosch/pongo2)
- [Liquid repository](https://github.com/osteele/liquid)
- [Afero repository](https://github.com/spf13/afero)
