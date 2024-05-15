-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount) 
-- sqlc.arg(amount) is a placeholder for the actual value, sqlc when generating the go
--code will name the function AddAccountBalance and the parameters required to call it
-- will be named AddAccountBalanceParams and the struct will have a field named Amount 
-- and ID instead of Balance and ID
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: ListAllAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListUserAccounts :many
SELECT * FROM accounts
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;