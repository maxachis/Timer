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
wails dev     # Full dev mode (Go backend + Vite frontend)
wails build   # Production build
```

## Conventions

- Svelte 5 runes (`$state`, `$derived`) for reactive state, not legacy `$:` syntax
- Timer logic uses monotonic time for elapsed calculation (immune to clock drift)
