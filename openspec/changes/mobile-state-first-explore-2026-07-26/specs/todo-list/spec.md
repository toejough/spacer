## ADDED Requirements

### Requirement: Mobile state-first UX
The system SHALL provide a mobile-first representation for changing todo state that is safe and discoverable.

#### Scenario: Mark Done via visible button
- **WHEN** a user views a todo in the mobile stacked swimlane list
- **THEN** the card exposes a prominent "Done" control
- **AND** tapping Done sets completed_at to now and shows an undo snackbar

#### Scenario: Abandon via menu with confirm
- **WHEN** a user needs to Abandon an item on mobile
- **THEN** they can open an overflow menu (or long-press) and select Abandon
- **AND** the system asks for confirm or provides an undo snackbar before permanently archiving

#### Scenario: Swipe to Done (optional variant)
- **WHEN** swipe-right-to-Done is enabled
- **THEN** a swipe sets completed_at to now and shows an undo snackbar
- **AND** no swipe action shall perform Abandon

#### Scenario: Bulk selection
- **WHEN** the user enters multi-select mode
- **THEN** they can apply Done or Abandon to multiple items and confirm Abandon actions

