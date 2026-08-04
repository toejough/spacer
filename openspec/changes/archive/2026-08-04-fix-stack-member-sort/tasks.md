## 1. Fix

- [x] 1.1 In `buildDisplayList` (`srv/static/script.js`, ~line 211), sort the stack `members` array with `compareByRelevance` before adding it to `entries` (`memberPool.filter(i => i.stack_id === stackId).sort(compareByRelevance)`)

## 2. Verify

- [x] 2.1 Run/use the app: create a stack with members in mixed states (open/done, differing review urgency, differing updated_at), expand it, and confirm member order matches the order the same cards would have unstacked
- [x] 2.2 Confirm `gcStacks`/`getStackMembers` behavior is unaffected (still an existence check, not display order)
