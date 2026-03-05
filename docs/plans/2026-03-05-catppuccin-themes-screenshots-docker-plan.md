# Catppuccin Themes, Screenshots & Docker Publish — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add Catppuccin Latte/Mocha theme variants with a palette toggle, capture 4 themed screenshots for the README, update README with inspiration section, and publish the Docker image to ghcr.io/paulofilip3/interloki.

**Architecture:** Theme system gains a second independent axis (`palette: default|catppuccin`) alongside the existing `theme: light|dark`. The combination drives the `data-theme` attribute. A new `PaletteToggle` pill sits beside the existing `ThemeToggle` in the header.

**Tech Stack:** Vue 3 + Pinia + CSS custom properties; Go/Docker for publish; browser automation for screenshots.

---

## Task 1: Add Catppuccin CSS variables

**Files:**
- Modify: `web/src/styles/variables.css`

**Step 1: Append the two new theme blocks**

Open `web/src/styles/variables.css` and append after the existing `[data-theme="dark"]` block:

```css
[data-theme="catppuccin-latte"] {
  --interloki-bg: #eff1f5;
  --interloki-bg-secondary: #e6e9ef;
  --interloki-fg: #4c4f69;
  --interloki-fg-secondary: #6c6f85;
  --interloki-accent: #1e66f5;
  --interloki-border: #ccd0da;
  --interloki-font-family: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  --interloki-font-size: 13px;
  --interloki-row-height: 24px;
  --interloki-level-debug: #9ca0b0;
  --interloki-level-info: #1e66f5;
  --interloki-level-warn: #df8e1d;
  --interloki-level-error: #d20f39;
  --interloki-level-fatal: #fe640b;
  --interloki-bg-hover: #ccd0da;
  --interloki-bg-active: #bcc0cc;
  --interloki-shadow: rgba(76, 79, 105, 0.1);
  --interloki-text-accent: #7287fd;
}

[data-theme="catppuccin-mocha"] {
  --interloki-bg: #1e1e2e;
  --interloki-bg-secondary: #181825;
  --interloki-fg: #cdd6f4;
  --interloki-fg-secondary: #a6adc8;
  --interloki-accent: #89b4fa;
  --interloki-border: #313244;
  --interloki-font-family: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  --interloki-font-size: 13px;
  --interloki-row-height: 24px;
  --interloki-level-debug: #6c7086;
  --interloki-level-info: #89b4fa;
  --interloki-level-warn: #f9e2af;
  --interloki-level-error: #f38ba8;
  --interloki-level-fatal: #fab387;
  --interloki-bg-hover: #313244;
  --interloki-bg-active: #45475a;
  --interloki-shadow: rgba(0, 0, 0, 0.4);
  --interloki-text-accent: #b4befe;
}
```

**Step 2: Verify visually (no automated test for pure CSS)**

After Task 5 the browser automation will validate this visually. No commit yet — batch with Task 2–5.

---

## Task 2: Extend settings store with `palette`

**Files:**
- Modify: `web/src/stores/settings.ts`

**Step 1: Update the interface and state**

The current `SettingsState` interface only has `theme: 'light' | 'dark'`. Update it and add `palette`:

```typescript
interface SettingsState {
  theme: 'light' | 'dark'
  palette: 'default' | 'catppuccin'   // ADD THIS
  autoFollow: boolean
  showTimestamp: boolean
  showLevel: boolean
  showSource: boolean
}
```

Add the ref after the existing `theme` ref:

```typescript
const theme = ref<'light' | 'dark'>(saved.theme ?? 'dark')
const palette = ref<'default' | 'catppuccin'>(saved.palette ?? 'default')  // ADD
```

**Step 2: Add `togglePalette` and update `persist`**

Update the `persist` function to include `palette.value`:

```typescript
function persist() {
  saveToStorage({
    theme: theme.value,
    palette: palette.value,    // ADD
    autoFollow: autoFollow.value,
    showTimestamp: showTimestamp.value,
    showLevel: showLevel.value,
    showSource: showSource.value,
  })
}
```

Update the `watch` call to include `palette`:

```typescript
watch([theme, palette, autoFollow, showTimestamp, showLevel, showSource], persist)
```

Add the `togglePalette` function after `toggleTheme`:

```typescript
function togglePalette() {
  palette.value = palette.value === 'default' ? 'catppuccin' : 'default'
}
```

Update the returned object:

```typescript
return {
  theme,
  palette,          // ADD
  autoFollow,
  showTimestamp,
  showLevel,
  showSource,
  toggleTheme,
  togglePalette,    // ADD
  setAutoFollow,
  setShowTimestamp,
  setShowLevel,
  setShowSource,
}
```

---

## Task 3: Update `useTheme` to compute effective theme key

**Files:**
- Modify: `web/src/composables/useTheme.ts`

**Step 1: Compute the effective `data-theme` value from both axes**

Replace the entire file content:

```typescript
import { watch, onMounted, computed } from 'vue'
import { useSettingsStore } from '../stores/settings'

export function useTheme() {
  const settingsStore = useSettingsStore()

  const effectiveTheme = computed(() => {
    if (settingsStore.palette === 'catppuccin') {
      return settingsStore.theme === 'dark' ? 'catppuccin-mocha' : 'catppuccin-latte'
    }
    return settingsStore.theme
  })

  function applyTheme(key: string) {
    document.documentElement.setAttribute('data-theme', key)
  }

  onMounted(() => {
    applyTheme(effectiveTheme.value)
  })

  watch(effectiveTheme, (key) => {
    applyTheme(key)
  })

  return {
    theme: settingsStore.theme,
    toggleTheme: settingsStore.toggleTheme,
    palette: settingsStore.palette,
    togglePalette: settingsStore.togglePalette,
  }
}
```

---

## Task 4: Create `PaletteToggle` component

**Files:**
- Create: `web/src/components/PaletteToggle.vue`

**Step 1: Write the component**

```vue
<script setup lang="ts">
import { useSettingsStore } from '../stores/settings'

const settings = useSettingsStore()
</script>

<template>
  <button
    class="palette-toggle"
    :class="{ 'palette-toggle--active': settings.palette === 'catppuccin' }"
    @click="settings.togglePalette()"
    :title="settings.palette === 'catppuccin' ? 'Switch to default palette' : 'Switch to Catppuccin palette'"
  >
    {{ settings.palette === 'catppuccin' ? 'Catppuccin' : 'Default' }}
  </button>
</template>

<style scoped>
.palette-toggle {
  background: none;
  border: 1px solid var(--interloki-border);
  border-radius: 4px;
  cursor: pointer;
  padding: 3px 8px;
  font-size: 11px;
  font-family: var(--interloki-font-family);
  color: var(--interloki-fg-secondary);
  transition: border-color 0.15s, color 0.15s, background-color 0.15s;
  white-space: nowrap;
}

.palette-toggle:hover {
  border-color: var(--interloki-accent);
  color: var(--interloki-fg);
}

.palette-toggle--active {
  border-color: var(--interloki-accent);
  color: var(--interloki-accent);
}
</style>
```

---

## Task 5: Wire `PaletteToggle` into `App.vue`

**Files:**
- Modify: `web/src/App.vue`

**Step 1: Import and place the component**

Add import after the `ThemeToggle` import:

```typescript
import PaletteToggle from './components/PaletteToggle.vue'
```

In the template, add `<PaletteToggle />` beside `<ThemeToggle />`:

```html
<div class="interloki-header__actions">
  <ColumnConfig />
  <PaletteToggle />
  <ThemeToggle />
</div>
```

**Step 2: Build frontend and verify**

```bash
cd /home/paulo/OwnCloud/developaulo/personal/work/interpt/tech/devs/interloki
make build-frontend
```

Expected: no errors, `web/dist/` updated.

**Step 3: Commit the theme changes**

```bash
cd /home/paulo/OwnCloud/developaulo/personal/work/interpt/tech/devs/interloki
git add web/src/styles/variables.css \
        web/src/stores/settings.ts \
        web/src/composables/useTheme.ts \
        web/src/components/PaletteToggle.vue \
        web/src/App.vue
git commit -m "feat: add Catppuccin Latte/Mocha themes with palette toggle"
```

---

## Task 6: Take screenshots in all 4 themes

**Files:**
- Create: `screenshots/` directory at repo root
- Create: `screenshots/light.png`, `screenshots/dark.png`, `screenshots/catppuccin-latte.png`, `screenshots/catppuccin-mocha.png`

**Step 1: Build the full binary with the updated frontend**

```bash
cd /home/paulo/OwnCloud/developaulo/personal/work/interpt/tech/devs/interloki
make build
```

**Step 2: Start demo in background**

```bash
./bin/interloki demo --rate=20 &
# Wait a couple seconds for it to start, then note the PID
```

Open `http://localhost:8080` in the browser via browser automation.

**Step 3: Capture screenshots via browser automation**

Use `mcp__claude-in-chrome__*` tools:
1. Navigate to `http://localhost:8080`
2. Resize window to 1280×800
3. Let demo logs fill for 2-3 seconds
4. Screenshot → `screenshots/dark.png` (default state is dark+default)
5. Click ThemeToggle → now light+default → screenshot → `screenshots/light.png`
6. Click PaletteToggle → now light+catppuccin → screenshot → `screenshots/catppuccin-latte.png`
7. Click ThemeToggle → now dark+catppuccin → screenshot → `screenshots/catppuccin-mocha.png`

**Step 4: Stop the demo backend**

```bash
kill %1
```

**Step 5: Commit screenshots**

```bash
cd /home/paulo/OwnCloud/developaulo/personal/work/interpt/tech/devs/interloki
git add screenshots/
git commit -m "docs: add screenshots for all 4 themes"
```

---

## Task 7: Update README

**Files:**
- Modify: `README.md`

**Step 1: Replace the screenshot placeholder with a 2×2 grid**

Find line 9: `<!-- TODO: Add screenshot -->`

Replace with:

```markdown
| Light | Dark |
|-------|------|
| ![Light theme](screenshots/light.png) | ![Dark theme](screenshots/dark.png) |
| ![Catppuccin Latte](screenshots/catppuccin-latte.png) | ![Catppuccin Mocha](screenshots/catppuccin-mocha.png) |
```

**Step 2: Add Inspiration section**

After the Features section (after the last `- **Single binary**` bullet and before `## Quick Start`), insert:

```markdown
## Inspiration

interloki was inspired by two projects:

- **[logdy](https://github.com/logdyhq/logdy-core)** — the single-binary, embedded-frontend, WebSocket-streaming architecture is directly modelled on logdy-core. If you want a more mature tool in the same space, check it out.
- **[Grafana Loki](https://grafana.com/oss/loki/)** — the label-based log aggregation concept and the project name. interloki is *not* a Loki datasource; it is a lightweight alternative for teams that do not run a full Loki stack.

```

**Step 3: Add Catppuccin themes to the Features list**

Find the line:
```
- **Light/dark themes** — switchable via the UI, driven by CSS custom properties
```

Replace with:
```
- **Light/dark themes with Catppuccin variants** — four palette options (Default Light, Default Dark, Catppuccin Latte, Catppuccin Mocha) switchable via the UI, driven by CSS custom properties
```

**Step 4: Commit**

```bash
cd /home/paulo/OwnCloud/developaulo/personal/work/interpt/tech/devs/interloki
git add README.md
git commit -m "docs: add inspiration section, screenshot grid, and Catppuccin theme mention"
```

---

## Task 8: Publish Docker image

**Prerequisites:** Must be logged in to ghcr.io. Run:

```bash
docker login ghcr.io -u paulofilip3
# Enter a GitHub personal access token with write:packages scope when prompted
```

**Step 1: Ensure buildx builder exists**

```bash
docker buildx ls
# If no multi-platform builder exists:
docker buildx create --use --name multiarch
```

**Step 2: Build and push multi-arch image**

Run from the repo root:

```bash
cd /home/paulo/OwnCloud/developaulo/personal/work/interpt/tech/devs/interloki
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t ghcr.io/paulofilip3/interloki:latest \
  -t ghcr.io/paulofilip3/interloki:0.1.0 \
  --push .
```

Expected: both architectures built and pushed. The Dockerfile already handles multi-arch (all base images are multi-arch).

**Step 3: Verify**

```bash
docker buildx imagetools inspect ghcr.io/paulofilip3/interloki:latest
```

Expected: shows `linux/amd64` and `linux/arm64` manifests.

**Step 4: Smoke-test the published image**

```bash
docker run --rm -p 8181:8080 ghcr.io/paulofilip3/interloki:latest demo --rate=10
# Open http://localhost:8181 and confirm it works, then ctrl-c
```

---

## Task 9: Commit design doc

```bash
cd /home/paulo/OwnCloud/developaulo/personal/work/interpt/tech/devs/interloki
git add docs/plans/
git commit -m "docs: add catppuccin/screenshots/docker design and implementation plan"
```
