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
    flashcards.forEach((card, idx) => {
        const cardDiv = document.createElement('div');
        cardDiv.className = 'flashcard';
        cardDiv.setAttribute('draggable', 'true');
        cardDiv.setAttribute('data-idx', idx);
        cardDiv.innerHTML = `
            <div class="question">${card.question}</div>
            <div class="answer">${card.answer}</div>
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

flashcardsList.addEventListener('click', function(e) {
    if (e.target.classList.contains('delete-btn')) {
        const idx = e.target.getAttribute('data-idx');
        const flashcards = getFlashcards();
        flashcards.splice(idx, 1);
        saveFlashcards(flashcards);
        renderFlashcards();
    }
});

// Initial render
renderFlashcards();
