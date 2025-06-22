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
        cardDiv.innerHTML = `
            <div class="question">${card.question}</div>
            <div class="answer">${card.answer}</div>
            <button class="delete-btn" title="Delete" data-idx="${idx}">&times;</button>
        `;
        flashcardsList.appendChild(cardDiv);
    });
}

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
