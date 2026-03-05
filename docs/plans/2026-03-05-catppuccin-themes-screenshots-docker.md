# Design: Catppuccin Themes, Screenshots, Docker Publish, README Updates

Date: 2026-03-05

## Goals

1. Add Catppuccin Latte (light) and Catppuccin Mocha (dark) theme variants
2. Capture screenshots in all 4 themes for the README
3. Publish the Docker image to ghcr.io/paulofilip3/interloki
4. Update README: inspiration section (logdy + Grafana Loki) + screenshot grid

---

## 1. Theme System

### State model

Add a second independent axis to the settings store alongside the existing `theme` (base):

```
palette: 'default' | 'catppuccin'   // new
theme:   'light'   | 'dark'         // existing, unchanged
```

The effective `data-theme` attribute is derived from both:

| theme | palette | data-theme |
|-------|---------|------------|
| light | default | `light` |
| dark  | default | `dark` |
| light | catppuccin | `catppuccin-latte` |
| dark  | catppuccin | `catppuccin-mocha` |

### Files changed

**`web/src/stores/settings.ts`**
- Add `palette` ref (`'default' | 'catppuccin'`, default `'default'`)
- Add `togglePalette()` function
- Include `palette` in localStorage persistence
- Update `SettingsState` interface

**`web/src/composables/useTheme.ts`**
- Compute effective theme key from `(theme, palette)`
- Watch both `theme` and `palette` and re-apply `data-theme`

**`web/src/styles/variables.css`**
- Add `[data-theme="catppuccin-latte"]` block (Catppuccin Latte palette)
- Add `[data-theme="catppuccin-mocha"]` block (Catppuccin Mocha palette)

Catppuccin variable mapping:

| CSS var | Latte | Mocha |
|---------|-------|-------|
| `--interloki-bg` | `#eff1f5` (Base) | `#1e1e2e` (Base) |
| `--interloki-bg-secondary` | `#e6e9ef` (Mantle) | `#181825` (Mantle) |
| `--interloki-fg` | `#4c4f69` (Text) | `#cdd6f4` (Text) |
| `--interloki-fg-secondary` | `#6c6f85` (Subtext0) | `#a6adc8` (Subtext0) |
| `--interloki-accent` | `#1e66f5` (Blue) | `#89b4fa` (Blue) |
| `--interloki-border` | `#ccd0da` (Surface0) | `#313244` (Surface0) |
| `--interloki-bg-hover` | `#ccd0da` (Surface0) | `#313244` (Surface0) |
| `--interloki-bg-active` | `#bcc0cc` (Surface1) | `#45475a` (Surface1) |
| `--interloki-shadow` | `rgba(76,79,105,0.1)` | `rgba(0,0,0,0.4)` |
| `--interloki-text-accent` | `#7287fd` (Lavender) | `#b4befe` (Lavender) |
| `--interloki-level-debug` | `#9ca0b0` (Overlay0) | `#6c7086` (Overlay0) |
| `--interloki-level-info` | `#1e66f5` (Blue) | `#89b4fa` (Blue) |
| `--interloki-level-warn` | `#df8e1d` (Yellow) | `#f9e2af` (Yellow) |
| `--interloki-level-error` | `#d20f39` (Red) | `#f38ba8` (Red) |
| `--interloki-level-fatal` | `#fe640b` (Peach) | `#fab387` (Peach) |

### UI: PaletteToggle component

New `web/src/components/PaletteToggle.vue`:
- A small pill button next to `ThemeToggle` in the header
- Shows label: `Default` or `Catppuccin`
- Calls `settings.togglePalette()` on click
- Styled consistently with `ThemeToggle` (border, hover accent)

`App.vue`: import and place `<PaletteToggle />` next to `<ThemeToggle />`.

---

## 2. Screenshots

- Directory: `screenshots/` at repo root
- Files: `light.png`, `dark.png`, `catppuccin-latte.png`, `catppuccin-mocha.png`
- Captured via browser automation with demo mode running (`interloki demo --rate=20`)
- Window size: 1280×800

README: replace `<!-- TODO: Add screenshot -->` with a 2×2 Markdown image grid:
```markdown
| Light | Dark |
|-------|------|
| ![Light](screenshots/light.png) | ![Dark](screenshots/dark.png) |
| ![Catppuccin Latte](screenshots/catppuccin-latte.png) | ![Catppuccin Mocha](screenshots/catppuccin-mocha.png) |
```

---

## 3. Docker Publish

Build multi-arch image and push to GitHub Container Registry:

```bash
docker login ghcr.io -u paulofilip3
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t ghcr.io/paulofilip3/interloki:latest \
  -t ghcr.io/paulofilip3/interloki:0.1.0 \
  --push .
```

The Dockerfile already exists in the repo root.

---

## 4. README: Inspiration Section

Insert after the Features list, before Quick Start:

```markdown
## Inspiration

interloki was inspired by two projects:

- **[logdy](https://github.com/logdyhq/logdy-core)** — the single-binary, embedded-frontend, WebSocket-streaming architecture is directly modelled on logdy-core. If you want a more mature tool in the same space, check it out.
- **[Grafana Loki](https://grafana.com/oss/loki/)** — the label-based log aggregation concept and the project name. interloki is *not* a Loki datasource; it is a lightweight alternative for teams that do not run a full Loki stack.
```
