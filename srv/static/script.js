// ===== Local Storage =====
const STORAGE_KEY = 'remember_everything_items';
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
      searchCount = searchResults.querySelectorAll('.item-card').length;
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
  const items = loadItems().filter(i => i.item_type === 'todo')
    .sort((a, b) => a.archived - b.archived || a.done - b.done || b.priority - a.priority || new Date(b.created_at) - new Date(a.created_at));
  if (items.length === 0) {
    container.innerHTML = '<div class="empty-state">No todos yet. Add one above!</div>';
    return;
  }
  container.innerHTML = items.map(renderTodoCard).join('');
}

function renderTodoCard(item) {
  const done = item.done === 1;
  const archived = item.archived === 1;
  let reviewInfo;
  if (item.review_enabled === false) {
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

  return `<div class="item-card${done ? ' done' : ''}${archived ? ' abandoned' : ''}">
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

// ===== Notes =====
function loadNotes() {
  const container = document.getElementById('noteList');
  const items = loadItems().filter(i => i.item_type === 'note' && !i.archived)
    .sort((a, b) => new Date(b.updated_at) - new Date(a.updated_at));
  if (items.length === 0) {
    container.innerHTML = '<div class="empty-state">No notes yet. Type \'/note your text\' above!</div>';
    return;
  }
  container.innerHTML = items.map(renderNoteCard).join('');
}

function renderNoteCard(item) {
  let reviewInfo;
  if (item.review_enabled === false) {
    reviewInfo = `<span class="item-sr">Reviews off</span>`;
  } else if (hasClozes(item.title)) {
    ensureClozeData(item);
    const dueCount = item.cloze_data.filter(cd => new Date(cd.next_review) <= new Date()).length;
    reviewInfo = `<span class="item-sr clickable" onclick="event.stopPropagation();showHistory(${item.id})">${getClozeCount(item.title)} cloze${getClozeCount(item.title) > 1 ? 's' : ''}${dueCount > 0 ? ` \u00b7 ${dueCount} due` : ''}</span>`;
  } else {
    reviewInfo = `<span class="item-sr clickable" onclick="event.stopPropagation();showHistory(${item.id})">Next review: ${formatDate(item.next_review)}</span>`;
  }

  return `<div class="item-card">
    <div class="item-body">
      <div class="item-title">${renderTitleWithClozeHints(item.title)}</div>
      <div class="item-meta">${reviewInfo}</div>
    </div>
    <div class="item-actions">
      <button class="btn-icon" onclick="openEdit(${item.id})" title="Edit">\u270f\ufe0f</button>
      <button class="btn-icon" onclick="archiveItem(${item.id})" title="Abandon">\ud83c\udff3\ufe0f</button>
    </div>
  </div>`;
}

// ===== Search =====
function doSearch() {
  const q = document.getElementById('searchInput').value.trim().toLowerCase();
  const container = document.getElementById('searchResults');
  if (!q) { container.innerHTML = ''; return; }
  const items = loadItems().filter(i => !i.archived &&
    (i.title.toLowerCase().includes(q) || i.content.toLowerCase().includes(q)));
  if (items.length === 0) {
    container.innerHTML = '<div class="empty-state">No results found.</div>';
    return;
  }
  container.innerHTML = items.map(item => {
    if (item.item_type === 'todo') return renderTodoCard(item);
    return renderNoteCard(item);
  }).join('');
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
  onEditTypeChange();
  updateModalAbandonButton();
  document.getElementById('editModal').style.display = 'flex';
}

function onEditTypeChange() {
  const isTodo = document.getElementById('editType').value === 'todo';
  document.getElementById('editDoneLabel').style.display = isTodo ? '' : 'none';
  const typeLabel = isTodo ? 'Todo' : 'Note';
  document.getElementById('modalTitle').textContent = `Edit ${typeLabel}`;
  updateModalAbandonButton();
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

  // Reset SM-2 if title was edited
  if (titleChanged) {
    resetSM2(item);
  }

  saveItems(items);
  closeModal();
  refreshCurrent();
  updateTabCounts();
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
