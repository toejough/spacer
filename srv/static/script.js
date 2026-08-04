// ===== Local Storage =====
const STORAGE_KEY = 'remember_everything_items';
const STACKS_STORAGE_KEY = 'remember_everything_stacks';
let nextId = 1;

function loadItems() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return [];
    return JSON.parse(raw);
  } catch { return []; }
}

function saveItems(items) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(items));
}

// ===== Stacks =====
function loadStacks() {
  try {
    const raw = localStorage.getItem(STACKS_STORAGE_KEY);
    if (!raw) return [];
    return JSON.parse(raw);
  } catch { return []; }
}

function saveStacks(stacks) {
  localStorage.setItem(STACKS_STORAGE_KEY, JSON.stringify(stacks));
}

function getNextStackId(stacks) {
  if (stacks.length === 0) return 1;
  return Math.max(...stacks.map(s => s.id)) + 1;
}

// Creates a stack and returns its id. Name must be non-empty — a stack must
// never exist without a name.
function createStack(name) {
  const trimmed = (name || '').trim();
  if (!trimmed) return null;
  const stacks = loadStacks();
  const now = new Date().toISOString();
  const stack = { id: getNextStackId(stacks), name: trimmed, created_at: now, updated_at: now };
  stacks.push(stack);
  saveStacks(stacks);
  return stack.id;
}

function renameStack(id, name) {
  const trimmed = (name || '').trim();
  if (!trimmed) return false;
  const stacks = loadStacks();
  const stack = stacks.find(s => s.id === id);
  if (!stack) return false;
  stack.name = trimmed;
  stack.updated_at = new Date().toISOString();
  saveStacks(stacks);
  return true;
}

// Membership is derived from item.stack_id, not stored on the stack, so a
// deleted/unstacked item can never leave a dangling reference behind.
function getStackMembers(stackId, items) {
  return items.filter(i => i.stack_id === stackId);
}

// Removes any stack with zero members. Called whenever membership changes
// (item removed from a stack, or an item is deleted) so stacks never linger
// as empty, orphaned records.
function gcStacks(items) {
  const stacks = loadStacks();
  const withMembers = stacks.filter(s => getStackMembers(s.id, items).length > 0);
  if (withMembers.length !== stacks.length) saveStacks(withMembers);
}

// Sets an item's stack_id, running stack GC afterward (covers both the
// "leaving a stack" and "joining a stack" cases uniformly).
function setItemStack(itemId, stackId) {
  const items = loadItems();
  const item = items.find(i => i.id === itemId);
  if (!item) return;
  item.stack_id = stackId || null;
  item.updated_at = new Date().toISOString();
  saveItems(items);
  gcStacks(loadItems());
}

// ===== Relevance sort =====
// state (open < done < archived) > review urgency (soonest due first, open items only) > recency
function getReviewUrgencyRank(item) {
  // Done/archived items are excluded from review, so next_review is stale —
  // neutralize this tier so recency governs ordering within those states.
  if (item.done || item.archived) return 0;
  if (item.review_enabled === false) return Infinity;
  if (hasClozes(item.title)) {
    ensureClozeData(item);
    if (!item.cloze_data.length) return Infinity;
    return Math.min(...item.cloze_data.map(cd => new Date(cd.next_review).getTime()));
  }
  return item.next_review ? new Date(item.next_review).getTime() : Infinity;
}

// Badge text for a done/archived item: relative time since the state change,
// not review info (which no longer applies once an item is closed).
function getClosedStateInfo(item) {
  const label = item.archived ? 'Abandoned' : 'Completed';
  return `<span class="item-sr">${label} ${formatDate(item.updated_at)}</span>`;
}

function compareByRelevance(a, b) {
  const stateRank = i => (i.archived ? 2 : (i.done ? 1 : 0));
  return stateRank(a) - stateRank(b)
    || getReviewUrgencyRank(a) - getReviewUrgencyRank(b)
    || new Date(b.updated_at) - new Date(a.updated_at);
}

function getNextId(items) {
  if (items.length === 0) return 1;
  return Math.max(...items.map(i => i.id)) + 1;
}

// ===== SM-2 Algorithm =====
function calculateSM2(quality, easeFactor, interval, repetitions) {
  let newEF, newInterval, newReps;
  if (quality < 3) {
    newReps = 0;
    newInterval = 0;
    newEF = easeFactor;
  } else {
    newReps = repetitions + 1;
    if (newReps === 1) newInterval = 1;
    else if (newReps === 2) newInterval = 6;
    else newInterval = interval * easeFactor;
    newEF = easeFactor + (0.1 - (5 - quality) * (0.08 + (5 - quality) * 0.02));
    if (newEF < 1.3) newEF = 1.3;
  }
  return { easeFactor: newEF, interval: newInterval, repetitions: newReps };
}

// ===== Cloze Helpers =====
const CLOZE_RE = /\{\{(.+?)\}\}/g;

function getClozeCount(title) {
  const m = title.match(CLOZE_RE);
  return m ? m.length : 0;
}

function hasClozes(title) {
  return getClozeCount(title) > 0;
}

// Return title HTML with one cloze blanked, rest revealed
function renderClozeTitle(title, blankIndex, revealed) {
  let idx = 0;
  return escHtml(title).replace(/\{\{(.+?)\}\}/g, (match, inner) => {
    if (idx++ === blankIndex) {
      return revealed
        ? `<span class="cloze-revealed">${inner}</span>`
        : `<span class="cloze-blank">[...]</span>`;
    }
    return `<span class="cloze-context">${inner}</span>`;
  });
}

// Render title in lists — show cloze markers subtly
function renderTitleWithClozeHints(title) {
  return escHtml(title).replace(/\{\{(.+?)\}\}/g,
    '<span class="cloze-hint">$1</span>');
}

// Ensure cloze_data array exists and is sized correctly
function ensureClozeData(item) {
  const count = getClozeCount(item.title);
  if (count === 0) { delete item.cloze_data; return; }
  if (!item.cloze_data) item.cloze_data = [];
  const now = new Date().toISOString();
  while (item.cloze_data.length < count) {
    item.cloze_data.push({
      ease_factor: 2.5, interval_days: 0, repetitions: 0,
      next_review: now, last_reviewed: null, review_history: [],
    });
  }
  // Trim extras if clozes were removed
  item.cloze_data.length = count;
}

// ===== State =====
let currentTab = 'review';
const expandedStacks = new Set();
let renamingStackId = null;

// ===== Stack display grouping =====
// Ranks `rankItems` by the existing relevance comparator, then walks the
// ranked list once: the first time any member of a stack is seen, it is
// replaced with a single stack entry (so the stack occupies its most
// urgent visible member's position); the stack's remaining members are
// skipped thereafter. `memberPool` supplies the (possibly wider) set of
// items used to compute the stack's full member list for rendering.
function buildDisplayList(rankItems, memberPool) {
  const stacks = loadStacks();
  const sorted = rankItems.slice().sort(compareByRelevance);
  const seen = new Set();
  const entries = [];
  for (const item of sorted) {
    const stackId = item.stack_id || null;
    if (!stackId) { entries.push({ type: 'item', item }); continue; }
    if (seen.has(stackId)) continue;
    seen.add(stackId);
    const stack = stacks.find(s => s.id === stackId);
    if (!stack) { entries.push({ type: 'item', item }); continue; }
    const members = memberPool.filter(i => i.stack_id === stackId);
    entries.push({ type: 'stack', stack, members });
  }
  return entries;
}

function renderCardForItem(item) {
  return item.item_type === 'todo' ? renderTodoCard(item) : renderNoteCard(item);
}

function renderEntries(entries) {
  return entries.map(e => e.type === 'item' ? renderCardForItem(e.item) : renderStackTile(e.stack, e.members)).join('');
}

function renderStackTile(stack, members) {
  const expanded = expandedStacks.has(stack.id);
  const renaming = renamingStackId === stack.id;
  const nameHtml = renaming
    ? `<input type="text" class="stack-rename-input" id="stackRenameInput-${stack.id}" value="${escHtml(stack.name)}"
        onclick="event.stopPropagation()"
        onkeydown="if(event.key==='Enter'){event.preventDefault();confirmStackRename(${stack.id})}else if(event.key==='Escape'){event.preventDefault();cancelStackRename()}"
        onblur="confirmStackRename(${stack.id})">`
    : `<span class="stack-name">${escHtml(stack.name)}</span>`;
  const renameBtn = renaming ? '' : `<button type="button" class="btn-icon stack-rename-btn" onclick="event.stopPropagation();startStackRename(${stack.id})" aria-label="Rename stack">✏️</button>`;

  const header = `<div class="stack-header" onclick="toggleStackExpand(${stack.id})">
    <span class="stack-expand-icon">${expanded ? '▾' : '▸'}</span>
    ${nameHtml}
    ${renameBtn}
    <span class="stack-count">${members.length}</span>
  </div>`;

  const membersHtml = expanded
    ? `<div class="stack-members" data-stack-id="${stack.id}">${members.map(renderCardForItem).join('')}</div>`
    : '';

  return `<div class="stack-tile${expanded ? ' expanded' : ''}" data-stack-id="${stack.id}">${header}${membersHtml}</div>`;
}

function toggleStackExpand(stackId) {
  if (expandedStacks.has(stackId)) expandedStacks.delete(stackId);
  else expandedStacks.add(stackId);
  refreshCurrent();
}

function startStackRename(stackId) {
  renamingStackId = stackId;
  refreshCurrent();
  const input = document.getElementById(`stackRenameInput-${stackId}`);
  if (input && input.focus) { input.focus(); if (input.select) input.select(); }
}

function cancelStackRename() {
  renamingStackId = null;
  refreshCurrent();
}

function confirmStackRename(stackId) {
  if (renamingStackId !== stackId) return; // already resolved (e.g. Escape then blur)
  const input = document.getElementById(`stackRenameInput-${stackId}`);
  const name = input ? input.value : '';
  renamingStackId = null;
  if ((name || '').trim()) renameStack(stackId, name);
  refreshCurrent();
}

// ===== Tab switching =====
function switchTab(tab) {
  currentTab = tab;
  document.querySelectorAll('.tab').forEach(t => t.classList.toggle('active', t.dataset.tab === tab));
  const helpBtn = document.getElementById('helpBtn');
  if (helpBtn) {
    helpBtn.classList.toggle('active', tab === 'help');
    helpBtn.setAttribute('aria-pressed', tab === 'help' ? 'true' : 'false');
  }
  document.querySelectorAll('.tab-content').forEach(s => s.style.display = 'none');
  document.getElementById('tab-' + tab).style.display = 'block';
  refreshCurrent();
  // Scroll the active tab into view within the tab bar (defensive: tests stub a minimal DOM).
  const activeBtn = document.querySelector && document.querySelector('.tab.active');
  if (activeBtn && activeBtn.scrollIntoView) activeBtn.scrollIntoView({ inline: 'center', block: 'nearest' });
}

// Toggle an edge-fade hint on the tab bar only when it actually overflows.
function updateTabScrollHint() {
  const tabs = document.querySelector && document.querySelector('.tabs');
  if (!tabs) return;
  tabs.classList.toggle('scrollable', tabs.scrollWidth > tabs.clientWidth + 1);
}
if (typeof window !== 'undefined' && window.addEventListener) {
  window.addEventListener('resize', updateTabScrollHint);
  window.addEventListener('load', updateTabScrollHint);
}

// ===== Quick Add =====
function quickAdd(itemType) {
  const inputId = itemType === 'todo' ? 'todoAddInput' : 'noteAddInput';
  const input = document.getElementById(inputId);
  const text = input.value.trim();
  if (!text) return;

  const items = loadItems();
  const now = new Date().toISOString();
  items.push({
    id: getNextId(items),
    item_type: itemType,
    title: text,
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
    stack_id: null,
  });
  saveItems(items);
  input.value = '';
  refreshCurrent();
  updateTabCounts();
}

// ===== Review =====
function getDueReviewEntries() {
  const now = new Date();
  const entries = [];
  for (const item of loadItems()) {
    if (item.archived || item.review_enabled === false) continue;
    if (item.item_type === 'todo' && item.done === 1) continue;
    if (hasClozes(item.title)) {
      ensureClozeData(item);
      // Collect due clozes for this item and pick the one due earliest
      const dueClozes = [];
      item.cloze_data.forEach((cd, idx) => {
        if (new Date(cd.next_review) <= now) {
          dueClozes.push({ item, clozeIndex: idx, cloze: cd,
            sortDate: new Date(cd.next_review) });
        }
      });
      if (dueClozes.length > 0) {
        dueClozes.sort((a, b) => a.sortDate - b.sortDate);
        entries.push(dueClozes[0]); // Only show one cloze per note per day
      }
    } else {
      if (new Date(item.next_review) <= now) {
        entries.push({ item, clozeIndex: -1, cloze: null,
          sortDate: new Date(item.next_review) });
      }
    }
  }
  return entries.sort((a, b) => a.sortDate - b.sortDate);
}

function getDueCount() {
  return getDueReviewEntries().length;
}

function loadReview() {
  const container = document.getElementById('reviewCards');
  const entries = getDueReviewEntries();
  if (entries.length === 0) {
    container.innerHTML = '<div class="empty-state">\ud83c\udf89 All caught up! Nothing to review.</div>';
    return;
  }
  container.innerHTML = entries.map(renderReviewEntry).join('');
  attachReviewButtonHover(container);
}

function attachReviewButtonHover(container) {
  if (!container) return;
  container.querySelectorAll('.review-btn').forEach(btn => {
    btn.addEventListener('mouseenter', () => btn.classList.add('hovered'));
    btn.addEventListener('mouseleave', () => btn.classList.remove('hovered'));
  });
}

function renderReviewEntry(entry) {
  const { item, clozeIndex, cloze } = entry;
  const typeBadge = `<span class="review-type-badge ${item.item_type}">${item.item_type}</span>`;

  if (clozeIndex >= 0) {
    // Cloze card
    const cd = cloze;
    const srInfo = `<span class="item-sr clickable" onclick="event.stopPropagation();showClozeHistory(${item.id},${clozeIndex})">interval: ${formatInterval(cd.interval_days)} \u00b7 ease: ${cd.ease_factor.toFixed(2)}</span>`;
    const clozeLabel = `<span class="cloze-label">Cloze ${clozeIndex + 1}/${getClozeCount(item.title)}</span>`;
    const titleHtml = renderClozeTitle(item.title, clozeIndex, false);
    const revealId = `reveal-${item.id}-${clozeIndex}`;

    return `<div class="review-card" id="${revealId}">
      ${typeBadge} ${clozeLabel}
      <div class="item-title">${titleHtml}</div>
      <div class="item-meta">${srInfo}</div>
      <div class="review-buttons">
        <button class="review-btn reveal-btn" onclick="revealCloze(${item.id},${clozeIndex})">Reveal</button>
      </div>
    </div>`;
  }

  // Normal card (no clozes)
  const srInfo = `<span class="item-sr clickable" onclick="event.stopPropagation();showHistory(${item.id})">interval: ${formatInterval(item.interval_days)} \u00b7 ease: ${item.ease_factor.toFixed(2)}</span>`;
  return `<div class="review-card">
    ${typeBadge}
    <div class="item-title">${escHtml(item.title)}</div>
    <div class="item-meta">${srInfo}</div>
    <div class="review-buttons">
      <button class="review-btn r0" onclick="submitReview(${item.id},-1,0)" title="Complete blackout">0 - Forgot</button>
      <button class="review-btn r1" onclick="submitReview(${item.id},-1,1)" title="Incorrect, but remembered">1 - Hard</button>
      <button class="review-btn r2" onclick="submitReview(${item.id},-1,2)" title="Incorrect, seemed easy">2 - Struggled</button>
      <button class="review-btn r3" onclick="submitReview(${item.id},-1,3)" title="Correct with difficulty">3 - OK</button>
      <button class="review-btn r4" onclick="submitReview(${item.id},-1,4)" title="Correct with hesitation">4 - Good</button>
      <button class="review-btn r5" onclick="submitReview(${item.id},-1,5)" title="Perfect recall">5 - Easy</button>
    </div>
  </div>`;
}

function revealCloze(itemId, clozeIndex) {
  const card = document.getElementById(`reveal-${itemId}-${clozeIndex}`);
  if (!card) return;
  const items = loadItems();
  const item = items.find(i => i.id === itemId);
  if (!item) return;

  const titleHtml = renderClozeTitle(item.title, clozeIndex, true);
  const titleEl = card.querySelector('.item-title');
  titleEl.innerHTML = titleHtml;

  const btns = card.querySelector('.review-buttons');
  btns.innerHTML = `
    <button class="review-btn r0" onclick="submitReview(${itemId},${clozeIndex},0)">0 - Forgot</button>
    <button class="review-btn r1" onclick="submitReview(${itemId},${clozeIndex},1)">1 - Hard</button>
    <button class="review-btn r2" onclick="submitReview(${itemId},${clozeIndex},2)">2 - Struggled</button>
    <button class="review-btn r3" onclick="submitReview(${itemId},${clozeIndex},3)">3 - OK</button>
    <button class="review-btn r4" onclick="submitReview(${itemId},${clozeIndex},4)">4 - Good</button>
    <button class="review-btn r5" onclick="submitReview(${itemId},${clozeIndex},5)">5 - Easy</button>`;
  attachReviewButtonHover(card);
}

function submitReview(id, clozeIndex, rating) {
  const items = loadItems();
  const item = items.find(i => i.id === id);
  if (!item) return;

  if (clozeIndex >= 0) {
    // Cloze review
    ensureClozeData(item);
    const cd = item.cloze_data[clozeIndex];
    if (!cd) return;
    const before = { ease: cd.ease_factor, interval: cd.interval_days, reps: cd.repetitions };
    const result = calculateSM2(rating, cd.ease_factor, cd.interval_days, cd.repetitions);
    cd.ease_factor = result.easeFactor;
    cd.interval_days = result.interval;
    cd.repetitions = result.repetitions;
    cd.next_review = new Date(Date.now() + result.interval * 24 * 60 * 60 * 1000).toISOString();
    cd.last_reviewed = new Date().toISOString();
    if (!cd.review_history) cd.review_history = [];
    cd.review_history.push({
      date: cd.last_reviewed, rating,
      ease_before: before.ease, ease_after: cd.ease_factor,
      interval_before: before.interval, interval_after: cd.interval_days,
      reps_after: cd.repetitions,
    });
    // Defer sibling clozes that are due today to tomorrow
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    tomorrow.setHours(0, 0, 0, 0);
    const now = new Date();
    item.cloze_data.forEach((sibling, idx) => {
      if (idx !== clozeIndex && new Date(sibling.next_review) <= now) {
        sibling.next_review = tomorrow.toISOString();
      }
    });
  } else {
    // Normal review
    if (!item.review_history) item.review_history = [];
    const before = { ease: item.ease_factor, interval: item.interval_days, reps: item.repetitions };
    const result = calculateSM2(rating, item.ease_factor, item.interval_days, item.repetitions);
    item.ease_factor = result.easeFactor;
    item.interval_days = result.interval;
    item.repetitions = result.repetitions;
    item.next_review = new Date(Date.now() + result.interval * 24 * 60 * 60 * 1000).toISOString();
    item.last_reviewed = new Date().toISOString();
    item.review_history.push({
      date: item.last_reviewed, rating,
      ease_before: before.ease, ease_after: item.ease_factor,
      interval_before: before.interval, interval_after: item.interval_days,
      reps_after: item.repetitions,
    });
  }

  item.updated_at = new Date().toISOString();
  saveItems(items);
  // Blur any focused review button so focus does not jump to the next card.
  if (document.activeElement && document.activeElement.classList.contains('review-btn')) {
    document.activeElement.blur();
  }
  loadReview();
  updateTabCounts();
}

function getOpenTodoCount() {
  return loadItems().filter(i => i.item_type === 'todo' && i.done === 0 && !i.archived).length;
}

function getNoteCount() {
  return loadItems().filter(i => i.item_type === 'note' && !i.archived).length;
}

function updateTabCounts() {
  const items = loadItems();
  const reviewCount = getDueCount();
  const todoCount = getOpenTodoCount();
  const noteCount = getNoteCount();

  let searchCount = 0;
  if (currentTab === 'search') {
    const searchResults = document.getElementById('searchResults');
    if (searchResults) {
      // Count top-level result entries (loose cards or collapsed stack
      // tiles), not nested member cards inside an expanded stack.
      searchCount = searchResults.children ? searchResults.children.length : 0;
    }
  }

  const badges = [
    { id: 'tabBadgeReview', count: reviewCount },
    { id: 'tabBadgeTodos', count: todoCount },
    { id: 'tabBadgeNotes', count: noteCount },
    { id: 'tabBadgeSearch', count: searchCount },
  ];

  for (const { id, count } of badges) {
    const el = document.getElementById(id);
    if (!el) continue;
    el.textContent = count;
    el.classList.toggle('visible', count > 0);
  }
}

// ===== Todos =====
function loadTodos() {
  const container = document.getElementById('todoList');
  const items = loadItems().filter(i => i.item_type === 'todo' && !isPendingDelete(i.id));
  if (items.length === 0) {
    container.innerHTML = '<div class="empty-state">No todos yet. Add one above!</div>';
    return;
  }
  container.innerHTML = renderEntries(buildDisplayList(items, items));
}

function renderTodoCard(item) {
  const done = item.done === 1;
  const archived = item.archived === 1;
  let reviewInfo;
  if (done || archived) {
    reviewInfo = getClosedStateInfo(item);
  } else if (item.review_enabled === false) {
    reviewInfo = `<span class="item-sr">Reviews off</span>`;
  } else if (hasClozes(item.title)) {
    ensureClozeData(item);
    const dueCount = item.cloze_data.filter(cd => new Date(cd.next_review) <= new Date()).length;
    reviewInfo = `<span class="item-sr clickable" onclick="event.stopPropagation();showHistory(${item.id})">${getClozeCount(item.title)} cloze${getClozeCount(item.title) > 1 ? 's' : ''}${dueCount > 0 ? ` · ${dueCount} due` : ''}</span>`;
  } else {
    reviewInfo = `<span class="item-sr clickable" onclick="event.stopPropagation();showHistory(${item.id})">Next review: ${formatDate(item.next_review)}</span>`;
  }

  // SVG icons (feather-style)
  const completeIcon = `<svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><path d="M9 12l2 2 4-4"/></svg>`;
  const abandonIcon = `<svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>`;
  const reopenIcon = `<svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"/></svg>`;
  const editIcon = `<svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z"/></svg>`;

  // Left-side action icons: Complete first, then Abandon for open; Reopen for done/abandoned
  let leftActions = '';
  if (archived || done) {
    leftActions = `<button class="btn-icon btn-reopen" onclick="reopenItem(${item.id})" aria-label="Reopen">${reopenIcon}</button>`;
  } else {
    leftActions = `<button class="btn-icon btn-complete" onclick="toggleTodo(${item.id})" aria-label="Complete">${completeIcon}</button><button class="btn-icon btn-abandon" onclick="archiveItem(${item.id})" aria-label="Abandon">${abandonIcon}</button>`;
  }

  const card = `<div class="item-card${done ? ' done' : ''}${archived ? ' abandoned' : ''}" data-id="${item.id}">
    ${dragHandleHtml(item.id)}
    <div class="item-icon-actions">
      ${leftActions}
    </div>
    <div class="item-body">
      <div class="item-title">${renderTitleWithClozeHints(item.title)}</div>
      <div class="item-meta">${reviewInfo}</div>
    </div>
    <div class="item-actions">
      <button class="btn-icon btn-edit" onclick="openEdit(${item.id})" aria-label="Edit">${editIcon}</button>
    </div>
  </div>`;
  return wrapWithSwipeReveal(card, done || archived);
}

function toggleTodo(id) {
  const items = loadItems();
  const item = items.find(i => i.id === id);
  if (!item) return;
  item.done = item.done === 0 ? 1 : 0;
  item.updated_at = new Date().toISOString();
  saveItems(items);
  loadTodos();
}

// ===== Permanent delete (swipe/modal, with undo) =====
const PENDING_DELETE_MS = 5000;
const pendingDeletes = new Map(); // id -> timeoutId

function isPendingDelete(id) {
  return pendingDeletes.has(id);
}

// Removes the item from view immediately and arms a timer; the item stays
// untouched in localStorage until the timer fires, so undo is just "cancel
// the timer" — no snapshot/restore logic needed.
function deleteItemPending(id) {
  if (pendingDeletes.has(id)) return;
  const timeoutId = setTimeout(() => commitDelete(id), PENDING_DELETE_MS);
  pendingDeletes.set(id, timeoutId);
  refreshCurrent();
  renderUndoToasts();
}

function commitDelete(id) {
  pendingDeletes.delete(id);
  const items = loadItems().filter(i => i.id !== id);
  saveItems(items);
  gcStacks(items);
  refreshCurrent();
  renderUndoToasts();
}

function undoDelete(id) {
  const timeoutId = pendingDeletes.get(id);
  if (timeoutId) clearTimeout(timeoutId);
  pendingDeletes.delete(id);
  refreshCurrent();
  renderUndoToasts();
}

// A dedicated grip icon, separate from the swipe-to-delete hit area (the
// card body), so drag and swipe never compete for the same gesture.
function dragHandleHtml(itemId) {
  return `<span class="drag-handle" data-drag-id="${itemId}" aria-label="Drag to stack">⠿</span>`;
}

// Wraps a rendered card with a reveal panel behind it, shown in the space
// the card vacates as it slides during a swipe — a red background with a
// trash icon, the standard "sliding drawer" swipe-to-delete treatment.
// Only cards eligible for delete get the wrapper; others render unwrapped.
function wrapWithSwipeReveal(cardHtml, deletable) {
  if (!deletable) return cardHtml;
  const trashIcon = `<svg viewBox="0 0 24 24" width="22" height="22" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/><path d="M10 11v6"/><path d="M14 11v6"/><path d="M9 6V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2"/></svg>`;
  return `<div class="item-card-wrapper">
    <div class="swipe-reveal" aria-hidden="true"><div class="swipe-reveal-icon">${trashIcon}</div></div>
    ${cardHtml}
  </div>`;
}

function itemAllowsDelete(item) {
  if (!item) return false;
  if (item.item_type === 'todo') return item.done === 1 || item.archived === 1;
  return item.archived === 1;
}

// Delegated pointer-gesture handling: cards are rebuilt on every refresh, so
// listeners live on the list container, not on individual cards.
const SWIPE_THRESHOLD_PX = 96;
const SWIPE_DEADZONE_PX = 8;
const SWIPE_DIRECTION_RATIO = 1.5;

function attachSwipeHandlers(containerId) {
  const container = document.getElementById(containerId);
  if (!container) return;
  let gesture = null;

  container.addEventListener('pointerdown', (e) => {
    // The drag handle owns its own gesture — never let a handle press also
    // arm swipe-to-delete on the same pointerdown.
    if (e.target.closest && e.target.closest('.drag-handle')) return;
    const card = e.target.closest ? e.target.closest('.item-card') : null;
    if (!card || !card.dataset) return;
    const id = parseInt(card.dataset.id, 10);
    const item = loadItems().find(i => i.id === id);
    if (!itemAllowsDelete(item)) return;
    gesture = { id, card, startX: e.clientX, startY: e.clientY, dx: 0, dy: 0, engaged: null };
  });

  container.addEventListener('pointermove', (e) => {
    if (!gesture) return;
    gesture.dx = e.clientX - gesture.startX;
    gesture.dy = e.clientY - gesture.startY;
    if (gesture.engaged === null) {
      if (Math.abs(gesture.dx) < SWIPE_DEADZONE_PX && Math.abs(gesture.dy) < SWIPE_DEADZONE_PX) return;
      // Only swipe-left engages delete (reveals the action on the right),
      // the standard convention — a rightward drag is treated like a
      // vertical one and left alone (no transform, no reveal).
      gesture.engaged = gesture.dx < 0 && Math.abs(gesture.dx) > Math.abs(gesture.dy) * SWIPE_DIRECTION_RATIO;
      if (!gesture.engaged) { gesture = null; return; }
      gesture.card.classList.add('swiping');
      const wrapper = gesture.card.parentElement;
      if (wrapper && wrapper.classList.contains('item-card-wrapper')) {
        gesture.wrapper = wrapper;
        wrapper.classList.add('swipe-active');
      }
    }
    const dx = Math.min(gesture.dx, 0);
    gesture.card.style.transform = `translateX(${dx}px)`;
    gesture.card.classList.toggle('swipe-armed', -dx >= SWIPE_THRESHOLD_PX);
  });

  function snapBack() {
    gesture.card.style.transform = '';
    gesture.card.classList.remove('swiping', 'swipe-armed');
    if (gesture.wrapper) gesture.wrapper.classList.remove('swipe-active');
  }

  // A deliberate release: commit the delete if past threshold, else snap back.
  function releaseGesture() {
    if (!gesture) return;
    if (gesture.engaged && -gesture.dx >= SWIPE_THRESHOLD_PX) {
      deleteItemPending(gesture.id);
    } else if (gesture.engaged) {
      snapBack();
    }
    gesture = null;
  }

  // An interrupted gesture — e.g. the browser handing off to native scroll —
  // must never commit, regardless of how far the drag had already gone.
  function abortGesture() {
    if (!gesture) return;
    if (gesture.engaged) snapBack();
    gesture = null;
  }

  container.addEventListener('pointerup', releaseGesture);
  container.addEventListener('pointercancel', abortGesture);
  container.addEventListener('pointerleave', abortGesture);
}

// ===== Drag-and-drop (stacking) =====
// Pointer-event-based, starting only from a card's drag handle — a distinct
// hit-target from the swipe-to-delete zone (the card body) — so the two
// gestures never compete. Works uniformly for mouse and touch pointers,
// mirroring the swipe gesture's own pointer-event approach.
let dragState = null;

function attachDragHandlers(containerId) {
  const container = document.getElementById(containerId);
  if (!container) return;

  container.addEventListener('pointerdown', (e) => {
    const handle = e.target.closest ? e.target.closest('.drag-handle') : null;
    if (!handle || !handle.dataset) return;
    const id = parseInt(handle.dataset.dragId, 10);
    e.preventDefault();
    dragState = { id, startX: e.clientX, startY: e.clientY, moved: false, ghost: null };
    if (handle.setPointerCapture) { try { handle.setPointerCapture(e.pointerId); } catch {} }
  });

  container.addEventListener('pointermove', (e) => {
    if (!dragState) return;
    const dx = e.clientX - dragState.startX;
    const dy = e.clientY - dragState.startY;
    if (!dragState.moved && Math.abs(dx) < SWIPE_DEADZONE_PX && Math.abs(dy) < SWIPE_DEADZONE_PX) return;
    dragState.moved = true;
    if (!dragState.ghost && document.body && document.body.appendChild) {
      const card = document.querySelector(`.item-card[data-id="${dragState.id}"]`);
      const ghost = document.createElement('div');
      ghost.className = 'drag-ghost';
      ghost.textContent = card && card.querySelector ? (card.querySelector('.item-title')?.textContent || '') : '';
      document.body.appendChild(ghost);
      dragState.ghost = ghost;
    }
    if (dragState.ghost && dragState.ghost.style) {
      dragState.ghost.style.left = `${e.clientX}px`;
      dragState.ghost.style.top = `${e.clientY}px`;
    }
    updateDropTargetHighlight(e.clientX, e.clientY, dragState.id);
  });

  function endDrag(e) {
    if (!dragState) return;
    const { id, moved, ghost } = dragState;
    if (ghost && ghost.remove) ghost.remove();
    clearDropTargetHighlight();
    if (moved && e) resolveDrop(id, e.clientX, e.clientY);
    dragState = null;
  }

  container.addEventListener('pointerup', endDrag);
  container.addEventListener('pointercancel', () => endDrag(null));
}

// Resolves the element under the pointer to a drop target, ignoring the
// dragged card itself.
function findDropTarget(clientX, clientY, draggedId) {
  if (!document.elementFromPoint) return null;
  const el = document.elementFromPoint(clientX, clientY);
  if (!el || !el.closest) return null;
  const cardEl = el.closest('.item-card');
  if (cardEl && cardEl.dataset && parseInt(cardEl.dataset.id, 10) !== draggedId) {
    return { kind: 'card', id: parseInt(cardEl.dataset.id, 10) };
  }
  const stackEl = el.closest('[data-stack-id]');
  if (stackEl && stackEl.dataset) {
    return { kind: 'stack', id: parseInt(stackEl.dataset.stackId, 10) };
  }
  return null;
}

function updateDropTargetHighlight(clientX, clientY, draggedId) {
  clearDropTargetHighlight();
  const target = findDropTarget(clientX, clientY, draggedId);
  if (!target) return;
  const el = target.kind === 'card'
    ? document.querySelector(`.item-card[data-id="${target.id}"]`)
    : document.querySelector(`[data-stack-id="${target.id}"]`);
  if (el && el.classList) el.classList.add('drop-target-active');
}

function clearDropTargetHighlight() {
  document.querySelectorAll('.drop-target-active').forEach(el => el.classList.remove('drop-target-active'));
}

function resolveDrop(draggedId, clientX, clientY) {
  const target = findDropTarget(clientX, clientY, draggedId);
  if (!target) {
    // Dropped outside any card or stack: un-stack the dragged card, if it
    // was in one.
    const dragged = loadItems().find(i => i.id === draggedId);
    if (dragged && dragged.stack_id) {
      setItemStack(draggedId, null);
      refreshCurrent();
    }
    return;
  }
  if (target.kind === 'stack') {
    setItemStack(draggedId, target.id);
    refreshCurrent();
    return;
  }
  // target.kind === 'card'
  const items = loadItems();
  const dragged = items.find(i => i.id === draggedId);
  const targetItem = items.find(i => i.id === target.id);
  if (!dragged || !targetItem) return;

  if (targetItem.stack_id) {
    setItemStack(draggedId, targetItem.stack_id);
    refreshCurrent();
  } else if (dragged.stack_id) {
    setItemStack(target.id, dragged.stack_id);
    refreshCurrent();
  } else {
    promptCreateStackAndMerge(draggedId, target.id);
  }
}

// Prompts for a stack name to create a new stack from two unstacked cards.
// Canceling leaves both cards untouched — a stack is never created without
// a name.
let pendingMergeIds = null;

function promptCreateStackAndMerge(draggedId, targetId) {
  pendingMergeIds = { draggedId, targetId };
  const input = document.getElementById('stackNameInput');
  if (input) input.value = '';
  const modal = document.getElementById('stackNameModal');
  if (modal && modal.style) modal.style.display = 'flex';
  if (input && input.focus) input.focus();
}

function confirmCreateStack() {
  if (!pendingMergeIds) return;
  const input = document.getElementById('stackNameInput');
  const name = input ? input.value : '';
  const stackId = createStack(name);
  if (!stackId) return; // empty name: no-op, modal stays open for correction
  const { draggedId, targetId } = pendingMergeIds;
  setItemStack(draggedId, stackId);
  setItemStack(targetId, stackId);
  cancelCreateStack();
  refreshCurrent();
}

function cancelCreateStack() {
  pendingMergeIds = null;
  const modal = document.getElementById('stackNameModal');
  if (modal && modal.style) modal.style.display = 'none';
}

function renderUndoToasts() {
  const container = document.getElementById('undoToastContainer');
  if (!container) return;
  const items = loadItems();
  container.innerHTML = [...pendingDeletes.keys()].map(id => {
    const item = items.find(i => i.id === id);
    const title = item ? escHtml(item.title) : 'Item';
    return `<div class="undo-toast" data-id="${id}" role="status">
      <span class="undo-toast-text">Deleted "${title}"</span>
      <button class="undo-toast-btn" onclick="undoDelete(${id})">Undo</button>
    </div>`;
  }).join('');
}

// ===== Notes =====
function loadNotes() {
  const container = document.getElementById('noteList');
  const items = loadItems().filter(i => i.item_type === 'note' && !isPendingDelete(i.id));
  if (items.length === 0) {
    container.innerHTML = '<div class="empty-state">No notes yet. Type \'/note your text\' above!</div>';
    return;
  }
  container.innerHTML = renderEntries(buildDisplayList(items, items));
}

function renderNoteCard(item) {
  const archived = item.archived === 1;
  let reviewInfo;
  if (archived) {
    reviewInfo = getClosedStateInfo(item);
  } else if (item.review_enabled === false) {
    reviewInfo = `<span class="item-sr">Reviews off</span>`;
  } else if (hasClozes(item.title)) {
    ensureClozeData(item);
    const dueCount = item.cloze_data.filter(cd => new Date(cd.next_review) <= new Date()).length;
    reviewInfo = `<span class="item-sr clickable" onclick="event.stopPropagation();showHistory(${item.id})">${getClozeCount(item.title)} cloze${getClozeCount(item.title) > 1 ? 's' : ''}${dueCount > 0 ? ` \u00b7 ${dueCount} due` : ''}</span>`;
  } else {
    reviewInfo = `<span class="item-sr clickable" onclick="event.stopPropagation();showHistory(${item.id})">Next review: ${formatDate(item.next_review)}</span>`;
  }

  const reopenIcon = `<svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"/></svg>`;
  const abandonOrReopenAction = archived
    ? `<button class="btn-icon btn-reopen" onclick="reopenItem(${item.id})" aria-label="Reopen">${reopenIcon}</button>`
    : `<button class="btn-icon" onclick="archiveItem(${item.id})" title="Abandon">\ud83c\udff3\ufe0f</button>`;

  const card = `<div class="item-card${archived ? ' abandoned' : ''}" data-id="${item.id}">
    ${dragHandleHtml(item.id)}
    <div class="item-body">
      <div class="item-title">${renderTitleWithClozeHints(item.title)}</div>
      <div class="item-meta">${reviewInfo}</div>
    </div>
    <div class="item-actions">
      <button class="btn-icon" onclick="openEdit(${item.id})" title="Edit">\u270f\ufe0f</button>
      ${abandonOrReopenAction}
    </div>
  </div>`;
  return wrapWithSwipeReveal(card, archived);
}

// ===== Search =====
function doSearch() {
  const q = document.getElementById('searchInput').value.trim().toLowerCase();
  const container = document.getElementById('searchResults');
  if (!q) { container.innerHTML = ''; return; }
  const allActive = loadItems().filter(i => !i.archived && !isPendingDelete(i.id));
  const matching = allActive.filter(i =>
    i.title.toLowerCase().includes(q) || i.content.toLowerCase().includes(q));
  if (matching.length === 0) {
    container.innerHTML = '<div class="empty-state">No results found.</div>';
    return;
  }
  // A stack appears if any member matches; expanding it shows all of the
  // stack's active members (any type), not just the ones that matched.
  container.innerHTML = renderEntries(buildDisplayList(matching, allActive));
  updateTabCounts();
}

// ===== Archive =====
function archiveItem(id) {
  const items = loadItems();
  const item = items.find(i => i.id === id);
  if (!item) return false;
  if (!confirm('Abandon this item?')) return false;
  item.archived = 1;
  item.done = 0; // clear done flag when abandoning
  item.updated_at = new Date().toISOString();
  saveItems(items);
  refreshCurrent();
  updateTabCounts();
  return true;
}

function reopenItem(id) {
  const items = loadItems();
  const item = items.find(i => i.id === id);
  if (!item) return false;
  item.archived = 0;
  item.done = 0;
  item.updated_at = new Date().toISOString();
  // Any prior review schedule is stale after the item was closed — reset to
  // fresh defaults so it becomes due for review immediately.
  const now = item.updated_at;
  item.ease_factor = 2.5;
  item.interval_days = 0;
  item.repetitions = 0;
  item.next_review = now;
  if (item.cloze_data) {
    item.cloze_data = item.cloze_data.map(() => ({
      ease_factor: 2.5, interval_days: 0, repetitions: 0,
      next_review: now, last_reviewed: null, review_history: [],
    }));
  }
  saveItems(items);
  refreshCurrent();
  updateTabCounts();
  return true;
}

// Show/hide the Abandon button inside the edit modal depending on item type and done state
function updateModalAbandonButton() {
  const btn = document.getElementById('abandonModalBtn');
  if (!btn) return;
  const id = parseInt(document.getElementById('editId').value || '0', 10);
  if (!id) {
    btn.style.display = 'none';
    return;
  }
  const items = loadItems();
  const item = items.find(i => i.id === id);
  if (!item) {
    btn.style.display = 'none';
    return;
  }
  const editTypeEl = document.getElementById('editType');
  if (editTypeEl && editTypeEl.value === 'note') {
    // Notes: always show Abandon
    btn.style.display = '';
    btn.textContent = 'Abandon';
    btn.setAttribute('onclick', 'abandonFromModal()');
    return;
  }
  // Todos: Abandon for open, Reopen for done/abandoned
  if (item.done === 1 || item.archived === 1) {
    btn.style.display = '';
    btn.textContent = 'Reopen';
    btn.setAttribute('onclick', 'reopenFromModal()');
  } else {
    btn.style.display = '';
    btn.textContent = 'Abandon';
    btn.setAttribute('onclick', 'abandonFromModal()');
  }
}

// Show/hide the Delete Forever button inside the edit modal — only visible
// when the item is Done/Abandoned (todos) or Abandoned (notes), giving
// keyboard/screen-reader users the same access the swipe gesture gives touch.
function updateModalDeleteButton() {
  const btn = document.getElementById('deleteForeverModalBtn');
  if (!btn) return;
  const id = parseInt(document.getElementById('editId').value || '0', 10);
  const item = id ? loadItems().find(i => i.id === id) : null;
  btn.style.display = itemAllowsDelete(item) ? '' : 'none';
}

function deleteForeverFromModal() {
  const id = parseInt(document.getElementById('editId').value, 10);
  if (!id) { alert('Item not found'); return false; }
  deleteItemPending(id);
  closeModal();
  return true;
}

function abandonFromModal() {
  const id = parseInt(document.getElementById('editId').value, 10);
  if (!id) { alert('Item not found'); return false; }
  const success = archiveItem(id);
  if (success) {
    closeModal();
    return true;
  }
  return false;
}

function reopenFromModal() {
  const id = parseInt(document.getElementById('editId').value, 10);
  if (!id) { alert('Item not found'); return false; }
  const success = reopenItem(id);
  if (success) {
    closeModal();
    return true;
  }
  return false;
}

// ===== Edit Modal =====
function openEdit(id) {
  const items = loadItems();
  const item = items.find(i => i.id === id);
  if (!item) { alert('Item not found'); return; }

  document.getElementById('editId').value = item.id;
  document.getElementById('editType').value = item.item_type;
  document.getElementById('editDone').value = item.done;
  document.getElementById('editTitle').value = item.title;
  document.getElementById('editReviewEnabled').value = item.review_enabled === false ? '0' : '1';
  document.getElementById('modalTitle').textContent = `Edit ${item.item_type === 'todo' ? 'Todo' : 'Note'}`;
  populateStackSelect(item.stack_id || null);
  onEditTypeChange();
  updateModalAbandonButton();
  updateModalDeleteButton();
  document.getElementById('editModal').style.display = 'flex';
}

// Fills the edit modal's stack <select> with "No stack" plus every existing
// stack, selecting the item's current one (if any).
function populateStackSelect(currentStackId) {
  const select = document.getElementById('editStack');
  if (!select) return;
  const stacks = loadStacks();
  const options = ['<option value="">No stack</option>']
    .concat(stacks.map(s => `<option value="${s.id}"${s.id === currentStackId ? ' selected' : ''}>${escHtml(s.name)}</option>`));
  select.innerHTML = options.join('');
  select.value = currentStackId ? String(currentStackId) : '';
}

function onEditTypeChange() {
  const isTodo = document.getElementById('editType').value === 'todo';
  document.getElementById('editDoneLabel').style.display = isTodo ? '' : 'none';
  const typeLabel = isTodo ? 'Todo' : 'Note';
  document.getElementById('modalTitle').textContent = `Edit ${typeLabel}`;
  updateModalAbandonButton();
  updateModalDeleteButton();
}

function closeModal() {
  document.getElementById('editModal').style.display = 'none';
}

function saveEdit() {
  const id = parseInt(document.getElementById('editId').value, 10);
  const title = document.getElementById('editTitle').value.trim();
  if (!title) { alert('Title is required'); return; }

  const items = loadItems();
  const item = items.find(i => i.id === id);
  if (!item) { alert('Item not found'); return; }

  const titleChanged = item.title !== title;
  item.item_type = document.getElementById('editType').value;
  item.done = item.item_type === 'todo' ? parseInt(document.getElementById('editDone').value, 10) : 0;
  item.title = title;
  item.review_enabled = document.getElementById('editReviewEnabled').value === '1';
  item.updated_at = new Date().toISOString();

  const stackSelect = document.getElementById('editStack');
  const newStackId = stackSelect && stackSelect.value ? parseInt(stackSelect.value, 10) : null;
  item.stack_id = newStackId;

  // Reset SM-2 if title was edited
  if (titleChanged) {
    resetSM2(item);
  }

  saveItems(items);
  gcStacks(items);
  closeModal();
  refreshCurrent();
  updateTabCounts();
}

// Renames the item's current stack from within the edit modal.
function renameStackFromModal() {
  const id = parseInt(document.getElementById('editId').value, 10);
  const item = id ? loadItems().find(i => i.id === id) : null;
  if (!item || !item.stack_id) { alert('This item is not in a stack'); return; }
  const stacks = loadStacks();
  const stack = stacks.find(s => s.id === item.stack_id);
  const name = prompt('Rename stack', stack ? stack.name : '');
  if (name === null) return; // cancelled
  if (renameStack(item.stack_id, name)) {
    populateStackSelect(item.stack_id);
  }
}

function resetSM2(item) {
  item.ease_factor = 2.5;
  item.interval_days = 0;
  item.repetitions = 0;
  item.next_review = new Date().toISOString();
  item.last_reviewed = null;
  if (item.review_history) {
    item.review_history.push({
      date: new Date().toISOString(),
      rating: -1,
      ease_before: item.ease_factor,
      ease_after: 2.5,
      interval_before: item.interval_days,
      interval_after: 0,
      reps_after: 0,
      note: 'Reset',
    });
  }
  // Also reset all cloze data
  delete item.cloze_data;
  ensureClozeData(item);
}

function resetReviews() {
  const id = parseInt(document.getElementById('editId').value, 10);
  const items = loadItems();
  const item = items.find(i => i.id === id);
  if (!item) return;
  resetSM2(item);
  saveItems(items);
  closeModal();
  refreshCurrent();
  updateTabCounts();
}

// ===== History Popup =====
function showClozeHistory(id, clozeIndex) {
  const items = loadItems();
  const item = items.find(i => i.id === id);
  if (!item) return;
  ensureClozeData(item);
  const cd = item.cloze_data[clozeIndex];
  if (!cd) return;

  // Extract the cloze text
  const matches = item.title.match(CLOZE_RE);
  const clozeText = matches && matches[clozeIndex]
    ? matches[clozeIndex].replace(/^\{\{|\}\}$/g, '') : '?';

  const h = cd.review_history || [];
  const ratingLabels = ['Forgot','Hard','Struggled','OK','Good','Easy'];

  let html = `<div class="history-popup-overlay" onclick="if(event.target===this)this.remove()">
    <div class="history-popup">
      <h3>Cloze ${clozeIndex + 1}: \u201c${escHtml(clozeText)}\u201d</h3>
      <div class="history-stats">
        <div><strong>Ease factor:</strong> ${cd.ease_factor.toFixed(2)}</div>
        <div><strong>Interval:</strong> ${formatInterval(cd.interval_days)}</div>
        <div><strong>Repetitions:</strong> ${cd.repetitions}</div>
        <div><strong>Next review:</strong> ${new Date(cd.next_review).toLocaleString()}</div>
      </div>`;

  if (h.length > 0) {
    html += `<h3>Review History</h3><div class="history-list">`;
    for (let i = h.length - 1; i >= 0; i--) {
      const r = h[i];
      const date = new Date(r.date).toLocaleDateString();
      const label = ratingLabels[r.rating] || r.rating;
      html += `<div class="history-entry">
        <span class="history-date">${date}</span>
        <span class="history-rating r${r.rating}">${r.rating} - ${label}</span>
        <span class="history-detail">ease ${r.ease_before.toFixed(2)} \u2192 ${r.ease_after.toFixed(2)} \u00b7 interval ${formatInterval(r.interval_before)} \u2192 ${formatInterval(r.interval_after)}</span>
      </div>`;
    }
    html += `</div>`;
  } else {
    html += `<p class="history-empty">No reviews yet for this cloze.</p>`;
  }

  html += `<button class="btn-primary" onclick="this.closest('.history-popup-overlay').remove()" style="margin-top:16px;width:100%">Close</button>
    </div></div>`;

  document.body.insertAdjacentHTML('beforeend', html);
}

function showHistory(id) {
  const items = loadItems();
  const item = items.find(i => i.id === id);
  if (!item) return;

  // If item has clozes, show cloze summary instead
  if (hasClozes(item.title)) {
    ensureClozeData(item);
    const matches = item.title.match(CLOZE_RE);
    let html = `<div class="history-popup-overlay" onclick="if(event.target===this)this.remove()">
      <div class="history-popup">
        <h3>Cloze Cards</h3>
        <p class="history-explainer">This item has ${matches.length} cloze deletion${matches.length > 1 ? 's' : ''}. Each one is reviewed independently. Tap a cloze to see its full history.</p>
        <div class="cloze-summary-list">`;
    matches.forEach((m, idx) => {
      const text = m.replace(/^\{\{|\}\}$/g, '');
      const cd = item.cloze_data[idx];
      const isDue = new Date(cd.next_review) <= new Date();
      html += `<div class="cloze-summary-item clickable" onclick="this.closest('.history-popup-overlay').remove();showClozeHistory(${item.id},${idx})">
        <span class="cloze-summary-num">${idx + 1}</span>
        <span class="cloze-summary-text">\u201c${escHtml(text)}\u201d</span>
        <span class="cloze-summary-stats">${formatInterval(cd.interval_days)} \u00b7 ${cd.ease_factor.toFixed(2)}${isDue ? ' \u00b7 <strong>due</strong>' : ''}</span>
      </div>`;
    });
    html += `</div>
        <h3>How SM-2 Works</h3>
        <p class="history-explainer">Each review adjusts two things: the <strong>ease factor</strong> (how quickly intervals grow) and the <strong>interval</strong> (days until next review). Rating 3+ means you remembered \u2014 intervals grow. Below 3 resets to the beginning. Higher ratings increase the ease factor, making future intervals grow faster.</p>
        <button class="btn-primary" onclick="this.closest('.history-popup-overlay').remove()" style="margin-top:16px;width:100%">Close</button>
      </div></div>`;
    document.body.insertAdjacentHTML('beforeend', html);
    return;
  }

  // Normal item history
  const h = item.review_history || [];
  const ratingLabels = ['Forgot','Hard','Struggled','OK','Good','Easy'];

  let html = `<div class="history-popup-overlay" onclick="if(event.target===this)this.remove()">
    <div class="history-popup">
      <h3>Review Stats</h3>
      <div class="history-stats">
        <div><strong>Ease factor:</strong> ${item.ease_factor.toFixed(2)}</div>
        <div><strong>Interval:</strong> ${formatInterval(item.interval_days)}</div>
        <div><strong>Repetitions:</strong> ${item.repetitions}</div>
        <div><strong>Next review:</strong> ${new Date(item.next_review).toLocaleString()}</div>
      </div>
      <h3>How SM-2 Works</h3>
      <p class="history-explainer">Each review adjusts two things: the <strong>ease factor</strong> (how quickly intervals grow) and the <strong>interval</strong> (days until next review). Rating 3+ means you remembered — intervals grow. Below 3 resets to the beginning. Higher ratings increase the ease factor, making future intervals grow faster.</p>`;

  if (h.length > 0) {
    html += `<h3>Review History</h3><div class="history-list">`;
    for (let i = h.length - 1; i >= 0; i--) {
      const r = h[i];
      const date = new Date(r.date).toLocaleDateString();
      const label = ratingLabels[r.rating] || r.rating;
      html += `<div class="history-entry">
        <span class="history-date">${date}</span>
        <span class="history-rating r${r.rating}">${r.rating} - ${label}</span>
        <span class="history-detail">ease ${r.ease_before.toFixed(2)} \u2192 ${r.ease_after.toFixed(2)} \u00b7 interval ${formatInterval(r.interval_before)} \u2192 ${formatInterval(r.interval_after)}</span>
      </div>`;
    }
    html += `</div>`;
  } else {
    html += `<p class="history-empty">No reviews yet.</p>`;
  }

  html += `<button class="btn-primary" onclick="this.closest('.history-popup-overlay').remove()" style="margin-top:16px;width:100%">Close</button>
    </div></div>`;

  document.body.insertAdjacentHTML('beforeend', html);
}

// ===== Helpers =====
function refreshCurrent() {
  if (currentTab === 'review') loadReview();
  else if (currentTab === 'todos') loadTodos();
  else if (currentTab === 'notes') loadNotes();
  updateTabCounts();
}

function escHtml(s) {
  if (!s) return '';
  return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
}

function formatDate(s) {
  if (!s) return '';
  try {
    const d = new Date(s);
    if (isNaN(d)) return s;
    const now = new Date();
    const diffMs = d - now;
    const diffDays = Math.round(diffMs / 86400000);
    if (diffDays === 0) return 'today';
    if (diffDays === 1) return 'tomorrow';
    if (diffDays === -1) return 'yesterday';
    if (diffDays > 0 && diffDays <= 7) return `in ${diffDays}d`;
    if (diffDays < 0 && diffDays >= -7) return `${-diffDays}d ago`;
    return d.toLocaleDateString();
  } catch { return s; }
}

function formatInterval(days) {
  if (days < 1) return 'now';
  if (days === 1) return '1 day';
  if (days < 30) return `${Math.round(days)} days`;
  if (days < 365) return `${Math.round(days / 30)} months`;
  return `${(days / 365).toFixed(1)} years`;
}

// ===== Cloze Editor Helpers =====
function toggleClozeSelection(textareaId) {
  const ta = document.getElementById(textareaId);
  if (!ta) return;
  const start = ta.selectionStart;
  const end = ta.selectionEnd;
  const val = ta.value;

  if (start === end) return; // nothing selected

  const selection = val.substring(start, end);
  let replacement;
  // Unwrap if already wrapped
  if (selection.startsWith('{{') && selection.endsWith('}}')) {
    replacement = selection.slice(2, -2);
  } else {
    replacement = '{{' + selection + '}}';
  }

  ta.setRangeText(replacement, start, end, 'end');
  // Restore selection around the replacement
  const newStart = start;
  const newEnd = start + replacement.length;
  ta.setSelectionRange(newStart, newEnd);
  ta.focus();
  updateClozeButtonLabel();
}

// Keyboard shortcut: Ctrl/Cmd+Shift+C in any textarea
function handleClozeKeydown(e) {
  if (e.key === 'C' && e.shiftKey && (e.ctrlKey || e.metaKey)) {
    const ta = e.target;
    if (ta.tagName === 'TEXTAREA' && ta.selectionStart !== ta.selectionEnd) {
      e.preventDefault();
      toggleClozeSelection(ta.id);
    }
  }
}
document.addEventListener('keydown', handleClozeKeydown);

// ===== Backup & Restore (export/import) =====
const EXPORT_SCHEMA_VERSION = 1;

// Build the exportable payload from current local data.
function buildExportPayload() {
  return {
    schema_version: EXPORT_SCHEMA_VERSION,
    exported_at: new Date().toISOString(),
    items: loadItems(),
  };
}

// Trigger a browser download of all local data as a JSON file.
function exportData() {
  const payload = buildExportPayload();
  const json = JSON.stringify(payload, null, 2);
  const blob = new Blob([json], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  const dateStr = new Date().toISOString().slice(0, 10);
  a.href = url;
  a.download = `remember-everything-export-${dateStr}.json`;
  document.body.appendChild(a);
  a.click();
  a.remove();
  URL.revokeObjectURL(url);
}

// Called by the Import button to open the file picker.
function triggerImport() {
  const input = document.getElementById('importFileInput');
  if (input) input.click();
}

// Handles the hidden file input's change event.
function handleImportFileChange(event) {
  const file = event.target.files && event.target.files[0];
  event.target.value = '';
  if (!file) return;
  const reader = new FileReader();
  reader.onload = () => {
    try {
      importDataFromText(reader.result);
    } catch (e) {
      alert('Import failed: ' + (e && e.message ? e.message : 'invalid file'));
    }
  };
  reader.onerror = () => alert('Import failed: could not read file');
  reader.readAsText(file);
}

// Validate & normalize an export payload into a current-shape items array.
// Throws on invalid input. Future schema_version migrations go here.
function upgradeExportData(data) {
  if (!data || typeof data !== 'object' || !Array.isArray(data.items)) {
    throw new Error('file does not look like a Remember Everything export');
  }
  return data.items;
}

// Parses exported JSON text and applies it. Exposed separately for testing.
function importDataFromText(text) {
  let data;
  try {
    data = JSON.parse(text);
  } catch {
    throw new Error('file is not valid JSON');
  }
  const items = upgradeExportData(data);
  applyImportedItems(items);
  return items;
}

// An item's "content" is its type plus title (including cloze markup).
// Everything else (done/archived, priority, tags, due date, spaced-repetition
// fields, timestamps) is metadata, not content. Trim + exact match, case-sensitive.
function itemContentKey(item) {
  return `${item.item_type}::${(item.title || '').trim()}`;
}

// Always append imported items that don't already exist locally.
// An imported item (todo or note) is skipped if its type+title matches an
// existing local item's; the existing local item (and all its metadata)
// wins and is left untouched.
function applyImportedItems(items) {
  const existing = loadItems();
  const existingContentKeys = new Set(existing.map(itemContentKey));
  let nextIdCounter = getNextId(existing);
  const toAppend = [];
  for (const item of items) {
    if (existingContentKeys.has(itemContentKey(item))) {
      continue; // duplicate content: existing local item wins
    }
    toAppend.push({ ...item, id: nextIdCounter++ });
  }
  saveItems(existing.concat(toAppend));
  refreshCurrent();
  updateTabCounts();
}


function updateClozeButtonLabel() {
  const btn = document.getElementById('clozeToggleBtn');
  const ta = document.getElementById('editTitle');
  const modal = document.getElementById('editModal');
  if (!btn || !ta || !modal || modal.style.display === 'none') return;

  const start = ta.selectionStart;
  const end = ta.selectionEnd;
  if (start === end || document.activeElement !== ta) {
    btn.textContent = 'Make cloze';
    btn.style.opacity = '.5';
    return;
  }
  btn.style.opacity = '';
  const selection = ta.value.substring(start, end);
  if (selection.startsWith('{{') && selection.endsWith('}}')) {
    btn.textContent = 'Remove cloze';
  } else {
    btn.textContent = 'Make cloze';
  }
}
document.addEventListener('selectionchange', updateClozeButtonLabel);

// ===== Init =====
loadReview();
updateTabCounts();
if (document.documentElement && document.documentElement.style.setProperty) {
  document.documentElement.style.setProperty('--swipe-threshold-px', SWIPE_THRESHOLD_PX + 'px');
}
attachSwipeHandlers('todoList');
attachSwipeHandlers('noteList');
attachDragHandlers('todoList');
attachDragHandlers('noteList');
attachDragHandlers('searchResults');
