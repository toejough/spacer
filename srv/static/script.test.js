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
    content: '',
    done: 0,
    priority: 0,
    due_date: null,
    tags: '',
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
  assert.strictEqual(payload.schema_version, 1, 'payload should carry a schema version');
  assert(typeof payload.exported_at === 'string' && payload.exported_at.length > 0, 'payload should have an export timestamp');
  assert.strictEqual(payload.items.length, 2, 'payload should include all current items');
});

test('exportData works when there are no items', () => {
  const { context } = runScriptWithItems([]);
  const payload = context.buildExportPayload();
  assert.deepStrictEqual(payload.items, [], 'payload items should be an empty array');
});

test('importDataFromText replace mode discards existing data', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1, title: 'old' }),
  ]);
  const importJson = JSON.stringify({
    schema_version: 1,
    exported_at: new Date().toISOString(),
    items: [baseItem({ id: 5, title: 'imported' })],
  });
  context.importDataFromText(importJson, 'replace');
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.length, 1, 'replace should leave only imported items');
  assert.strictEqual(items[0].id, 5, 'replace should keep imported item');
});

test('importDataFromText replace mode is skipped if user cancels confirmation', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1, title: 'old' }),
  ]);
  context.confirm = () => false;
  const importJson = JSON.stringify({
    schema_version: 1,
    exported_at: new Date().toISOString(),
    items: [baseItem({ id: 5, title: 'imported' })],
  });
  context.importDataFromText(importJson, 'replace');
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.length, 1, 'declined replace should leave data untouched');
  assert.strictEqual(items[0].id, 1, 'declined replace should keep original item');
});

test('importDataFromText merge mode adds new items and keeps existing ones on id collision', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1, title: 'existing' }),
  ]);
  const importJson = JSON.stringify({
    schema_version: 1,
    exported_at: new Date().toISOString(),
    items: [
      baseItem({ id: 1, title: 'colliding-imported' }),
      baseItem({ id: 2, title: 'new-imported' }),
    ],
  });
  context.importDataFromText(importJson, 'merge');
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.length, 2, 'merge should add only the non-colliding item');
  assert.strictEqual(items.find(i => i.id === 1).title, 'existing', 'existing item should win on id collision');
  assert(items.find(i => i.id === 2), 'new item should be added by merge');
});

test('importDataFromText rejects invalid JSON without modifying data', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1 }),
  ]);
  assert.throws(() => context.importDataFromText('not json', 'merge'), /not valid JSON/);
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.length, 1, 'invalid JSON import should not modify existing data');
});

test('importDataFromText rejects a file missing an items array', () => {
  const { context, store } = runScriptWithItems([
    baseItem({ id: 1 }),
  ]);
  assert.throws(() => context.importDataFromText(JSON.stringify({ schema_version: 1 }), 'merge'), /does not look like/);
  const items = JSON.parse(store['remember_everything_items']);
  assert.strictEqual(items.length, 1, 'invalid shape import should not modify existing data');
});

