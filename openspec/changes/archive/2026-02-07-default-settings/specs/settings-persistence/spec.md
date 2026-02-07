## ADDED Requirements

### Requirement: Settings storage via Tauri store
The application SHALL persist user settings to disk using the Tauri store plugin.

#### Scenario: Settings saved to store
- **WHEN** the user saves settings with default start time of 60 minutes and increment of 10 minutes
- **THEN** the values `default_duration_secs: 3600` and `default_increment_secs: 600` SHALL be written to the store

#### Scenario: Settings loaded from store on startup
- **WHEN** the application starts and a settings store file exists
- **THEN** the saved `default_duration_secs` and `default_increment_secs` SHALL be loaded into managed state

### Requirement: Settings fallback to defaults
The application SHALL use hardcoded defaults when no saved settings exist.

#### Scenario: First launch with no store file
- **WHEN** the application starts and no settings store file exists
- **THEN** `default_duration_secs` SHALL be 300 and `default_increment_secs` SHALL be 300

#### Scenario: Corrupted or invalid store values
- **WHEN** the application starts and the store file contains invalid or out-of-range values
- **THEN** the application SHALL fall back to the default values of 300 seconds each

### Requirement: Get settings command
The application SHALL expose a `get_settings` Tauri command that returns the current settings.

#### Scenario: Get settings returns current values
- **WHEN** the frontend invokes `get_settings`
- **THEN** the response SHALL include `default_duration_secs` and `default_increment_secs` as numbers

### Requirement: Save settings command
The application SHALL expose a `save_settings` Tauri command that validates, persists, and updates settings.

#### Scenario: Save valid settings
- **WHEN** the frontend invokes `save_settings` with `default_duration_secs: 3600` and `default_increment_secs: 600`
- **THEN** the values SHALL be persisted to the store and the managed state SHALL be updated

#### Scenario: Save settings with out-of-range values
- **WHEN** the frontend invokes `save_settings` with `default_duration_secs: 0`
- **THEN** the command SHALL return an error and the settings SHALL not be changed

### Requirement: Store plugin permissions
The application SHALL declare the required Tauri store plugin permissions in capabilities.

#### Scenario: Store permissions configured
- **WHEN** the application starts
- **THEN** the store plugin permissions SHALL be declared in the Tauri capabilities configuration
