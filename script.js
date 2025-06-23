// Flashcard app logic
const form = document.getElementById('flashcard-form');
const questionInput = document.getElementById('question');
const answerInput = document.getElementById('answer');
const flashcardsList = document.getElementById('flashcards-list');

// Store flashcards in localStorage for persistence
function getFlashcards() {
    return JSON.parse(localStorage.getItem('flashcards') || '[]');
}

function saveFlashcards(flashcards) {
    localStorage.setItem('flashcards', JSON.stringify(flashcards));
}

function renderFlashcards() {
    const flashcards = getFlashcards();
    flashcardsList.innerHTML = '';
    // Track blur state for each card
    if (!window.blurStates || window.blurStates.length !== flashcards.length) {
        window.blurStates = Array(flashcards.length).fill(true);
    }
    flashcards.forEach((card, idx) => {
        const cardDiv = document.createElement('div');
        cardDiv.className = 'flashcard';
        cardDiv.setAttribute('draggable', 'true');
        cardDiv.setAttribute('data-idx', idx);
        cardDiv.tabIndex = 0;
        let answerHtml = `<div class="answer${window.blurStates[idx] ? ' blurred' : ''}">${card.answer}</div>`;
        if (window.blurStates[idx]) {
            answerHtml += `<button class="reveal-btn" data-idx="${idx}">Reveal</button>`;
        }
        cardDiv.innerHTML = `
            <div class="question">${card.question}</div>
            ${answerHtml}
            <button class="delete-btn" title="Delete" data-idx="${idx}">&times;</button>
        `;
        flashcardsList.appendChild(cardDiv);
    });
}

// Drag & Drop logic
let dragSrcIdx = null;

flashcardsList.addEventListener('dragstart', function(e) {
    const card = e.target.closest('.flashcard');
    if (!card) return;
    dragSrcIdx = card.getAttribute('data-idx');
    card.classList.add('dragging');
});

flashcardsList.addEventListener('dragend', function(e) {
    const card = e.target.closest('.flashcard');
    if (card) card.classList.remove('dragging');
    dragSrcIdx = null;
});

flashcardsList.addEventListener('dragover', function(e) {
    e.preventDefault();
    const card = e.target.closest('.flashcard');
    if (!card) return;
    card.classList.add('drag-over');
});

flashcardsList.addEventListener('dragleave', function(e) {
    const card = e.target.closest('.flashcard');
    if (card) card.classList.remove('drag-over');
});

flashcardsList.addEventListener('drop', function(e) {
    e.preventDefault();
    const card = e.target.closest('.flashcard');
    if (!card || dragSrcIdx === null) return;
    card.classList.remove('drag-over');
    const dropIdx = card.getAttribute('data-idx');
    if (dragSrcIdx === dropIdx) return;
    const flashcards = getFlashcards();
    const [moved] = flashcards.splice(dragSrcIdx, 1);
    flashcards.splice(dropIdx, 0, moved);
    saveFlashcards(flashcards);
    renderFlashcards();
});

form.addEventListener('submit', function(e) {
    e.preventDefault();
    const question = questionInput.value.trim();
    const answer = answerInput.value.trim();
    if (!question || !answer) return;
    const flashcards = getFlashcards();
    flashcards.push({ question, answer });
    saveFlashcards(flashcards);
    renderFlashcards();
    form.reset();
});

// Reveal answer logic
flashcardsList.addEventListener('click', function(e) {
    if (e.target.classList.contains('delete-btn')) {
        const idx = e.target.getAttribute('data-idx');
        const flashcards = getFlashcards();
        flashcards.splice(idx, 1);
        saveFlashcards(flashcards);
        window.blurStates.splice(idx, 1);
        renderFlashcards();
        return;
    }
    if (e.target.classList.contains('reveal-btn')) {
        const idx = e.target.getAttribute('data-idx');
        window.blurStates[idx] = false;
        renderFlashcards();
        return;
    }
});

// Blur answer when clicking outside
window.addEventListener('mousedown', function(e) {
    const cards = document.querySelectorAll('.flashcard');
    let changed = false;
    cards.forEach((card, idx) => {
        if (!card.contains(e.target) && window.blurStates[idx] === false) {
            window.blurStates[idx] = true;
            changed = true;
        }
    });
    if (changed) renderFlashcards();
});

// Initial render
renderFlashcards();
