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

	err = qtx.CreateTenant(ctx, args)

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

	err = qtx.UpdateHouseById(ctx, UpdateHouseByIdParams{
		Occupied:  true,
		ID:        args.HouseID,
		Location:  house.Location,
		Block:     house.Block,
		Partition: house.Partition,
		Version:   house.Version,
	})

	if err != nil {

		return err
	}

	return tx.Commit()

}

func (store *SQLStore) TxnUpdateTenantHouse(ctx context.Context, args UpdateTenantParams, prev_house_id uuid.UUID) error {

	tx, err := store.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	qtx := New(tx)

	// rare cases when tenant shifts houses
	if prev_house_id != args.HouseID {

		// old house tenant is moving from
		house, err := qtx.GetHouseById(ctx, prev_house_id)

		if err != nil {
			return err
		}

		hargs := UpdateHouseByIdParams{
			Occupied:  false,
			ID:        prev_house_id,
			Location:  house.Location,
			Block:     house.Block,
			Partition: house.Partition,
			Version:   house.Version,
		}

		err = qtx.UpdateHouseById(ctx, hargs)

		if err != nil {
			return err
		}

		// new house tenant is moving to

		nh, err := qtx.GetHouseById(ctx, args.HouseID)

		if err != nil {
			return err
		}

		nhargs := UpdateHouseByIdParams{
			Occupied:  true,
			ID:        args.HouseID,
			Location:  nh.Location,
			Block:     nh.Block,
			Partition: nh.Partition,
			Version:   nh.Version,
		}

		err = qtx.UpdateHouseById(ctx, nhargs)

		if err != nil {
			return err
		}

	}

	err = qtx.UpdateTenant(ctx, args)

	if err != nil {
		return err
	}

	return tx.Commit()

}

func (store *SQLStore) TxnRemoveTenantHouse(ctx context.Context, args UpdateTenantParams) error {

	tx, err := store.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	qtx := New(tx)

	args.Active = false

	err = qtx.UpdateTenant(ctx, args)

	if err != nil {
		return err
	}

	house, err := qtx.GetHouseById(ctx, args.HouseID)

	if err != nil {
		return err
	}

	hargs := UpdateHouseByIdParams{
		Occupied:  false,
		ID:        args.HouseID,
		Location:  house.Location,
		Block:     house.Block,
		Partition: house.Partition,
		Version:   house.Version,
	}

	err = qtx.UpdateHouseById(ctx, hargs)

	if err != nil {
		return err
	}

	return tx.Commit()

}
