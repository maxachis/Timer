# Timer

Desktop countdown timer app built with Tauri 2 (Rust backend + Svelte 5 frontend).

## Tech Stack

- **Frontend:** Svelte 5 (runes), TypeScript, Vite 6
- **Backend:** Rust (edition 2024), Tauri 2
- **Persistence:** tauri-plugin-store (JSON settings), tauri-plugin-notification

## Project Structure

- `src/` — Frontend (TypeScript + Svelte). `App.svelte` is the main UI component.
- `src-tauri/src/` — Rust backend. `lib.rs` (Tauri commands/state), `timer.rs` (core timer logic).
- `openspec/` — Specification documents.

## Commands

```bash
npm run dev          # Vite dev server (frontend only, http://localhost:5173)
npm run tauri dev    # Full Tauri dev mode (frontend + Rust backend)
npm run build        # Build frontend
npm run tauri build  # Production build (desktop app)
```

## Rust Tests

```bash
cd src-tauri && cargo test
```

Unit tests live in `src-tauri/src/timer.rs` (inline `#[cfg(test)]` module).

## Conventions

- Tauri IPC commands are in `src-tauri/src/lib.rs`, decorated with `#[tauri::command]`
- State is managed via `Mutex<T>` passed through Tauri's state system
- Svelte 5 runes (`$state`, `$derived`) for reactive state, not legacy `$:` syntax
- Timer logic uses monotonic `Instant` for elapsed time (immune to clock drift)
- 80ms polling interval for timer status updates
