# CLI Design Direction

## Intent

Aruo is a focused engineering instrument. Its interface should make consequential actions easy to understand and repository health easy to scan. It should feel calm, precise, and trustworthy rather than theatrical or "hacker" themed.

## Experience

- Register: developer tool.
- Density: compact, with blank space reserved for separating stages rather than decorating them.
- Visual character: terminal-native text with restrained semantic outcome colors.
- Typography: the terminal's own typeface; hierarchy comes from weight, case, spacing, and alignment.
- Motion: only live work receives a spinner. Completed, failed, cancelled, and redirected output is durable and static.
- Responsive behavior: lower-priority table columns disappear before content becomes unreadable.

## Interface contract

- Commands follow title, context, detail, next-action hierarchy.
- Interactive prompts use terminal-default foreground and background colors. Short categorical choices use a filled active state. Framework and library lists retain their established ecosystem colors and underline the complete active option.
- Green, amber, and red are reserved for success, warning, and failure.
- Markers and labels carry outcome meaning when color and Unicode are unavailable.
- Prompt descriptions explain consequences or provide context; they do not repeat the label.
- Errors identify the summary, effect, next action, and stable reference when those values exist.
- Tables hide sensitive fields and adapt to width. Machine-readable output never contains styling.
- The plain and accessible renderers preserve information and stable wording, not visual imitation.
- No gradients, banners, boxed logos, decorative animation, or novelty glyphs.

## Accessibility baseline

Preserve keyboard navigation and screen-reader-compatible prompts. Honor `NO_COLOR`, `--color`, `--unicode`, `--no-input`, redirected streams, and terminal capability detection. Never use color, alignment, or animation as the only carrier of meaning. Test the accessible flow independently from the rich terminal flow.
