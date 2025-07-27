// Service worker version (keep in sync with service-worker.js)
const SW_VERSION = '0.1.0.2025-07-26.23';

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
        // If answer is revealed, show prompt
        let promptHtml = '';
        if (!card.blur) {
            promptHtml = `<div class=\"remember-prompt-flex\"><span class=\"remember-prompt-text\">Did you remember it?</span><div class=\"remember-prompt-btns\"><button class=\"remember-btn yes\" data-idx=\"${idx}\" data-remembered=\"yes\">Yes</button> <button class=\"remember-btn no\" data-idx=\"${idx}\" data-remembered=\"no\">No</button></div></div>`;
        }
        cardDiv.innerHTML = `
            ${questionHtml}
            ${answerHtml}
            ${promptHtml}
        `;
        flashcardsList.appendChild(cardDiv);
    });
}

// Drag & Drop logic
let dragSrcIdx = null;
let touchDrag = {
    dragging: false,
    startIdx: null,
    currentIdx: null,
    ghost: null
};

function getCardIdxFromTouch(touch) {
    const el = document.elementFromPoint(touch.clientX, touch.clientY);
    const card = el && el.closest('.flashcard');
    return card ? parseInt(card.getAttribute('data-idx')) : null;
}

// Touch events for mobile drag & drop
flashcardsList.addEventListener('touchstart', function(e) {
    const card = e.target.closest('.flashcard');
    if (!card) return;
    touchDrag.dragging = true;
    touchDrag.startIdx = parseInt(card.getAttribute('data-idx'));
    touchDrag.currentIdx = touchDrag.startIdx;
    card.classList.add('dragging');
    // Create ghost
    touchDrag.ghost = card.cloneNode(true);
    touchDrag.ghost.style.position = 'fixed';
    touchDrag.ghost.style.pointerEvents = 'none';
    touchDrag.ghost.style.opacity = '0.7';
    touchDrag.ghost.style.zIndex = '9999';
    document.body.appendChild(touchDrag.ghost);
});

flashcardsList.addEventListener('touchmove', function(e) {
    if (!touchDrag.dragging || !touchDrag.ghost) return;
    const touch = e.touches[0];
    touchDrag.ghost.style.left = (touch.clientX - touchDrag.ghost.offsetWidth/2) + 'px';
    touchDrag.ghost.style.top = (touch.clientY - touchDrag.ghost.offsetHeight/2) + 'px';
    const overIdx = getCardIdxFromTouch(touch);
    if (overIdx !== null && overIdx !== touchDrag.currentIdx) {
        // Remove drag-over from all
        document.querySelectorAll('.flashcard.drag-over').forEach(c => c.classList.remove('drag-over'));
        // Add drag-over to new
        const overCard = flashcardsList.querySelector(`.flashcard[data-idx="${overIdx}"]`);
        if (overCard) overCard.classList.add('drag-over');
        touchDrag.currentIdx = overIdx;
    }
    e.preventDefault();
}, {passive: false});

flashcardsList.addEventListener('touchend', function(e) {
    if (!touchDrag.dragging) return;
    // Remove ghost
    if (touchDrag.ghost) {
        document.body.removeChild(touchDrag.ghost);
        touchDrag.ghost = null;
    }
    // Remove dragging/drag-over classes
    document.querySelectorAll('.flashcard.dragging').forEach(c => c.classList.remove('dragging'));
    document.querySelectorAll('.flashcard.drag-over').forEach(c => c.classList.remove('drag-over'));
    // Perform reorder if needed
    if (touchDrag.startIdx !== null && touchDrag.currentIdx !== null && touchDrag.startIdx !== touchDrag.currentIdx) {
        const flashcards = getFlashcards();
        const [moved] = flashcards.splice(touchDrag.startIdx, 1);
        flashcards.splice(touchDrag.currentIdx, 0, moved);
        saveFlashcards(flashcards);
        animateWithViewTransition(() => renderFlashcards());
    }
    touchDrag.dragging = false;
    touchDrag.startIdx = null;
    touchDrag.currentIdx = null;
});

// View Transitions API for drag & drop reorder animation
async function animateWithViewTransition(callback) {
    if (document.startViewTransition) {
        await document.startViewTransition(callback);
    } else {
        callback();
    }
}

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
    if (e.target.classList.contains('remember-btn')) {
        const idx = e.target.getAttribute('data-idx');
        const flashcards = getFlashcards();
        flashcards[idx].blur = true;
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

// Add styles for the remember prompt and buttons
const style = document.createElement('style');
style.textContent = `
.remember-prompt-flex {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: #e3f2fd;
    border-left: 4px solid #34a853;
    border-radius: 6px;
    margin-top: 12px;
    padding: 12px 10px 10px 14px;
    font-size: 1.05rem;
    color: #1765c1;
    font-weight: 500;
    border: 1px solid #b3d8fd;
    gap: 16px;
}
.remember-prompt-text {
    flex: 1 1 auto;
    text-align: left;
}
.remember-prompt-btns {
    display: flex;
    gap: 10px;
}
.remember-btn.yes {
    background: #34a853;
    color: #fff;
    border: none;
    border-radius: 6px;
    padding: 6px 18px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.2s, box-shadow 0.2s;
    box-shadow: 0 2px 8px rgba(52, 168, 83, 0.10);
}
.remember-btn.yes:hover {
    background: #2e8c46;
}
.remember-btn.no {
    background: #b71c1c;
    color: #fff;
    border: none;
    border-radius: 6px;
    padding: 6px 18px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.2s, box-shadow 0.2s;
    box-shadow: 0 2px 8px rgba(183, 28, 28, 0.10);
}
.remember-btn.no:hover {
    background: #7f1010;
}
`;
document.head.appendChild(style);

// Add version number to bottom left corner
function addVersionFooter() {
    const version = SW_VERSION;
    let footer = document.getElementById('version-footer');
    if (!footer) {
        footer = document.createElement('div');
        footer.id = 'version-footer';
        document.body.appendChild(footer);
    }
    footer.innerHTML = `<span>${version}</span>`;
    footer.style.position = 'fixed';
    footer.style.left = '18px';
    footer.style.bottom = '10px';
    footer.style.fontSize = '0.95rem';
    footer.style.color = '#1a73e8';
    footer.style.background = 'rgba(255,255,255,0.85)';
    footer.style.padding = '4px 12px';
    footer.style.borderRadius = '8px';
    footer.style.boxShadow = '0 2px 8px rgba(26,115,232,0.08)';
    footer.style.zIndex = '1000';
    footer.style.userSelect = 'all';
}
addVersionFooter();

// --- Pull-to-refresh: reload page on pull down gesture (mobile-friendly) with animation and blur ---
(function enablePullToRefresh() {
    let startY = null;
    let pulling = false;
    let threshold = 60; // px to trigger refresh
    let scrollable = () => window.scrollY === 0;
    let indicator = null;
    let blurOverlay = null;
    // Create pull-to-refresh indicator
    function showIndicator(offset, spinning) {
        if (!indicator) {
            indicator = document.createElement('div');
            indicator.id = 'pull-to-refresh-indicator';
            indicator.style.position = 'fixed';
            indicator.style.top = '0';
            indicator.style.left = '0';
            indicator.style.right = '0';
            indicator.style.height = '48px';
            indicator.style.display = 'flex';
            indicator.style.alignItems = 'center';
            indicator.style.justifyContent = 'center';
            indicator.style.zIndex = '2000';
            indicator.style.pointerEvents = 'none';
            indicator.style.transition = 'transform 0.2s cubic-bezier(.4,2,.6,1), opacity 0.2s';
            indicator.innerHTML = '<svg width="28" height="28" viewBox="0 0 28 28" style="transform:rotate(0deg);transition:transform 0.3s;"><circle cx="14" cy="14" r="12" stroke="#1a73e8" stroke-width="3" fill="none" opacity="0.18"/><path d="M7 14a7 7 0 0 1 14 0" stroke="#1a73e8" stroke-width="3" fill="none" stroke-linecap="round"/></svg>';
            document.body.appendChild(indicator);
        }
        indicator.style.opacity = '1';
        indicator.style.transform = `translateY(${offset}px)`;
        const svg = indicator.querySelector('svg');
        if (spinning) {
            svg.style.transition = 'transform 0.5s linear';
            svg.style.transform = 'rotate(720deg)';
        } else {
            svg.style.transition = 'transform 0.2s';
            svg.style.transform = `rotate(${Math.min(offset, threshold) * 2}deg)`;
        }
    }
    function showBlur() {
        if (!blurOverlay) {
            blurOverlay = document.createElement('div');
            blurOverlay.id = 'pull-to-refresh-blur';
            blurOverlay.style.position = 'fixed';
            blurOverlay.style.top = '0';
            blurOverlay.style.left = '0';
            blurOverlay.style.width = '100vw';
            blurOverlay.style.height = '100vh';
            blurOverlay.style.zIndex = '1500';
            blurOverlay.style.backdropFilter = 'blur(7px)';
            blurOverlay.style.background = 'rgba(255,255,255,0.18)';
            blurOverlay.style.pointerEvents = 'none';
            blurOverlay.style.transition = 'opacity 0.3s';
            blurOverlay.style.opacity = '1';
            document.body.appendChild(blurOverlay);
        } else {
            blurOverlay.style.opacity = '1';
        }
    }
    function hideBlur() {
        if (blurOverlay) {
            blurOverlay.style.opacity = '0';
            setTimeout(() => {
                if (blurOverlay && blurOverlay.parentNode) blurOverlay.parentNode.removeChild(blurOverlay);
                blurOverlay = null;
            }, 350);
        }
    }
    function hideIndicator() {
        if (indicator) {
            indicator.style.opacity = '0';
            indicator.style.transform = 'translateY(0)';
            setTimeout(() => {
                if (indicator && indicator.parentNode) indicator.parentNode.removeChild(indicator);
                indicator = null;
            }, 350);
        }
    }
    function isPullingOnCard(touchEvent) {
        // Check if the touchstart or touchmove target is inside a flashcard
        const touch = touchEvent.touches && touchEvent.touches[0];
        if (!touch) return false;
        const el = document.elementFromPoint(touch.clientX, touch.clientY);
        return el && el.closest('.flashcard');
    }
    document.addEventListener('touchstart', function(e) {
        if (e.touches.length !== 1) return;
        if (!scrollable()) return;
        if (isPullingOnCard(e)) return; // Don't start pull-to-refresh on cards
        startY = e.touches[0].clientY;
        pulling = false;
    }, {passive: true});
    document.addEventListener('touchmove', function(e) {
        if (startY === null) return;
        if (isPullingOnCard(e)) return; // Don't trigger pull-to-refresh if finger is on a card
        let deltaY = e.touches[0].clientY - startY;
        if (deltaY > 0 && scrollable()) {
            showIndicator(Math.min(deltaY, threshold * 1.5), false);
            showBlur();
        }
        if (deltaY > threshold && scrollable()) {
            pulling = true;
        }
    }, {passive: false});
    document.addEventListener('touchend', function(e) {
        if (pulling) {
            showIndicator(threshold, true);
            showBlur();
            setTimeout(() => {
                hideIndicator();
                // Blur stays until navigation
                // If offline, reload without cb param to use cache
                if (!navigator.onLine) {
                    window.location.reload();
                } else {
                    const url = new URL(window.location.href);
                    url.searchParams.set('cb', Date.now());
                    window.location.replace(url.toString());
                }
            }, 400);
        } else {
            hideIndicator();
            hideBlur();
        }
        startY = null;
        pulling = false;
    }, {passive: true});
})();

// --- Restart service worker on refresh with cache-busting param or reload ---
function restartServiceWorkerOnRefresh() {
    const url = new URL(window.location.href);
    const cb = url.searchParams.get('cb');
    if (!('serviceWorker' in navigator)) return;
    if (cb && navigator.onLine) { // Only refresh SW if online
        // Unregister and re-register with cache-busting param
        navigator.serviceWorker.getRegistrations().then(regs => {
            Promise.all(regs.map(reg => reg.unregister())).then(() => {
                navigator.serviceWorker.register('service-worker.js?cb=' + cb).then(reg => {
                    // After registration, send a message to force activation
                    if (reg.active) {
                        reg.active.postMessage({ type: 'RESTART_SW' });
                    } else if (reg.installing) {
                        reg.installing.addEventListener('statechange', function(e) {
                            if (e.target.state === 'activated') {
                                reg.active && reg.active.postMessage({ type: 'RESTART_SW' });
                            }
                        });
                    }
                });
            });
        });
    }
}
restartServiceWorkerOnRefresh();

// Register service worker for offline PWA support (if not already registered)
if ('serviceWorker' in navigator) {
    window.addEventListener('load', function() {
        if (!navigator.serviceWorker.controller) {
            navigator.serviceWorker.register('service-worker.js').catch(() => {});
        }
    });
}
