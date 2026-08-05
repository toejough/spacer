## Purpose

Governs how the app represents icons — branding, navigation, actions, and the browser/app icon — so the served app presents one consistent, theme-aware icon system instead of platform-dependent emoji glyphs.

## ADDED Requirements

### Requirement: All in-app iconography uses the SVG icon system, not emoji

The app SHALL represent every icon — branding, navigation, buttons, and status indicators — as an inline SVG in the same visual language (outline/stroke style, theme-color-aware via `currentColor`) already used for todo-card actions. No emoji character SHALL be used as an icon anywhere in the served app.

#### Scenario: Note cards match todo cards
- **WHEN** a note card is rendered
- **THEN** its abandon and edit buttons use the same SVG icons the equivalent todo-card buttons use, not emoji

#### Scenario: Stack controls use the icon system
- **WHEN** a stack's rename control is rendered
- **THEN** it uses an SVG icon, not an emoji

#### Scenario: Branding and navigation use the icon system
- **WHEN** the header, the four tab buttons, the help button, the empty-review state, or any heading in the help content is rendered
- **THEN** each uses an SVG icon from the same system, not an emoji

#### Scenario: Structural affordances use the icon system
- **WHEN** a stack's expand/collapse control or an item's drag handle is rendered
- **THEN** each uses an SVG icon, not a Unicode symbol standing in for one

### Requirement: The app has a dedicated favicon and app icon

The app SHALL serve a dedicated icon asset for the browser tab and PWA install, distinct from generating one by drawing an emoji character as SVG text.

#### Scenario: The browser tab shows a real icon
- **WHEN** the app's entry page is loaded
- **THEN** the response includes an icon link the browser can use for the tab favicon

#### Scenario: The PWA install icon is a real icon
- **WHEN** the web app manifest is requested
- **THEN** its icon is a dedicated icon asset, not an emoji character rendered as SVG text
