## MODIFIED Requirements

### Requirement: Mobile parity for Done and Abandon (IMPLEMENTED)
The system SHALL make both Done and Abandon equally accessible on mobile while providing safeguards for Abandon, and SHALL provide a Reopen action from either closed state.

#### Scenarios implemented
- Two visible actions with Reopen on closed cards.
- Split-swipe left/right is NOT part of the default implementation; swipe-based interactions are archived as experimental. The default implementation uses two visible buttons per card (Done and Abandon) and Reopen on closed cards.
- Keyboard and ARIA announcements implemented: controls are focusable and state changes use ARIA live regions.
- Visual treatments for each state provided and contrast verified.

