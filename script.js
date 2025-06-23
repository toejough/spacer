// Flashcard app logic
const form = document.getElementById('flashcard-form');
const questionInput = document.getElementById('question');
const answerInput = document.getElementById('answer');
const flashcardsList = document.getElementById('flashcards-list');

// Store flashcards in localStorage for persistence
function getFlashcards() {
    // Ensure each card has a blur property (default true)
    let cards = JSON.parse(localStorage.getItem('flashcards') || '[]');
    let changed = false;
    cards.forEach(card => {
        if (typeof card.blur === 'undefined') {
            card.blur = true;
            changed = true;
        }
    });
    if (changed) saveFlashcards(cards);
    return cards;
}

function saveFlashcards(flashcards) {
    localStorage.setItem('flashcards', JSON.stringify(flashcards));
}

function hashString(str) {
    let hash = 0, i, chr;
    if (str.length === 0) return hash;
    for (i = 0; i < str.length; i++) {
        chr = str.charCodeAt(i);
        hash = ((hash << 5) - hash) + chr;
        hash |= 0; // Convert to 32bit integer
    }
    return Math.abs(hash).toString(36);
}

function renderFlashcards() {
    const flashcards = getFlashcards();
    flashcardsList.innerHTML = '';
    flashcards.forEach((card, idx) => {
        const cardDiv = document.createElement('div');
        cardDiv.className = 'flashcard';
        cardDiv.setAttribute('draggable', 'true');
        cardDiv.setAttribute('data-idx', idx);
        cardDiv.tabIndex = 0;
        const hash = hashString(card.question + '||' + card.answer);
        cardDiv.style.viewTransitionName = `flashcard-${hash}`;
        // Question row with delete button right-aligned
        let questionHtml = `<div class=\"question-row\"><div class=\"question\">${card.question}</div><button class=\"delete-btn\" title=\"Delete\" data-idx=\"${idx}\">&times;</button></div>`;
        // Answer row with reveal button (if needed)
        let answerHtml = `<div class=\"answer-row\">`;
        answerHtml += `<div class=\"answer${card.blur ? ' blurred' : ''}\">${card.answer}</div>`;
        if (card.blur) {
            answerHtml += `<button class=\"reveal-btn\" data-idx=\"${idx}\">Reveal</button>`;
        }
        answerHtml += `</div>`;
        cardDiv.innerHTML = `
            ${questionHtml}
            ${answerHtml}
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

// View Transitions API for drag & drop reorder animation
async function animateWithViewTransition(callback) {
    if (document.startViewTransition) {
        await document.startViewTransition(callback);
    } else {
        callback();
    }
}

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
    animateWithViewTransition(() => renderFlashcards());
});

form.addEventListener('submit', function(e) {
    e.preventDefault();
    const question = questionInput.value.trim();
    const answer = answerInput.value.trim();
    if (!question || !answer) return;
    const flashcards = getFlashcards();
    flashcards.push({ question, answer, blur: true });
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
        renderFlashcards();
        return;
    }
    if (e.target.classList.contains('reveal-btn')) {
        const idx = e.target.getAttribute('data-idx');
        const flashcards = getFlashcards();
        flashcards[idx].blur = false;
        saveFlashcards(flashcards);
        renderFlashcards();
        return;
    }
});

// Blur answer when clicking outside
window.addEventListener('mousedown', function(e) {
    const cards = document.querySelectorAll('.flashcard');
    let changed = false;
    const flashcards = getFlashcards();
    cards.forEach((card, idx) => {
        if (!card.contains(e.target) && flashcards[idx] && flashcards[idx].blur === false) {
            flashcards[idx].blur = true;
            changed = true;
        }
    });
    if (changed) {
        saveFlashcards(flashcards);
        renderFlashcards();
    }
});

// Initial render
renderFlashcards();
