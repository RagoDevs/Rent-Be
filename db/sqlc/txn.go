package db

import (
	"context"
	"fmt"

	"github.com/lib/pq"
)

type HouseBulk struct {
	Location  string
	Block     string
	Partition int
	Occupied  bool
}

func (s *SQLStore) BulkInsert(ctx context.Context, houses []HouseBulk) error {

	fail := func(err error) error {
		return fmt.Errorf("BulkInsert: %v", err)
	}

	txn, err := s.db.Begin()

	if err != nil {
		return fail(err)
	}

	defer txn.Rollback()

	stmt, err := txn.PrepareContext(ctx, pq.CopyIn("house", "location", "block", "partition", "occupied"))

	if err != nil {
		return fmt.Errorf("BulkInsert: %v", err)
	}

	for _, house := range houses {
		_, err = stmt.Exec(house.Location, house.Block, house.Partition, house.Occupied)
		if err != nil {
			return fail(fmt.Errorf("error inserting house: %v", err))
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return fail(fmt.Errorf("error executing stmt: %v", err))
	}

	err = stmt.Close()
	if err != nil {
		return fail(fmt.Errorf("error closing stmt: %v", err))
	}

	err = txn.Commit()
	if err != nil {
		return fail(fmt.Errorf("error commiting stmt: %v", err))
	}

	return nil

}
