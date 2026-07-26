# Prevent abandoning completed todos

## Why

Abandoning should be an alternative to completing a todo. Allowing abandonment for already-completed todos is confusing and can suggest irreversible deletion. The UX should prevent abandoning a done todo and inform the user.

## What Changes

- Disallow the Abandon action for todos with done === 1.
- Show an informative message when the user attempts to abandon a completed todo.
- Update the todo-list spec to require this behavior.

