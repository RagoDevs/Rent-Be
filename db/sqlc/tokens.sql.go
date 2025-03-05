// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: tokens.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createToken = `-- name: CreateToken :exec
INSERT INTO token (hash, id, expiry, scope) VALUES ($1, $2, $3, $4)
`

type CreateTokenParams struct {
	Hash   []byte    `json:"hash"`
	ID     uuid.UUID `json:"id"`
	Expiry time.Time `json:"expiry"`
	Scope  string    `json:"scope"`
}

func (q *Queries) CreateToken(ctx context.Context, arg CreateTokenParams) error {
	_, err := q.db.ExecContext(ctx, createToken,
		arg.Hash,
		arg.ID,
		arg.Expiry,
		arg.Scope,
	)
	return err
}

const deleteAllToken = `-- name: DeleteAllToken :exec
DELETE FROM token WHERE scope = $1 AND id = $2
`

type DeleteAllTokenParams struct {
	Scope string    `json:"scope"`
	ID    uuid.UUID `json:"id"`
}

func (q *Queries) DeleteAllToken(ctx context.Context, arg DeleteAllTokenParams) error {
	_, err := q.db.ExecContext(ctx, deleteAllToken, arg.Scope, arg.ID)
	return err
}
