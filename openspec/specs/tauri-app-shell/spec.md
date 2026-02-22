# tauri-app-shell Specification

## Purpose
Tauri 2 application bootstrap — window creation, frontend dev server integration, and build pipeline.

## Requirements

### Requirement: Tauri project structure
The project SHALL follow the canonical Tauri 2 directory layout with Rust code under `src-tauri/` and frontend code at the project root.

#### Scenario: Standard directory layout
- **WHEN** the project root is inspected
- **THEN** `src-tauri/` SHALL contain `Cargo.toml`, `build.rs`, `tauri.conf.json`, `capabilities/`, and `src/`
- **AND** the project root SHALL contain `package.json`, `vite.config.ts`, `index.html`, and a `src/` directory for frontend code

### Requirement: Tauri backend entry point
The Rust backend SHALL use the Tauri 2 builder pattern with a thin `main.rs` and shared `lib.rs`.

#### Scenario: Application startup
- **WHEN** the application is launched via `cargo tauri dev`
- **THEN** `main.rs` SHALL call the `run()` function from `lib.rs`
- **AND** `lib.rs` SHALL initialize `tauri::Builder` and run the application

#### Scenario: Windows console suppression
- **WHEN** the application is built in release mode on Windows
- **THEN** the console window SHALL NOT appear (via `windows_subsystem = "windows"` attribute)

### Requirement: Tauri build configuration
The project SHALL include a `build.rs` that invokes `tauri_build::build()` for code generation.

#### Scenario: Build script execution
- **WHEN** the project is compiled
- **THEN** `build.rs` SHALL run `tauri_build::build()` to generate capability schemas and platform configuration

### Requirement: Tauri window configuration
The application SHALL launch a single window with appropriate defaults for a timer application.

#### Scenario: Default window properties
- **WHEN** the application launches
- **THEN** the window title SHALL be "Timer"
- **AND** the window size SHALL be 400x600 pixels
- **AND** the window SHALL be resizable

#### Scenario: Minimum window size
- **WHEN** the user resizes the window
- **THEN** the window SHALL NOT shrink below 100x100 pixels
- **AND** `minWidth` SHALL be set to 100 and `minHeight` SHALL be set to 100 in the window configuration

### Requirement: Tauri capabilities and permissions
The application SHALL declare a default capability granting core permissions to the main window.

#### Scenario: Default capability file
- **WHEN** the application starts
- **THEN** `src-tauri/capabilities/default.json` SHALL grant `core:default` permissions to the main window

### Requirement: Svelte frontend scaffold
The project SHALL include a minimal Svelte 5 + Vite frontend that renders in the Tauri window.

#### Scenario: Frontend dev server
- **WHEN** `npm run dev` is executed
- **THEN** Vite SHALL start a dev server on port 5173
- **AND** the Svelte app SHALL be served with hot module replacement

#### Scenario: Frontend production build
- **WHEN** `npm run build` is executed
- **THEN** Vite SHALL produce a static build in the `dist/` directory

#### Scenario: Minimal placeholder UI
- **WHEN** the frontend loads
- **THEN** the page SHALL display a "Timer" heading as a placeholder

### Requirement: Rust dependencies
The `src-tauri/Cargo.toml` SHALL declare the required Tauri 2 dependencies.

#### Scenario: Required crate dependencies
- **WHEN** `src-tauri/Cargo.toml` is inspected
- **THEN** it SHALL include `tauri` version 2 as a runtime dependency
- **AND** it SHALL include `serde` version 1 with the `derive` feature
- **AND** it SHALL include `tauri-build` version 2 as a build dependency

### Requirement: Existing timer module preserved
The existing `CountdownTimer` module SHALL be relocated to `src-tauri/src/timer.rs` without modification.

#### Scenario: Timer module accessible from backend
- **WHEN** `src-tauri/src/lib.rs` is compiled
- **THEN** the `timer` module SHALL be declared and accessible
- **AND** all existing timer tests SHALL pass
