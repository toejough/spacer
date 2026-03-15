# File Format Templates

## Issue File

`docs/issues/{number}-{slug}.md`:

```markdown
# {Title}

**Status:** open | in-progress | done | blocked | wont-fix
**Priority:** p0 | p1 | p2
**Labels:** feature | bug | tooling | docs | infra | issue-N-retro
**Created:** {date}
**Closed:** {date}

## Description
{What and why}

## Acceptance Criteria
- [ ] {Criterion 1}
- [ ] {Criterion 2}

## Notes
{Anything learned during attempts}
```

## Status File

`docs/status.md`:

```markdown
# Spacer — Project Status

**Last updated:** {date+time}
**Current increment:** {number}
**Streak:** {N consecutive successful increments}

## Done
- #{number} {title} ({date})

## In Progress
- #{number} {title} — {current phase}

## Up Next
- #{number} {title}

## Blocked
- #{number} {title} — {reason}
```

## Retro Entry

Append to `docs/retros.md`:

```markdown
### Increment #{number}: {title} — {pass|fail}
**What worked:** {1 sentence}
**What to improve:** {1 sentence}
**Action items:**
- {Concrete action → issue #{number} or skill change}
```
