package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

func (store *SQLStore) TxnCreateTenant(ctx context.Context, args CreateTenantParams) error {

	tx, err := store.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := New(tx)

	tenant_id, err := qtx.CreateTenant(ctx, args)

	if err != nil {
		return err
	}

	house, err := qtx.GetHouseById(ctx, args.HouseID)

	if err != nil {
		return err
	}

	if house.Occupied {
		return fmt.Errorf("house is already occupied")
	}

	occupiedBy := uuid.NullUUID{
		UUID:  tenant_id,
		Valid: true,
	}

	err = qtx.UpdateHouseById(ctx, UpdateHouseByIdParams{
		Occupied:   true,
		ID:         args.HouseID,
		Location:   house.Location,
		Block:      house.Block,
		Partition:  house.Partition,
		Version:    house.Version,
		Occupiedby: occupiedBy,
	})

	if err != nil {

		return err
	}

	return tx.Commit()

}

func (store *SQLStore) TxnUpdateTenantHouse(ctx context.Context, args UpdateTenantParams, isDelTenant bool) error {

	tx, err := store.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	qtx := New(tx)

	var occupiedBy uuid.NullUUID

	var occupied bool

	if isDelTenant {

		args.Active = false

		occupiedBy = uuid.NullUUID{
			UUID:  args.ID,
			Valid: false,
		}

		occupied = false

	} else {

		args.Active = true

		occupiedBy = uuid.NullUUID{
			UUID:  args.ID,
			Valid: true,
		}

		occupied = true

	}

	err = qtx.UpdateTenant(ctx, args)

	if err != nil {
		return err
	}

	house, err := qtx.GetHouseById(ctx, args.HouseID)

	if err != nil {
		return err
	}

	hargs := UpdateHouseByIdParams{
		Occupied:   occupied,
		ID:         args.HouseID,
		Location:   house.Location,
		Block:      house.Block,
		Partition:  house.Partition,
		Version:    house.Version,
		Occupiedby: occupiedBy,
	}

	err = qtx.UpdateHouseById(ctx, hargs)

	if err != nil {
       return err
	}

	return tx.Commit()

}
