# 036 — No consistent pattern for empty/error/loading states

**Status:** closed
**Type:** ux / prevention
**Source:** #015 UX/design premortem

## Problem

HomeView and DeckView have ad-hoc empty states ("No decks yet", "No cards yet"), but there's no loading state, no error state, and invalid routes render blank (DeckView's `v-if="deck"` shows nothing when the deck doesn't exist). Each view handles state display differently.

After 10+ features, each new view invents its own approach to loading spinners, empty illustrations, error messages, and invalid-input handling. The inconsistency is disorienting — some screens show spinners, some flash blank, some show stale data.

## Principle

Every data-driven view goes through a lifecycle: loading → loaded (with content or empty) → possibly error. A consistent pattern for representing these states prevents each feature from reinventing the wheel and gives users predictable visual cues.

## Guidance

Before implementing, review patterns for view state machines in Vue apps — how to represent loading/error/empty/loaded as explicit states rather than combinations of boolean refs. Also look at mobile UX conventions for skeleton screens vs. spinners vs. progressive loading. Read the current codebase to understand what state patterns have emerged since this issue was filed.
