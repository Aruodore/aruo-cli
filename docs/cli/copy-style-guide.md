# CLI copy style guide

This is the short reference for the actual words Aruo shows a user: prompt
labels, descriptions, examples, confirmations, errors, and status text. The
[Terminal UX specification](terminal-ux.md) is the normative contract for
*mechanism and behavior* (which widget to use, how cancellation propagates,
what a screen reader needs); this page is for *wording*: the sentences that
go inside that mechanism. If a rule needs a `MUST`/`SHOULD` and a paragraph
of justification, it belongs in the spec, not here. This page works by
example: find the closest pattern below and match it.

## Prompt labels and descriptions

- Use the words of the target ecosystem, not Aruo's internal model names.
  `create`'s `--module` flag means three different real things depending on
  the template, and the prompt says so explicitly instead of using one
  generic label:

  | Template | Label | Example shown |
  | --- | --- | --- |
  | `go-library` | `Go module path` | `github.com/your-name/my-library` |
  | `js-library`, `ts-library`, `react`, ... | `npm package name` | `my-library` |
  | `python-library` | `PyPI package name` | `my-library` |

  The flag's own `--help` text stays ecosystem-neutral since it's shared
  across all of them (`package or module identifier for the template's
  ecosystem`); the prompt shown at the point of use is where the specific
  word choice belongs.

- Every prompt shows a realistic `Example`, not an abstract placeholder like
  `<value>` or `foo`. Use a name a real project would have.

- If one option among several has a clear recommendation, put it first and
  say `Recommended`. Don't do this for unrelated choices. Aruo's 9 templates
  aren't "recommended vs. not," they're different ecosystems; this rule
  applies when there's a real default a user would usually want.

- Before a flat list of every template gets long enough to scroll, split it:
  `create`'s first interactive screen asks `What are you building?` with two
  options, `Application` and `Library`, each carrying a one-line description
  built from the real template names in that kind (`Next.js application,
  Nuxt application, React application`) rather than a generic blurb. The
  answer filters the second screen down to that kind's 3-5 templates. This
  step is skipped whenever `--kind`, `--template`, or a non-interactive
  session already answers the question, and disappears entirely if the
  catalog ever shrinks to one kind.

## Optional fields

Say `Optional` on the label itself, not buried mid-sentence in the
description. A user scanning prompts reads the label first.

```
Bad:   Label: "Short description"
       Description: "One sentence explaining what the project does.
                      Optional; leave blank to fill in later."

Good:  Label: "Short description (Optional)"
       Description: "One sentence explaining what the project does."
```

An early version of this defaulted to visible `TODO: ...` placeholder text
written straight into generated files (`Copyright (c) TODO: set an
author`). Don't do that: it reads as broken, not helpful, and it's not
automatic; it just moves the manual step into the generated output instead
of removing it. Prefer a real value derived without asking:

```
--description left blank  ->  the catalog entry's own Description
                               ("A production-ready Go library with...")
--author left blank       ->  `git config --get user.name`, best-effort,
                               falling back to an empty string
                               (not a placeholder) if git or the config
                               entry isn't there
--module left blank       ->  the project's own name (the go.mod module
                               path / package.json "name" / pyproject.toml
                               "name" all start out equal to it)
```

An empty string in the output (`Copyright (c) ` with nothing after) is a
better failure mode than fabricated text: it's incomplete rather than
looking like plausible but wrong content.

`--module` isn't just defaulted, it isn't asked about at all: there's no
module/package-name screen in the guided flow anymore, interactive or not.
It silently becomes the project's own name unless `--module` is passed
explicitly. This is a deliberate exception to "only default a field when a
wrong default is merely incomplete, not broken": a bare project name is a
valid `go.mod` module path, `name` in `package.json`, and `name` in
`pyproject.toml` for every template today, so it's not *broken*, just not
the ideal value for a Go module meant to be published (`go get`-able
projects want a domain-qualified path like `github.com/you/name`, which
this can't guess). Asking the same question twice, once for the project
name and once for the package name, was judged worse than that gap. Pass
`--module` explicitly for a real import path.

A default never occupies the editable input as pre-filled text the user
has to delete to leave the field blank. Show it as a hint instead, and
apply it after the fact only if the user submits nothing:

```
Bad (rich prompter): input box already contains "Aruodore",
                      cursor at the end; the user backspaces
                      four characters to leave it blank.

Good: Author or organization (Optional) [Aruodore]:
      (empty input box; pressing Enter without typing
       submits Aruodore, but the box was never occupied)
```

The plain adapter already worked this way (`Label [Default]` is
informational text in the prompt line, not something sitting in a buffer);
match its contract when writing another prompter rather than reinventing
one where a default subtly requires deleting.

## Back-navigation

`create`'s guided flow (`Prompter.Guide`) lets a user return to an earlier
screen instead of only being able to move forward. The rich prompter gets
this for free from Huh's own `shift+tab`, with no copy of its own to write.
Its default help bar already advertises `shift+tab back`. The accessible
adapter has no such built-in, so it prints one line before the first
question, and only when there's more than one question to navigate between:

```
Type back at any prompt to return to the previous question.
```

Typing `back` (case-insensitive, trimmed) at any prompt in that flow is a
reserved navigation command, not a literal value: it never becomes a
project name, module path, description, or author. That's a deliberate
trade-off: if a user wants the literal text "back" as a value,
the corresponding flag (`--name`, `--module`, `--description`, `--author`)
bypasses the prompt entirely and takes the value as given. Revisiting an
already-answered screen shows that answer as its default (`[Aruodore]`,
`(default)` on a list), so pressing Enter alone keeps it rather than
forcing a retype.

## Confirmations

`create`'s confirmation today is generic: `Create this project?
[y/N]`, preceded by a summary block naming the destination, template, and
license. It isn't target/count-specific, since there's exactly one
confirmation in the product right now. Once a command confirms an
*effect* rather than a one-time creation, name the exact target and count:

```
Apply these 4 changes to github.com/acme/api? [y/N]
```

Uppercase marks the default; Enter accepts it. Use `y`/`yes`/`n`/`no`: they
stay stable across locales even before Aruo supports translated prompts.

## Errors

Today's renderer (`internal/cli/run.go`) is one line,
`fmt.Fprintf(ErrOut, "Error: %v\n", err)`, so the entire error *is* the
wrapped Go error string. That constrains the writing more, not less: it has
to say everything in one plain-language sentence, since there's no second
line to fall back on.

```
Error: destination already exists: /work/my-library
Error: Go module path is required; provide --module
Error: unknown flag: --bogus-flag
Error: end of input
```

That last one is what a closed/EOF stdin produces today
(`tux.ErrEndOfInput`, unwrapped by the same one-line renderer), not a
friendlier "Input ended; project creation was cancelled" message.
`terminal-ux.md` names that friendlier wording as the target; nothing
renders it yet, so don't assume it exists when writing code that expects
a specific EOF message.

- Say what failed in plain language, not the internal error type.
- Name the exact path/value involved.
- Lowercase, no trailing punctuation, matches `STYLE_GUIDE.md`'s Go error
  string convention, since these strings usually *are* Go error strings.
- Give a next action where the sentence has room for one (`provide
  --module`), rather than a bare fact with no path forward.
- An expected cancellation prints `Cancelled.`, never a red `Error:` line.

**Not implemented yet**: planned, multi-line anatomy (separate cause,
consequence, and next-action lines, a debug reference) from
`terminal-ux.md`'s "User-facing error anatomy". Don't write code today as
though this format exists. When it lands, this section should show the
real multi-line output, not the aspirational one.

## Status and progress text

No command calls `ProgressSink`/`Session.Progress()` yet: `create`'s write
and `doctor`'s audit both finish under the "<200ms: no progress" threshold
`terminal-ux.md` itself sets, so there's no live progress text to check
against reality today. When a command's work is slow enough to need it,
name the actual work, not a generic wait state (`Downloading blueprint`,
not `Please wait...`). Until then this is a rule for the future, not a
description of current output.

The one status line that *is* real is the final success message, shown
below:

```
Created go-library with 24 files at ./my-library
```

A durable summary a user could paste into a chat, not just a checkmark:
`presentCreated` in `internal/cli/command/create.go`.

## Cancellation messages

```
Cancelling... Press Ctrl+C again to exit immediately.
Cancelled.
```

Both verbatim (`internal/tux/lifecycle/manager.go`, `internal/cli/run.go`):
three ASCII periods, not a Unicode ellipsis (`…`), matching what's printed
rather than what reads better in prose. `terminal-ux.md`'s per-outcome
variants (`Cancelled. No repository changes were applied.` /
`Cancelled after 3 of 5 operations...`) are the target once a command has
partial-completion states to distinguish; today every cancellation prints
the same plain `Cancelled.`

## Keeping this page honest

When you add a template, a prompt, or a command, and you catch yourself
writing wording that doesn't match a pattern above: fix the pattern here
too, in the same change. This page is only useful if it reflects what the
CLI says. See the [terminal UX specification's implementation
status](terminal-ux.md#implementation-status-2026-08-06) for the equivalent
discipline applied to behavior.
