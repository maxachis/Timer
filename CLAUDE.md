# Timer

Desktop countdown timer app. **Active implementation:** Wails 2 (Go backend + Svelte 5 frontend) in `wails-app/`.

> **Deprecated:** The original Tauri 2 (Rust) implementation in `src-tauri/` and the root-level `src/` frontend are deprecated and kept only for reference. Do not add new features there. See `src-tauri/DEPRECATED.md`.

## Tech Stack (active — `wails-app/`)

- **Frontend:** Svelte 5 (runes), TypeScript, Vite
- **Backend:** Go, Wails 2

## Project Structure

- `wails-app/` — **Active app.** Go/Wails backend + Svelte frontend in `wails-app/frontend/`.
- `src-tauri/` — *Deprecated* Rust/Tauri backend.
- `src/` — *Deprecated* root-level Svelte frontend (Tauri-era).
- `openspec/` — Specification documents.

## Commands

Run from `wails-app/`:

```bash
wails dev -tags webkit2_41     # Full dev mode (Go backend + Vite frontend)
wails build -tags webkit2_41   # Production build
```

The `-tags webkit2_41` flag is required on Ubuntu 24.04+, which ships `libwebkit2gtk-4.1-dev` instead of the 4.0 series Wails defaults to. Without the tag, the build fails with `Package 'webkit2gtk-4.0' ... not found` from pkg-config.

## Conventions

- Svelte 5 runes (`$state`, `$derived`) for reactive state, not legacy `$:` syntax
- Timer logic uses monotonic time for elapsed calculation (immune to clock drift)

## Mistakes

- **[tooling]**: On Ubuntu 24.04+ only `libwebkit2gtk-4.1-dev` is available, not 4.0. `wails dev` → `wails dev -tags webkit2_41`.
