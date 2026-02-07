## 1. Directory restructure

- [x] 1.1 Create `src-tauri/src/` directory structure
- [x] 1.2 Move `src/timer.rs` to `src-tauri/src/timer.rs` (unchanged)
- [x] 1.3 Remove old root `Cargo.toml` and `src/main.rs`

## 2. Rust backend (src-tauri)

- [x] 2.1 Create `src-tauri/Cargo.toml` with tauri 2, serde 1 (derive), serde_json 1, and tauri-build 2
- [x] 2.2 Create `src-tauri/build.rs` with `tauri_build::build()`
- [x] 2.3 Create `src-tauri/src/main.rs` — thin entry point with `windows_subsystem` attribute, calls `lib::run()`
- [x] 2.4 Create `src-tauri/src/lib.rs` — `mod timer;`, Tauri builder with `run()` function

## 3. Tauri configuration

- [x] 3.1 Create `src-tauri/tauri.conf.json` — app identifier, window (title "Timer", 400x600), build commands pointing to Vite
- [x] 3.2 Create `src-tauri/capabilities/default.json` — grant `core:default` to main window

## 4. Frontend scaffold (Svelte 5 + Vite)

- [x] 4.1 Create `package.json` with svelte 5, vite, @sveltejs/vite-plugin-svelte, @tauri-apps/api 2, @tauri-apps/cli 2
- [x] 4.2 Create `vite.config.ts` with Svelte plugin and Tauri-compatible server settings
- [x] 4.3 Create `tsconfig.json` for TypeScript configuration
- [x] 4.4 Create `index.html` — Vite entry point loading `src/main.ts`
- [x] 4.5 Create `src/main.ts` — imports and mounts `App.svelte`
- [x] 4.6 Create `src/App.svelte` — minimal placeholder with "Timer" heading

## 5. Verify

- [x] 5.1 Run `npm install` to install frontend dependencies
- [x] 5.2 Compile Rust backend (`cargo build` in src-tauri) to verify Tauri setup
- [x] 5.3 Run timer module tests (`cargo test` in src-tauri) to verify module preserved
