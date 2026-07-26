## ADDED Requirements

### Requirement: Mobile parity for Done and Abandon
The system SHALL make both Done and Abandon equally accessible on mobile while providing safeguards for Abandon.

#### Scenario: Two visible actions
- **WHEN** a user views a todo card on mobile
- **THEN** the card exposes two visible, focusable controls: Done and Abandon
- **AND** both controls are reachable within 2 taps
- **AND** activating Abandon shows either a confirm or an undo snackbar depending on sensitivity

#### Scenario: Split-swipe symmetric gestures
- **WHEN** split-swipe is enabled
- **THEN** swipe-right marks Done and swipe-left marks Abandon
- **AND** both actions show undo snackbar; Abandon may additionally require confirmation when enabled
- **AND** keyboard/menu alternatives exist for both actions

#### Scenario: Accessibility parity
- **WHEN** a screen reader user navigates a card
- **THEN** the Done and Abandon controls are announced and reachable via standard keyboard focus
- **AND** ARIA live regions announce successful state changes and undos

