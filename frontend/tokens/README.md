Tokens for cross-platform design system

- `tokens.json` contains color, spacing, fontSize and shadow primitives derived from `src/theme.js`.
- `components.json` maps component-level tokens for buttons, inputs, avatars and cards.

Usage
- Consume `tokens.json` in Flutter to configure `ThemeData` colors, spacing constants, and text styles.
- Map `components.json` to widget styles for consistent visual parity.

Next steps
- Add a small Dart generator that converts `tokens.json` into Flutter `ThemeData` and constants.
- Keep tokens as the single source of truth; update `src/theme.js` and regenerate tokens when visual changes are made.
