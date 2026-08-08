const fs = require('fs');
const path = require('path');
const vm = require('vm');
const assert = require('assert');

const scriptPath = path.join(__dirname, 'script.js');
const scriptSource = fs.readFileSync(scriptPath, 'utf8');

function makeStubElement() {
  const classes = new Set();
  return {
    innerHTML: '',
    style: {},
    textContent: '',
    value: '',
    classList: {
      add: (c) => classes.add(c),
      remove: (c) => classes.delete(c),
      toggle: (c, force) => {
        const on = force === undefined ? !classes.has(c) : !!force;
        if (on) classes.add(c); else classes.delete(c);
        return on;
      },
      contains: (c) => classes.has(c),
    },
    querySelector: () => null,
    querySelectorAll: () => [],
    addEventListener: () => {},
    setSelectionRange: () => {},
    focus: () => {},
  };
}

function runScriptWithItems(items) {
  const store = { ['remember_everything_items']: JSON.stringify(items) };
  const alerts = [];
  const confirms = [];

  const elements = new Map();
  function getElementById(id) {
    if (!elements.has(id)) {
      const el = makeStubElement();
      el.id = id;
      elements.set(id, el);
    }
    return elements.get(id);
  }

  const globals = {
    console,
    Date,
    Math,
    JSON,
    setInterval: () => {},
    navigator: {},
    window: {},
    alert: (m) => alerts.push(m),
    confirm: (m) => { confirms.push(m); return true; },
    localStorage: {
      getItem: (k) => store[k] ?? null,
      setItem: (k, v) => { store[k] = v; },
    },
    document: {
      addEventListener: () => {},
      querySelectorAll: () => [],
      activeElement: null,
      getElementById,
    },
  };
  const context = vm.createContext(globals);
  vm.runInContext(scriptSource, context, { filename: 'script.js' });
  return { context, store, alerts, confirms };
}

function baseItem(overrides = {}) {
  const now = new Date().toISOString();
  return {
    id: 1,
    item_type: 'todo',
    title: 'Test todo',
    done: 0,
    ease_factor: 2.5,
    interval_days: 0,
    repetitions: 0,
    next_review: now,
    last_reviewed: null,
    created_at: now,
    updated_at: now,
    archived: 0,
    ...overrides,
  };
}

function runScriptWithItemsAndDOM(items) {
  const store = { ['remember_everything_items']: JSON.stringify(items) };
  const timeouts = [];
  const elements = new Map();

  function makeElement(id) {
    const classes = new Set();
    const listeners = [];
    return {
      id,
      innerHTML: '',
      style: {},
      textContent: '',
      value: '',
      classList: {
        add: (c) => classes.add(c),
        remove: (c) => classes.delete(c),
        toggle: (c, force) => {
          const on = force === undefined ? !classes.has(c) : !!force;
          if (on) classes.add(c); else classes.delete(c);
          return on;
        },
        contains: (c) => classes.has(c),
      },
      querySelector: () => null,
      querySelectorAll: () => [],
      addEventListener: (type, fn) => listeners.push({ type, fn }),
      getListeners: () => listeners,
      focus: () => {},
      blur: () => {},
    };
  }

  const document = {
    addEventListener: (type, fn) => {
      const el = elements.get('document') || makeElement('document');
      elements.set('document', el);
      el.addEventListener(type, fn);
    },
    querySelectorAll: () => [],
    activeElement: null,
    getElementById: (id) => {
      if (!elements.has(id)) elements.set(id, makeElement(id));
      return elements.get(id);
    },
  };

  const alerts = [];
  const confirms = [];

  const globals = {
    console,
    Date,
    Math,
    JSON,
    setInterval: () => {},
    setTimeout: (fn, ms) => { timeouts.push(fn); return timeouts.length; },
    navigator: {},
    window: {},
    alert: (m) => alerts.push(m),
    confirm: (m) => { confirms.push(m); return true; },
    localStorage: {
      getItem: (k) => store[k] ?? null,
      setItem: (k, v) => { store[k] = v; },
    },
    document,
  };
  const context = vm.createContext(globals);
  vm.runInContext(scriptSource, context, { filename: 'script.js' });
  return { context, store, document, timeouts, alerts, confirms };
}
function test(name, fn) {
  try {
    fn();
    console.log('✓', name);
  } catch (err) {
    console.error('✗', name);
    console.error(err);
    process.exitCode = 1;
  }
}

test('finished todo is excluded from review entries', () => {
  const { context } = runScriptWithItems([
    baseItem({ id: 1, done: 1 }),
  ]);
  const entries = context.getDueReviewEntries();
  assert.strictEqual(entries.length, 0, 'expected 0 review entries for finished todo');
});

test('unfinished todo is included in review entries', () => {
  const { context } = runScriptWithItems([
    baseItem({ id: 1, done: 0 }),
  ]);
  const entries = context.getDueReviewEntries();
  assert.strictEqual(entries.length, 1, 'expected 1 review entry for unfinished todo');
});

test('finished note remains in review entries', () => {
  const { context } = runScriptWithItems([
    baseItem({ id: 1, item_type: 'note', done: 1 }),
  ]);
  const entries = context.getDueReviewEntries();
  assert.strictEqual(entries.length, 1, 'notes should still be reviewed regardless of done state');
});

test('review badge count excludes finished todos', () => {
  const { context } = runScriptWithItems([
    baseItem({ id: 1, done: 1 }),
    baseItem({ id: 2, done: 0 }),
  ]);
  assert.strictEqual(context.getDueCount(), 1, 'expected count of 1, excluding finished todo');
});

test('SM-2 resets on failed review', () => {
  const { context } = runScriptWithItems([]);
  const result = context.calculateSM2(2, 2.5, 6, 3);
  assert.strictEqual(result.repetitions, 0);
  assert.strictEqual(result.interval, 0);
});

test('SM-2 advances on successful review', () => {
  const { context } = runScriptWithItems([]);
  const result = context.calculateSM2(4, 2.5, 0, 0);
  assert.strictEqual(result.repetitions, 1);
  assert.strictEqual(result.interval, 1);
});

test('review buttons render without hovered class after submission', () => {
  const now = new Date().toISOString();
  const { context, document } = runScriptWithItemsAndDOM([
    baseItem({ id: 1, next_review: now }),
    baseItem({ id: 2, next_review: now }),
  ]);

  const activeBtn = document.getElementById('activeBtn');
  activeBtn.classList.add('review-btn');
  document.activeElement = activeBtn;

  context.submitReview(1, -1, 4);

  const reviewCards = document.getElementById('reviewCards');
  assert(!reviewCards.innerHTML.includes('hovered'), 'rendered cards should not start with hovered class');
  assert(!reviewCards.classList.contains('review-pointer-locked'), 'review-pointer-locked should not be used');
  assert.strictEqual(typeof context.attachReviewButtonHover, 'function', 'attachReviewButtonHover should be defined');
});
test('done todo shows Reopen button instead of Abandon', () => {
  const { context } = runScriptWithItems([
    baseItem({ id: 1, done: 1 }),
  ]);
  const html = context.renderTodoCard(baseItem({ id: 1, done: 1 }));
  assert(!html.includes('archiveItem('), 'abandon button should not be rendered for completed todo');
  assert(html.includes('reopenItem('), 'reopen button should be rendered for completed todo');
});

test('archived todo shows Reopen button', () => {
  const { context } = runScriptWithItems([
    baseItem({ id: 1, archived: 1 }),
  ]);
  const html = context.renderTodoCard(baseItem({ id: 1, archived: 1 }));
  assert(html.includes('reopenItem('), 'reopen button should be rendered for archived todo');
  assert(!html.includes('archiveItem('), 'abandon button should not be rendered for archived todo');
});

test('archiveItem allows abandoning completed todo', () => {
  const item = baseItem({ id: 1, done: 1 });
  const { context, store, confirms } = runScriptWithItems([item]);
  context.archiveItem(1);
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.find(i => i.id === 1).archived, 1, 'completed todo should be archived');
  assert.strictEqual(items.find(i => i.id === 1).done, 0, 'completed todo should be cleared from done');
  assert.strictEqual(confirms.length, 1, 'should confirm abandonment');
});

test('reopenItem restores archived todo to open', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1, done: 0, archived: 1 }),
  ]);
  context.reopenItem(1);
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.find(i => i.id === 1).archived, 0, 'archived todo should be reopened');
  assert.strictEqual(items.find(i => i.id === 1).done, 0, 'reopened todo should be not done');
});

test('reopenItem restores done todo to open', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1, done: 1, archived: 0 }),
  ]);
  context.reopenItem(1);
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.find(i => i.id === 1).archived, 0, 'done todo should be reopened');
  assert.strictEqual(items.find(i => i.id === 1).done, 0, 'reopened todo should be not done');
});
test('abandonFromModal archives active todo', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1, done: 0 }),
  ]);
  const editId = context.document.getElementById('editId');
  editId.value = '1';
  context.abandonFromModal();
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.find(i => i.id === 1).archived, 1, 'active todo should be archived from modal');
});

test('abandonFromModal allows archiving done todo', () => {
  const { context, store, confirms } = runScriptWithItems([
    baseItem({ id: 1, done: 1 }),
  ]);
  const editId = context.document.getElementById('editId');
  editId.value = '1';
  context.abandonFromModal();
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.find(i => i.id === 1).archived, 1, 'done todo should be archived from modal');
  assert.strictEqual(items.find(i => i.id === 1).done, 0, 'done todo should be cleared from done');
  assert.strictEqual(confirms.length, 1, 'should confirm abandonment');
});

test('reopenFromModal reopens done todo', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1, done: 1 }),
  ]);
  const editId = context.document.getElementById('editId');
  editId.value = '1';
  context.reopenFromModal();
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.find(i => i.id === 1).done, 0, 'done todo should be reopened from modal');
  assert.strictEqual(items.find(i => i.id === 1).archived, 0, 'done todo should be unarchived from modal');
});

test('exportData builds a payload with schema_version and current items', () => {
  const { context } = runScriptWithItems([
    baseItem({ id: 1 }),
    baseItem({ id: 2, item_type: 'note' }),
  ]);
  const payload = context.buildExportPayload();
  assert.strictEqual(payload.schema_version, 2, 'payload should carry a schema version');
  assert(typeof payload.exported_at === 'string' && payload.exported_at.length > 0, 'payload should have an export timestamp');
  assert.strictEqual(payload.items.length, 2, 'payload should include all current items');
});

test('exportData works when there are no items', () => {
  const { context } = runScriptWithItems([]);
  const payload = context.buildExportPayload();
  assert.deepStrictEqual(payload.items, [], 'payload items should be an empty array');
});

test('buildExportPayload includes a stacks array reflecting current local stacks', () => {
  const { context } = runScriptWithItems([baseItem({ id: 1 })]);
  context.createStack('Errands');
  const payload = context.buildExportPayload();
  assert.strictEqual(payload.stacks.length, 1, 'payload should include current stacks');
  assert.strictEqual(payload.stacks[0].name, 'Errands');
});

test('buildExportPayload includes an empty stacks array when no stacks exist', () => {
  const { context } = runScriptWithItems([]);
  const payload = context.buildExportPayload();
  assert.strictEqual(payload.stacks.length, 0, 'payload stacks should be an empty array');
});

test('importDataFromText skips a todo that duplicates an existing todo title, keeping local metadata', () => {
  const { context, store } = runScriptWithItems([
    baseItem({
      id: 1, item_type: 'todo', title: 'Buy milk', done: 1, ease_factor: 3.1,
    }),
  ]);
  const importJson = JSON.stringify({
    schema_version: 1,
    exported_at: new Date().toISOString(),
    items: [baseItem({ id: 9, item_type: 'todo', title: 'Buy milk', done: 0, ease_factor: 2.5 })],
  });
  context.importDataFromText(importJson);
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.length, 1, 'duplicate-title todo should not be added');
  assert.strictEqual(items[0].done, 1, 'existing local todo metadata should win');
  assert.strictEqual(items[0].ease_factor, 3.1, 'existing local todo ease factor should win');
});

test('importDataFromText adds a todo whose title is new', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1, item_type: 'todo', title: 'Buy milk' }),
  ]);
  const importJson = JSON.stringify({
    schema_version: 1,
    exported_at: new Date().toISOString(),
    items: [baseItem({ id: 9, item_type: 'todo', title: 'Walk the dog' })],
  });
  context.importDataFromText(importJson);
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.length, 2, 'new-title todo should be appended');
  assert(items.find(i => i.title === 'Walk the dog'), 'new todo should be present');
});

test('importDataFromText does not treat a todo and a note with the same title as duplicates', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1, item_type: 'todo', title: 'Same title' }),
  ]);
  const importJson = JSON.stringify({
    schema_version: 1,
    exported_at: new Date().toISOString(),
    items: [baseItem({ id: 9, item_type: 'note', title: 'Same title' })],
  });
  context.importDataFromText(importJson);
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.length, 2, 'same title but different type should not be deduped');
});

test('importDataFromText assigns fresh ids to imported items so they never collide', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1, item_type: 'todo', title: 'existing' }),
  ]);
  const importJson = JSON.stringify({
    schema_version: 1,
    exported_at: new Date().toISOString(),
    items: [baseItem({ id: 1, item_type: 'todo', title: 'imported' })],
  });
  context.importDataFromText(importJson);
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.length, 2, 'both items should be present');
  const ids = items.map(i => i.id);
  assert.strictEqual(new Set(ids).size, 2, 'imported item should get a new, non-colliding id');
});


test('importDataFromText skips a note whose content duplicates an existing note, keeping local metadata', () => {
  const { context, store } = runScriptWithItems([
    baseItem({
      id: 1, item_type: 'note', title: 'The capital of France is Paris',
      next_review: '2020-01-01T00:00:00.000Z', ease_factor: 3.1, repetitions: 5,
    }),
  ]);
  const importJson = JSON.stringify({
    schema_version: 1,
    exported_at: new Date().toISOString(),
    items: [
      baseItem({
        id: 9, item_type: 'note', title: 'The capital of France is Paris',
        next_review: '2099-01-01T00:00:00.000Z', ease_factor: 2.5, repetitions: 0,
      }),
    ],
  });
  context.importDataFromText(importJson);
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.length, 1, 'duplicate-content note should not be added');
  assert.strictEqual(items[0].next_review, '2020-01-01T00:00:00.000Z', 'existing local review metadata should win');
  assert.strictEqual(items[0].ease_factor, 3.1, 'existing local ease factor should win');
});

test('importDataFromText adds a note whose content is new', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1, item_type: 'note', title: 'Existing note' }),
  ]);
  const importJson = JSON.stringify({
    schema_version: 1,
    exported_at: new Date().toISOString(),
    items: [baseItem({ id: 9, item_type: 'note', title: 'A brand new note' })],
  });
  context.importDataFromText(importJson);
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.length, 2, 'new-content note should be appended');
  assert(items.find(i => i.title === 'A brand new note'), 'new note should be present');
});

test('importDataFromText rejects invalid JSON without modifying data', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1 }),
  ]);
  assert.throws(() => context.importDataFromText('not json'), /not valid JSON/);
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.length, 1, 'invalid JSON import should not modify existing data');
});

test('importDataFromText rejects a file missing an items array', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1 }),
  ]);
  assert.throws(() => context.importDataFromText(JSON.stringify({ schema_version: 1 })), /does not look like/);
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.length, 1, 'invalid shape import should not modify existing data');
});

// ===== Import + Stacks =====

test('importDataFromText restores a stack that has no local counterpart', () => {
  const { context, store } = runScriptWithItems([]);
  const importJson = JSON.stringify({
    schema_version: 2,
    exported_at: new Date().toISOString(),
    items: [
      baseItem({ id: 9, item_type: 'todo', title: 'Buy milk', stack_id: 5 }),
      baseItem({ id: 10, item_type: 'todo', title: 'Buy eggs', stack_id: 5 }),
    ],
    stacks: [{ id: 5, name: 'Errands', created_at: '2020-01-01T00:00:00.000Z', updated_at: '2020-01-01T00:00:00.000Z' }],
  });
  context.importDataFromText(importJson);

  const stacks = JSON.parse(store['remember_everything_stacks']);
  assert.strictEqual(stacks.length, 1, 'imported stack should be created locally');
  assert.strictEqual(stacks[0].name, 'Errands');
  const newStackId = stacks[0].id;
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.length, 2, 'both imported items should be added');
  assert(items.every(i => i.stack_id === newStackId), 'both imported items should link to the new local stack');
});

test('importDataFromText merges an imported stack into an existing local stack with the same name', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1, item_type: 'todo', title: 'Existing local item' }),
  ]);
  const localStackId = context.createStack('Errands');
  context.setItemStack(1, localStackId);

  const importJson = JSON.stringify({
    schema_version: 2,
    exported_at: new Date().toISOString(),
    items: [baseItem({ id: 9, item_type: 'todo', title: 'Buy milk', stack_id: 5 })],
    stacks: [{ id: 5, name: 'Errands', created_at: '2020-01-01T00:00:00.000Z', updated_at: '2020-01-01T00:00:00.000Z' }],
  });
  context.importDataFromText(importJson);

  const stacks = JSON.parse(store['remember_everything_stacks']);
  assert.strictEqual(stacks.length, 1, 'no duplicate stack should be created');
  assert.strictEqual(stacks[0].id, localStackId, 'existing local stack id should be preserved');
  assert.strictEqual(stacks[0].name, 'Errands', 'existing local stack name should be unchanged');

  const items = JSON.parse(store['remember_everything_items']);
  const imported = items.find(i => i.title === 'Buy milk');
  assert.strictEqual(imported.stack_id, localStackId, 'imported item should link to the existing local stack');
});

test('importDataFromText does not change stack_id of an existing item skipped as a duplicate', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1, item_type: 'todo', title: 'Buy milk' }),
  ]);
  const localStackId = context.createStack('Errands');
  context.setItemStack(1, localStackId);

  const importJson = JSON.stringify({
    schema_version: 2,
    exported_at: new Date().toISOString(),
    items: [baseItem({ id: 9, item_type: 'todo', title: 'Buy milk', stack_id: null })],
    stacks: [],
  });
  context.importDataFromText(importJson);

  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.length, 1, 'duplicate should not be added');
  assert.strictEqual(items[0].stack_id, localStackId, 'existing item stack membership should be untouched by the import');
});

test('importDataFromText imports a pre-stacks (v1) export file with no stacks field', () => {
  const { context, store } = runScriptWithItems([]);
  const importJson = JSON.stringify({
    schema_version: 1,
    exported_at: new Date().toISOString(),
    items: [baseItem({ id: 9, item_type: 'todo', title: 'Buy milk' })],
  });
  context.importDataFromText(importJson);

  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.length, 1, 'item should import normally');
  assert.strictEqual(store['remember_everything_stacks'], undefined, 'no stacks should be created from a file with no stacks field');
});

test('importDataFromText returns a summary with correct counts in a mixed scenario', () => {
  const { context } = runScriptWithItems([
    baseItem({ id: 1, item_type: 'todo', title: 'Existing todo' }),
  ]);
  context.createStack('Errands');

  const importJson = JSON.stringify({
    schema_version: 2,
    exported_at: new Date().toISOString(),
    items: [
      baseItem({ id: 9, item_type: 'todo', title: 'Existing todo' }), // duplicate -> skipped
      baseItem({ id: 10, item_type: 'todo', title: 'New todo', stack_id: 5 }), // new, merges into existing stack
      baseItem({ id: 11, item_type: 'note', title: 'New note', stack_id: 6 }), // new, new stack
    ],
    stacks: [
      { id: 5, name: 'Errands', created_at: '2020-01-01T00:00:00.000Z', updated_at: '2020-01-01T00:00:00.000Z' },
      { id: 6, name: 'Reading', created_at: '2020-01-01T00:00:00.000Z', updated_at: '2020-01-01T00:00:00.000Z' },
    ],
  });
  const summary = context.importDataFromText(importJson);

  assert.strictEqual(summary.itemsAdded, 2, 'two new items should be added');
  assert.strictEqual(summary.itemsSkipped, 1, 'one duplicate item should be skipped');
  assert.strictEqual(summary.stacksAdded, 1, 'one new stack should be added');
  assert.strictEqual(summary.stacksMerged, 1, 'one stack should be merged into existing');
});

// ===== Stacks =====

test('createStack requires a non-empty name', () => {
  const { context, store } = runScriptWithItems([]);
  const id = context.createStack('   ');
  assert.strictEqual(id, null, 'blank name should not create a stack');
  assert.strictEqual(store['remember_everything_stacks'], undefined, 'no stack should be persisted');
});

test('createStack persists a named stack', () => {
  const { context, store } = runScriptWithItems([]);
  const id = context.createStack('Errands');
  assert.strictEqual(typeof id, 'number');
  const stacks = JSON.parse(store['remember_everything_stacks']);
  assert.strictEqual(stacks.length, 1);
  assert.strictEqual(stacks[0].name, 'Errands');
});

test('setItemStack assigns stack_id and getStackMembers reflects it', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1 }),
    baseItem({ id: 2 }),
  ]);
  const stackId = context.createStack('Errands');
  context.setItemStack(1, stackId);
  context.setItemStack(2, stackId);
  const items = JSON.parse(store['remember_everything_items']);
  const members = context.getStackMembers(stackId, items);
  assert.strictEqual(members.length, 2, 'both items should be members of the stack');
});

test('gcStacks removes a stack once its last member leaves', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1 }),
  ]);
  const stackId = context.createStack('Errands');
  context.setItemStack(1, stackId);
  let stacks = JSON.parse(store['remember_everything_stacks']);
  assert.strictEqual(stacks.length, 1, 'stack should exist while it has a member');

  context.setItemStack(1, null);
  stacks = JSON.parse(store['remember_everything_stacks']);
  assert.strictEqual(stacks.length, 0, 'empty stack should be garbage collected');
});

test('renameStack requires a non-empty name and updates existing stacks', () => {
  const { context, store } = runScriptWithItems([]);
  const stackId = context.createStack('Errands');
  const ok = context.renameStack(stackId, 'Weekend Errands');
  assert.strictEqual(ok, true);
  let stacks = JSON.parse(store['remember_everything_stacks']);
  assert.strictEqual(stacks[0].name, 'Weekend Errands');

  const rejected = context.renameStack(stackId, '  ');
  assert.strictEqual(rejected, false, 'blank rename should be rejected');
  stacks = JSON.parse(store['remember_everything_stacks']);
  assert.strictEqual(stacks[0].name, 'Weekend Errands', 'name should be unchanged after rejected rename');
});

test('buildDisplayList collapses stacked items to a single entry at the most urgent member\'s position', () => {
  const now = new Date();
  const soon = new Date(now.getTime() + 1000).toISOString();
  const later = new Date(now.getTime() + 10 * 24 * 60 * 60 * 1000).toISOString();
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1, next_review: later }),
    baseItem({ id: 2, next_review: soon }),
    baseItem({ id: 3, next_review: later }),
  ]);
  const stackId = context.createStack('Errands');
  context.setItemStack(1, stackId);
  context.setItemStack(2, stackId); // most urgent member of the stack

  const items = JSON.parse(store['remember_everything_items']);
  const entries = context.buildDisplayList(items, items);

  assert.strictEqual(entries.length, 2, 'stack should collapse to one entry, plus the unstacked item');
  assert.strictEqual(entries[0].type, 'stack', 'stack should occupy the position of its most urgent member');
  assert.strictEqual(entries[0].members.length, 2, 'stack entry should include both members');
  assert.strictEqual(entries[1].type, 'item');
  assert.strictEqual(entries[1].item.id, 3);
});

test('buildDisplayList omits a stack with no members in the given pool', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1, item_type: 'note' }),
  ]);
  const stackId = context.createStack('Notes only');
  context.setItemStack(1, stackId);

  const items = JSON.parse(store['remember_everything_items']);
  // Simulate the Todos tab: only todo-type items are in the ranking/member pool.
  const todosOnly = items.filter(i => i.item_type === 'todo');
  const entries = context.buildDisplayList(todosOnly, todosOnly);
  assert.strictEqual(entries.length, 0, 'stack with no todo members should not appear on the Todos tab');
});


