-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;




-- But for this user has to provide all the fields, email, full_name, hashed_password
-- What if user wants to update only email or only full_name or only hashed_password
-- For that we have to use other approaches like below
UPDATE users
SET
  hashed_password = $1,
  full_name = $2,
  email = $3
WHERE
  username = $4
RETURNING *;


-- Using conditional logic to update only the fields that are provided by the user
-- Now along with the fields, user has to provide a boolean value to indicate 
-- whether to update the field or not:

-- $1,                  $2,              $3,             $4,       $5,        $6,    $7
-- set_hashed_password, hashed_password, set_full_name, full_name, set_email, email, username
UPDATE users
SET
  hashed_password = CASE WHEN $1 = TRUE THEN $2 ELSE hashed_password END,
  full_name = CASE WHEN $3 = true THEN $4 ELSE full_name END,
  email = CASE WHEN $5 IS NOT NULL THEN $6 ELSE email END
WHERE
  username = $7
RETURNING *;


-- This can be written as:
UPDATE users
SET 
-- if the field is provided in the query then update the field with the provided value
-- else keep current value of the field from the database
  hashed_password = CASE WHEN sqlc.arg(set_hashed_password)::boolean = TRUE THEN sqlc.arg(hashed_password) ELSE hashed_password END,
  full_name = CASE WHEN @set_full_name::boolean = TRUE THEN @full_name ELSE full_name END,
  email = CASE WHEN sqlc.arg(set_email) = TRUE THEN sqlc.arg(email) ELSE email END
WHERE
  username = sqlc.arg(username)
RETURNING *;

-- Note you cannot mix the positional parameters($1, $2, $3) with
-- named parameters(@set_full_name, @full_name, @set_email
-- Also sqlc.arg(set_hashed_password) is similar to @set_hashed_password



-- Using COALESCE function of postgresql to update only the fields that 
-- are provided by the user.
-- COALESCE function returns the first non-null value from the list of arguments
-- provided to it.
-- sqlc.narg is the same as sqlc. arg , but always marks the parameter as nullable.

-- name: UpdateUser :one
UPDATE users
SET
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
  full_name = COALESCE(sqlc.narg(full_name), full_name),
  email = COALESCE(sqlc.narg(email), email)
WHERE
  username = sqlc.arg(username)
RETURNING *;