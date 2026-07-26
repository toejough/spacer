# Add Abandon action to Edit modal

## Why

The spec requires that items can be Abandoned from the edit modal, but the UI did not provide a control in the modal. Users expect to perform the Abandon action from the item edit view as an alternative to the list-level control.

## What Changes

- Add an Abandon control (button) to the edit modal for todos and notes.
- Hide/disable the Abandon control when a todo is already completed (done === 1).
- Confirm Abandon and close the modal on success.

## Impact

- Frontend only (srv/static/script.js, srv/templates/index.html)
- Tests updated: srv/static/script.test.js

