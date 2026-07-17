-- name: CreateItem :one
INSERT INTO items (item_type, title, content, done, priority, due_date, tags, next_review, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'), datetime('now'))
RETURNING *;

-- name: GetItem :one
SELECT * FROM items WHERE id = ? AND archived = 0;

-- name: ListTodos :many
SELECT * FROM items WHERE item_type = 'todo' AND archived = 0 ORDER BY done ASC, priority DESC, created_at DESC;

-- name: ListNotes :many
SELECT * FROM items WHERE item_type = 'note' AND archived = 0 ORDER BY updated_at DESC;

-- name: ListDueForReview :many
SELECT * FROM items WHERE archived = 0 AND next_review <= datetime('now') ORDER BY next_review ASC;

-- name: UpdateItem :one
UPDATE items SET title = ?, content = ?, priority = ?, due_date = ?, tags = ?, updated_at = datetime('now')
WHERE id = ? AND archived = 0
RETURNING *;

-- name: ToggleTodoDone :one
UPDATE items SET done = CASE WHEN done = 0 THEN 1 ELSE 0 END, updated_at = datetime('now')
WHERE id = ? AND item_type = 'todo' AND archived = 0
RETURNING *;

-- name: ReviewItem :one
UPDATE items SET ease_factor = ?, interval_days = ?, repetitions = ?, next_review = ?, last_reviewed = datetime('now'), updated_at = datetime('now')
WHERE id = ? AND archived = 0
RETURNING *;

-- name: ArchiveItem :exec
UPDATE items SET archived = 1, updated_at = datetime('now') WHERE id = ?;

-- name: SearchItems :many
SELECT * FROM items WHERE archived = 0 AND (title LIKE '%' || ? || '%' OR content LIKE '%' || ? || '%') ORDER BY updated_at DESC;

-- name: CountDueForReview :one
SELECT COUNT(*) FROM items WHERE archived = 0 AND next_review <= datetime('now');
