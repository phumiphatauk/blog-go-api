-- Queries to fix import errors in querier.go file
-- Please not use this queries for other purpose
-- name: FixErrorImportPGType :one
SELECT NOW()::TIMESTAMP == $1::TIMESTAMP;

-- Queries to fix import errors in querier.go file
-- Please not use this queries for other purpose
-- name: FixErrorImportTime :one
SELECT id FROM users where created_at == $1;
