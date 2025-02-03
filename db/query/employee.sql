-- name: CreateEmployee :one
INSERT INTO employees (
    first_name,
    last_name,
    email,
    password_hash,
    role
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateEmployee :one
UPDATE employees
    SET first_name = $2,
    last_name = $3 
WHERE id = $1
RETURNING *;

-- name: GetEmployee :one
SELECT * FROM employees
WHERE id = $1 LIMIT 1;

-- name: ListEmployees :many
SELECT * FROM employees;

-- name: DeleteEmployee :exec
DELETE FROM employees
WHERE id = $1;
