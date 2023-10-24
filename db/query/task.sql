-- name: CreateTask :one
INSERT INTO tasks (
title,
illustration,
category_id,
status_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1 LIMIT 1;

-- name: GetTaskForUpdate :one
SELECT * FROM tasks
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListTasks :many
SELECT * FROM tasks
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateTask :one 
UPDATE tasks SET title = $2, illustration = $3, category_id = $4, status_id = $5  WHERE id = $1 RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM tasks WHERE id = $1;

-- name: CreateCategory :one
INSERT INTO categories (
	name
) VALUES (
    $1
) 
RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM categories
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;

-- name: CreateStatus :one
INSERT INTO statuses (
	name
) VALUES (
    $1
) 
RETURNING *;

-- name: GetStatus :one
SELECT * FROM statuses
WHERE id = $1 LIMIT 1;

-- name: ListStatuses :many
SELECT * FROM statuses
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteStatus :exec
DELETE FROM statuses WHERE id = $1;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (
firstname,
lastname,
password,
email,
phone,
token,
user_type_id,
refresh_token
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;