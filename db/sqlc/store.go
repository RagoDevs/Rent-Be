package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Store interface {
	Querier
	NewToken(id uuid.UUID, expiry time.Time, scope string) (*TokenLoc, error)
	BulkInsert(ctx context.Context, houses []HouseBulk) error
	TxnCreateTenant(ctx context.Context, args CreateTenantParams) error
	TxnUpdateTenantHouse(ctx context.Context, args UpdateTenantParams, prev_house_id uuid.UUID) error
	TxnRemoveTenantHouse(ctx context.Context, args UpdateTenantParams) error
}

type SQLStore struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
