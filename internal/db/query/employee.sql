-- name: CreateEmployee :one
INSERT INTO employees (
    id,
    name,
    email,
    role
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateEmployeeName :one
UPDATE employees
SET name = $2
WHERE id = $1
RETURNING *;

-- name: UpdateEmployeeEmail :one
UPDATE employees
SET email = $2 
WHERE id = $1
RETURNING *;

-- name: UpdateEmployeeRole :one
UPDATE employees
SET role = $2 
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
