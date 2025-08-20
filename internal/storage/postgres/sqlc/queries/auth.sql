-- name: IsNameTaken :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE name = sqlc.arg('name')
) AS exists;

-- name: IsEmailTaken :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE email = sqlc.arg('email')
) AS exists;

-- name: CreateUser :one
INSERT INTO users (name, password, email, role)
VALUES (sqlc.arg('name'), sqlc.arg('password'), sqlc.arg('email'), sqlc.arg('role'))
RETURNING id;

-- name: GetUserByName :one
SELECT id, name, password, role
FROM users
WHERE name = sqlc.arg('name');
