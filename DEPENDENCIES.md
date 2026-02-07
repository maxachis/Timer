# Dependencies

## Frontend (npm)

### Runtime Dependencies

| Package | Version | Purpose |
|---|---|---|
| `@tauri-apps/api` | ^2.0.0 | Core Tauri API for frontend-backend IPC (`invoke`) |
| `@tauri-apps/plugin-dialog` | ^2.0.0 | Native OS dialogs (file open/save, alerts) |
| `@tauri-apps/plugin-fs` | ^2.0.0 | Filesystem access from the frontend |
| `@tauri-apps/plugin-sql` | ^2.0.0 | SQL database access (SQLite) from the frontend |

### Dev Dependencies

| Package | Version | Purpose |
|---|---|---|
| `svelte` | ^5.0.0 | UI framework (Svelte 5 with runes) |
| `@sveltejs/vite-plugin-svelte` | ^4.0.0 | Vite plugin for Svelte compilation |
| `vite` | ^5.4.0 | Build tool and dev server |
| `typescript` | ^5.0.0 | TypeScript compiler |
| `svelte-check` | ^4.0.0 | Svelte type checking and diagnostics |
| `@tsconfig/svelte` | ^5.0.4 | Shared TypeScript config base for Svelte |
| `tslib` | ^2.6.0 | TypeScript runtime helpers |
| `@tauri-apps/cli` | ^2.0.0 | Tauri CLI for building and running the app |
| `vitest` | ^4.0.18 | Test runner |
| `@testing-library/svelte` | ^5.3.1 | Svelte component testing utilities |
| `jsdom` | ^28.0.0 | DOM environment for tests |

## Backend (Cargo / Rust)

### Runtime Dependencies

| Crate | Version | Features | Purpose |
|---|---|---|---|
| `tauri` | 2 | — | Core Tauri runtime and command system |
| `tauri-plugin-sql` | 2 | `sqlite` | SQLite database plugin |
| `tauri-plugin-dialog` | 2 | — | Native OS dialog plugin |
| `tauri-plugin-fs` | 2 | — | Filesystem plugin |
| `serde` | 1 | `derive` | Serialization/deserialization of data structures |
| `serde_json` | 1 | — | JSON parsing and generation |
| `uuid` | 1 | `v4` | UUID v4 generation for book IDs |
| `chrono` | 0.4 | `serde` | Date/time handling with serde support |
| `csv` | 1.3 | — | CSV parsing and writing for import/export |

### Build Dependencies

| Crate | Version | Purpose |
|---|---|---|
| `tauri-build` | 2 | Tauri build-time code generation |
