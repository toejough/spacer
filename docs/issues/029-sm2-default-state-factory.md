# 029 — SM-2 default state factory

**Status:** closed
**Type:** tech-debt
**Source:** #014 architecture premortem

## Problem

SM-2 initial state (`easeFactor: 2.5, interval: 0, repetitions: 0, nextReview: new Date()`) is hardcoded in `DeckView.vue:24-30` and repeated in `flow.test.ts:27-34`. Adding card-creation paths (import, bulk add, API sync) multiplies the duplication. Changing defaults requires hunting every occurrence.

## Recommendation

Add a `newSM2State()` factory function to `sm2.ts` that returns the canonical initial state. All card-creation sites use it.

## Research

- SM-2 algorithm reference — review whether 2.5 is the standard initial ease factor or whether other implementations use different defaults (Anki uses 2.5 but allows user override).
- Look at how Anki's source handles [default deck config](https://github.com/ankitects/anki) for initial ease/interval settings to inform whether the factory should accept optional overrides.
