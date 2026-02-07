## Context

The Timer project is a Rust binary with a `CountdownTimer` module in `src/timer.rs`.
Currently there is no window, no frontend, and no Tauri integration — just a console
`println!`. The `DEPENDENCIES.md` specifies a Tauri 2 + Svelte 5 + Vite stack.

Tauri 2 projects have a specific directory layout: Rust code lives under `src-tauri/`,
the frontend lives at the project root (`src/`, `index.html`, `package.json`), and
`tauri.conf.json` + `capabilities/` live inside `src-tauri/`.

This means our existing `src/main.rs` and `src/timer.rs` need to move into `src-tauri/src/`.

## Goals / Non-Goals

**Goals:**
- Standard Tauri 2 project layout that works with `cargo tauri dev` and `cargo tauri build`
- Svelte 5 + Vite frontend with a minimal placeholder page
- Existing `timer.rs` module preserved and accessible from the Tauri backend
- App launches a window titled "Timer" at a reasonable default size

**Non-Goals:**
- Timer UI components (separate change)
- Tauri commands to expose timer logic to the frontend (separate change)
- Plugins beyond the default opener (dialog, fs, sql come later)
- Mobile support configuration
- App icons (placeholder/defaults are fine)

## Decisions

### Decision 1: Plain Svelte 5 + Vite (not SvelteKit)

**Approach:** Use a plain Svelte 5 + Vite setup without SvelteKit.

**Why not SvelteKit?** A timer app is a single-page tool — no routing, no SSR, no
adapter-static complexity. Plain Svelte + Vite is simpler, faster to build, and has
fewer moving parts. SvelteKit can be introduced later if the app grows.

### Decision 2: Move Rust code into `src-tauri/`

**Approach:** Relocate `Cargo.toml`, `src/main.rs`, and `src/timer.rs` into the
`src-tauri/` directory to match the canonical Tauri 2 layout.

- `src-tauri/Cargo.toml` — with Tauri dependencies added
- `src-tauri/src/main.rs` — minimal desktop entry point calling `lib::run()`
- `src-tauri/src/lib.rs` — Tauri builder setup, module declarations
- `src-tauri/src/timer.rs` — unchanged timer module
- `src-tauri/build.rs` — Tauri build codegen

**Why not keep Rust at the root?** Tauri 2 expects `tauri.conf.json` adjacent to
`Cargo.toml` inside `src-tauri/`. Fighting this convention causes build issues and
makes the project unfamiliar to anyone who knows Tauri.

### Decision 3: Split main.rs into main.rs + lib.rs

**Approach:** Follow the Tauri 2 convention of a thin `main.rs` that calls
`lib::run()`, with all app logic in `lib.rs`. This supports potential mobile
targets in the future.

- `main.rs`: just `#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]`
  and calls `run()`
- `lib.rs`: `tauri::Builder`, plugin registration, module declarations

### Decision 4: Minimal frontend placeholder

**Approach:** Create a bare-bones Svelte app (`App.svelte` showing "Timer" heading)
with `index.html`, `main.ts`, `vite.config.ts`, `tsconfig.json`, and `package.json`.

No timer UI yet — just proof that the Svelte → Tauri pipeline works.

### Decision 5: Default window configuration

**Approach:** Single window, title "Timer", 400x600 (portrait orientation suits a
countdown timer display). Resizable.

## Risks / Trade-offs

- [Directory restructure] → The existing `Cargo.toml` and `src/` at root will be
  replaced. The root `Cargo.toml` must be removed or replaced with a workspace config.
  Since there are no commits yet, this is low risk.
- [Node.js required] → The frontend needs npm installed. This is standard for Tauri
  projects but is a new system dependency.
- [Edition 2024 compatibility] → Tauri 2 crates target edition 2021. Our code uses
  edition 2024. This should be fine — editions are per-crate, not per-workspace.
  If issues arise, we can fall back to 2021.
