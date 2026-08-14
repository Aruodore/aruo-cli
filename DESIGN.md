# Design Direction

## Intent

Generated applications are working product surfaces, not marketing demos. Their initial screen exists to prove the UI system is usable and accessible without inventing a business domain. It should feel calm, precise, and ready to be replaced feature by feature.

## Experience

- Register: product.
- Density: moderate; compact controls, comfortable reading measure.
- Visual character: restrained neutral surfaces, one accessible accent, visible state changes.
- Typography: system sans-serif stack, fixed type scale, no decorative display face.
- Motion: 150–200 ms and only for state feedback; respect reduced-motion preferences.
- Responsive behavior: structural layout changes rather than fluid type tricks.

## Interface contract

- Use semantic landmarks, native controls, explicit labels, and visible keyboard focus.
- Provide loading, empty, invalid, unavailable, unauthorized, and unexpected states when a feature can reach them.
- Keep controls visually consistent across screens; do not create one-off button or form vocabularies.
- Use color as reinforcement, never as the only state signal.
- Avoid fake metrics, fictional records, generic dashboards, testimonials, pricing, and decorative product claims.

## Accessibility baseline

Target WCAG 2.2 AA. Preserve keyboard navigation, readable status announcements, sufficient contrast, 44px coarse-pointer targets where practical, and reduced-motion behavior. Automated checks support review but do not replace keyboard and screen-reader testing.

