# 034 — Navigation is linear with no structure for growth

**Status:** open
**Type:** ux / prevention
**Source:** #015 UX/design premortem

## Problem

Navigation follows a single drill-down path: Home → Deck → Review. Back-navigation is a text link in DeckView and a button in Review's done state. There's no global navigation, no breadcrumbs, no way to jump between decks or access app-level features (settings, stats, search) from anywhere.

Adding any top-level feature requires bolting navigation onto the header or individual views ad-hoc. After several features, some are reachable from some places but not others, and there's no consistent pattern for where things live.

## Principle

A mobile-first PWA needs a navigation structure that accommodates growth without redesigning on every feature addition. The pattern should make it obvious where a new top-level feature's entry point goes.

## Guidance

Before implementing, review navigation patterns for mobile-first PWAs — bottom tab bars, hamburger menus, and their trade-offs. Consider how many top-level destinations are likely (decks, stats, settings, search) and pick a pattern that fits that range. Read the current codebase to understand what views and routes exist when this issue is picked up.
