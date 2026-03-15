# 035 — Review quality ratings are unlabeled

**Status:** open
**Type:** ux / prevention
**Source:** #015 UX/design premortem

## Problem

The review screen presents rating buttons as bare numbers `1 2 3 4 5`. SM-2 quality scores aren't intuitive — users don't know what 1 vs. 3 vs. 5 means, whether higher is better or worse, or how their choice affects scheduling. This is confusing now and gets worse if alternative review modes with different scales are added.

## Principle

Interactive controls should communicate their meaning without requiring users to learn an external system. The labels should convey the effect on the user's experience ("Again", "Hard", "Good", "Easy") rather than exposing the algorithm's internal scale.

## Guidance

Before implementing, look at how Anki and other spaced repetition apps label their rating buttons — both the labels used and how many choices they present. Consider whether all 5 SM-2 quality levels need to be exposed or whether a reduced set (e.g., Again/Hard/Good/Easy mapping to SM-2 scores) is better for the user. Read the current review flow to understand what's changed since this issue was filed.
