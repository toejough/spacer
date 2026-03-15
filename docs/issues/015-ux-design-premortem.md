# 015 — UX & Visual Design Premortem

**Status:** open
**Type:** design / prevention

## Context

Spacer is a mobile-first PWA for spaced repetition flashcards. The current UI is functional but minimal: 3 views, inline Tailwind classes, a single accent color, no shared component abstractions. This is fine at current scale.

## The Exercise

Perform a premortem on the visual and UX design. Assume we've shipped 10+ features (stats, tags, search, import/export, settings, media cards, cram mode, streak tracking, etc.) and the UX is now incoherent — inconsistent styling, confusing navigation, poor discoverability, and an app that looks like it was built by 10 different people. Users are frustrated.

**Your job:** Figure out what led to that state.

### How to run the premortem

1. **Audit the current UI** — read every view and template. Note how styles are applied (inline classes? shared components? copy-paste?), what interactive states exist (hover, focus, disabled, loading, empty, error), how navigation works, and what feedback the user gets for their actions.
2. **Imagine 10+ features added** — stats dashboard, tags/search, import/export, sync, settings, media cards, cram mode, streak tracking, sharing, onboarding, etc. Consider this as a mobile-first PWA that users interact with daily.
3. **Identify 3-5 specific UX/design weaknesses** in the current UI that would compound badly under that growth. Focus on things that are fine at current scale but would create incoherence, confusion, or frustration. Be concrete — reference actual templates, class strings, interaction patterns, and missing states.
4. **For each weakness**, describe: what specifically goes wrong as features multiply, why the current UI enables it, and a concrete mitigation (with a recommendation on whether to adopt now or defer with a specific trigger).

### What makes a good premortem item

- Rooted in what the templates and interactions actually look like today
- Specific enough to point to the exact view/element where inconsistency would start
- Considers the user's experience, not just code cleanliness
- The mitigation is proportional — lightweight for "adopt now", clearly scoped for "defer"

## Deliverable

- 3-5 premortem items with analysis and mitigations
- For each: a decision recommendation (adopt now / defer with trigger / reject)
- Any "adopt now" mitigations implemented
- Design conventions documented for future feature work

## Acceptance Criteria

- [ ] Current UI fully audited (every view, every template)
- [ ] 3-5 risks identified with concrete references to current templates/interactions
- [ ] Each risk has a mitigation with adopt/defer/reject recommendation
- [ ] Decisions recorded and any immediate mitigations implemented
