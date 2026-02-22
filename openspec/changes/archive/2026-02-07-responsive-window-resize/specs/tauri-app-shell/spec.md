## MODIFIED Requirements

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
