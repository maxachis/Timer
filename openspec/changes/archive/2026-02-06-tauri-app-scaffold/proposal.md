## Why

The Timer app currently runs as a plain Rust binary with no window or UI. To deliver a desktop application with a Svelte frontend, we need the Tauri 2 project structure in place — Rust backend entry point, frontend scaffolding, build configuration, and dev tooling. This is the foundation that all UI work depends on.

## What Changes

- **BREAKING**: Replace the current `main.rs` console binary with a Tauri application entry point
- Add Tauri 2 and Serde dependencies to `Cargo.toml`, plus `tauri-build` as a build dependency
- Create `build.rs` for Tauri's build-time code generation
- Create `src-tauri/tauri.conf.json` with window configuration (title: "Timer", default size)
- Scaffold the Svelte 5 + Vite frontend: `package.json`, `vite.config.ts`, `tsconfig.json`, `src/` with a minimal `App.svelte` and `index.html`
- Add Tauri capability/permission files as required by Tauri 2's security model

## Capabilities

### New Capabilities
- `tauri-app-shell`: Tauri 2 application bootstrap — window creation, frontend dev server integration, and build pipeline

### Modified Capabilities
<!-- None — existing timer module is preserved, just the binary entry point changes -->

## Impact

- `Cargo.toml`: Add `tauri`, `serde` dependencies; add `[build-dependencies]` section with `tauri-build`
- `src/main.rs`: Replace `println` with Tauri application builder
- `build.rs`: New file for Tauri build codegen
- `src-tauri/tauri.conf.json`: New Tauri configuration
- `src-tauri/capabilities/`: New Tauri 2 permission files
- `package.json`, `vite.config.ts`, `tsconfig.json`: New frontend tooling config
- `src/` (frontend): `index.html`, `main.ts`, `App.svelte` — minimal Svelte shell
- `src/timer.rs`: Unchanged (preserved as-is)
