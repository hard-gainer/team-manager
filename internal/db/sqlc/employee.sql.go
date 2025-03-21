// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: employee.sql

package db

import (
	"context"
)

const createEmployee = `-- name: CreateEmployee :one
INSERT INTO employees (
    id,
    name,
    email,
    role
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, name, email, role
`

type CreateEmployeeParams struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (q *Queries) CreateEmployee(ctx context.Context, arg CreateEmployeeParams) (Employee, error) {
	row := q.db.QueryRow(ctx, createEmployee,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Role,
	)
	var i Employee
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Role,
	)
	return i, err
}

const deleteEmployee = `-- name: DeleteEmployee :exec
DELETE FROM employees
WHERE id = $1
`

func (q *Queries) DeleteEmployee(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteEmployee, id)
	return err
}

const getEmployee = `-- name: GetEmployee :one
SELECT id, name, email, role FROM employees
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetEmployee(ctx context.Context, id int32) (Employee, error) {
	row := q.db.QueryRow(ctx, getEmployee, id)
	var i Employee
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Role,
	)
	return i, err
}

const listEmployees = `-- name: ListEmployees :many
SELECT id, name, email, role FROM employees
`

func (q *Queries) ListEmployees(ctx context.Context) ([]Employee, error) {
	rows, err := q.db.Query(ctx, listEmployees)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Employee{}
	for rows.Next() {
		var i Employee
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Role,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateEmployee = `-- name: UpdateEmployee :one
UPDATE employees
SET name = $2,
    email = $3,
    role = $4
WHERE id = $1
RETURNING id, name, email, role
`

type UpdateEmployeeParams struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (q *Queries) UpdateEmployee(ctx context.Context, arg UpdateEmployeeParams) (Employee, error) {
	row := q.db.QueryRow(ctx, updateEmployee,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Role,
	)
	var i Employee
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Role,
	)
	return i, err
}

const updateEmployeeEmail = `-- name: UpdateEmployeeEmail :one
UPDATE employees
SET email = $2 
WHERE id = $1
RETURNING id, name, email, role
`

type UpdateEmployeeEmailParams struct {
	ID    int32  `json:"id"`
	Email string `json:"email"`
}

func (q *Queries) UpdateEmployeeEmail(ctx context.Context, arg UpdateEmployeeEmailParams) (Employee, error) {
	row := q.db.QueryRow(ctx, updateEmployeeEmail, arg.ID, arg.Email)
	var i Employee
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Role,
	)
	return i, err
}

const updateEmployeeName = `-- name: UpdateEmployeeName :one
UPDATE employees
SET name = $2
WHERE id = $1
RETURNING id, name, email, role
`

type UpdateEmployeeNameParams struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) UpdateEmployeeName(ctx context.Context, arg UpdateEmployeeNameParams) (Employee, error) {
	row := q.db.QueryRow(ctx, updateEmployeeName, arg.ID, arg.Name)
	var i Employee
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Role,
	)
	return i, err
}

const updateEmployeeRole = `-- name: UpdateEmployeeRole :one
UPDATE employees
SET role = $2 
WHERE id = $1
RETURNING id, name, email, role
`

type UpdateEmployeeRoleParams struct {
	ID   int32  `json:"id"`
	Role string `json:"role"`
}

func (q *Queries) UpdateEmployeeRole(ctx context.Context, arg UpdateEmployeeRoleParams) (Employee, error) {
	row := q.db.QueryRow(ctx, updateEmployeeRole, arg.ID, arg.Role)
	var i Employee
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Role,
	)
	return i, err
}
