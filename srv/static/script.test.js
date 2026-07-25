const fs = require('fs');
const path = require('path');
const vm = require('vm');
const assert = require('assert');

const scriptPath = path.join(__dirname, 'script.js');
const scriptSource = fs.readFileSync(scriptPath, 'utf8');

function makeStubElement() {
  return {
    innerHTML: '',
    style: {},
    textContent: '',
    value: '',
    querySelector: () => null,
    addEventListener: () => {},
    setSelectionRange: () => {},
    focus: () => {},
  };
}

function runScriptWithItems(items) {
  const store = { ['remember_everything_items']: JSON.stringify(items) };
  const globals = {
    console,
    Date,
    Math,
    JSON,
    setInterval: () => {},
    navigator: {},
    window: {},
    localStorage: {
      getItem: (k) => store[k] ?? null,
      setItem: (k, v) => { store[k] = v; },
    },
    document: {
      addEventListener: () => {},
      querySelectorAll: () => [],
      activeElement: null,
      getElementById: () => makeStubElement(),
    },
  };
  const context = vm.createContext(globals);
  vm.runInContext(scriptSource, context, { filename: 'script.js' });
  return { context, store };
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
        contains: (c) => classes.has(c),
      },
      querySelector: () => null,
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

  const globals = {
    console,
    Date,
    Math,
    JSON,
    setInterval: () => {},
    setTimeout: (fn, ms) => { timeouts.push(fn); return timeouts.length; },
    navigator: {},
    window: {},
    localStorage: {
      getItem: (k) => store[k] ?? null,
      setItem: (k, v) => { store[k] = v; },
    },
    document,
  };
  const context = vm.createContext(globals);
  vm.runInContext(scriptSource, context, { filename: 'script.js' });
  return { context, store, document, timeouts };
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

test('submitReview locks pointer events to prevent hover carry-over', () => {
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
  assert(reviewCards.classList.contains('review-pointer-locked'), 'review cards should be pointer-locked after submission');
});
