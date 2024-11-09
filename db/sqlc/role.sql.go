// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: role.sql

package db

import (
	"context"
)

const countAllRole = `-- name: CountAllRole :one
SELECT COUNT(1) AS count 
FROM "role"
WHERE deleted IS False
AND LOWER(name) LIKE '%' || LOWER($1) || '%'
LIMIT 1
`

func (q *Queries) CountAllRole(ctx context.Context, lower string) (int64, error) {
	row := q.db.QueryRow(ctx, countAllRole, lower)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createRole = `-- name: CreateRole :one
INSERT INTO "role" (name, created_at)
VALUES ($1, NOW()::TIMESTAMP)
RETURNING id, name, created_at, updated_at, deleted
`

func (q *Queries) CreateRole(ctx context.Context, name string) (Role, error) {
	row := q.db.QueryRow(ctx, createRole, name)
	var i Role
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Deleted,
	)
	return i, err
}

const deleteRole = `-- name: DeleteRole :exec
UPDATE "role"
SET updated_at = NOW()::TIMESTAMP,
deleted = TRUE
WHERE deleted IS False 
AND id = $1
`

func (q *Queries) DeleteRole(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteRole, id)
	return err
}

const getAllRole = `-- name: GetAllRole :many
SELECT 
id, name
FROM "role"
WHERE deleted IS False
AND LOWER(name) LIKE '%' || LOWER($1) || '%'
ORDER BY name ASC
OFFSET $2
LIMIT $3
`

type GetAllRoleParams struct {
	Lower  string `json:"lower"`
	Offset int32  `json:"offset"`
	Limit  int32  `json:"limit"`
}

type GetAllRoleRow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) GetAllRole(ctx context.Context, arg GetAllRoleParams) ([]GetAllRoleRow, error) {
	rows, err := q.db.Query(ctx, getAllRole, arg.Lower, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllRoleRow{}
	for rows.Next() {
		var i GetAllRoleRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoleById = `-- name: GetRoleById :one
SELECT
id, name
FROM "role"
WHERE deleted IS False 
AND id = $1
LIMIT 1
`

type GetRoleByIdRow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) GetRoleById(ctx context.Context, id int64) (GetRoleByIdRow, error) {
	row := q.db.QueryRow(ctx, getRoleById, id)
	var i GetRoleByIdRow
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getRoleByUserId = `-- name: GetRoleByUserId :many
SELECT r.id
, r.name
FROM "role" r
INNER JOIN "user_role" ur ON r.id = ur.role_id AND ur.deleted IS false
WHERE r.deleted IS false
AND ur.user_id = $1
ORDER BY r.name ASC
`

type GetRoleByUserIdRow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) GetRoleByUserId(ctx context.Context, userID int64) ([]GetRoleByUserIdRow, error) {
	rows, err := q.db.Query(ctx, getRoleByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetRoleByUserIdRow{}
	for rows.Next() {
		var i GetRoleByUserIdRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoleForDropDownList = `-- name: GetRoleForDropDownList :many
SELECT 
id, name
FROM "role"
WHERE deleted IS False
ORDER BY name ASC
`

type GetRoleForDropDownListRow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) GetRoleForDropDownList(ctx context.Context) ([]GetRoleForDropDownListRow, error) {
	rows, err := q.db.Query(ctx, getRoleForDropDownList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetRoleForDropDownListRow{}
	for rows.Next() {
		var i GetRoleForDropDownListRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateRole = `-- name: UpdateRole :exec
UPDATE "role"
SET name = $2,
updated_at = NOW()::TIMESTAMP
WHERE deleted IS False 
AND id = $1
`

type UpdateRoleParams struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) UpdateRole(ctx context.Context, arg UpdateRoleParams) error {
	_, err := q.db.Exec(ctx, updateRole, arg.ID, arg.Name)
	return err
}