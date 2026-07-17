-- Todos and Notes with Spaced Repetition
CREATE TABLE IF NOT EXISTS items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    item_type TEXT NOT NULL CHECK(item_type IN ('todo', 'note')),
    title TEXT NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    done INTEGER NOT NULL DEFAULT 0,
    priority INTEGER NOT NULL DEFAULT 0,
    due_date TEXT,
    tags TEXT NOT NULL DEFAULT '',
    -- Spaced repetition fields (SM-2)
    ease_factor REAL NOT NULL DEFAULT 2.5,
    interval_days REAL NOT NULL DEFAULT 0,
    repetitions INTEGER NOT NULL DEFAULT 0,
    next_review TEXT NOT NULL DEFAULT (datetime('now')),
    last_reviewed TEXT,
    -- Metadata
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    archived INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_items_type ON items(item_type);
CREATE INDEX IF NOT EXISTS idx_items_next_review ON items(next_review);
CREATE INDEX IF NOT EXISTS idx_items_done ON items(done);
CREATE INDEX IF NOT EXISTS idx_items_archived ON items(archived);

INSERT OR IGNORE INTO migrations (migration_number, migration_name)
VALUES (002, '002-todos-notes');
