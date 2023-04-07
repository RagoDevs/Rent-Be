package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrDuplicatePhoneNumber = errors.New("duplicate phone number")
)

type Tenant struct {
	TenantId       string    `json:"tenant_id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Phone          string    `json:"phone"`
	HouseId        string    `json:"house_id"`
	PersonalIdType string    `json:"personal_id_type"`
	PersonalId     string    `json:"personal_id"`
	Photo          byte      `json:"photo"`
	Active         bool      `json:"active"`
	Eos            time.Time `json:"eos"`
}

type TenantModel struct {
	DB *sql.DB
}

func (t TenantModel) Insert(tenant *Tenant) error {
	query := `INSERT INTO TENANTS (house_id,phone,personal_id_type,personal_id,active,eos) VALUES ($1, $2, $3, $4, $5, $6)`

	args := []interface{}{tenant.HouseId, tenant.Phone, tenant.PersonalIdType, tenant.PersonalId, tenant.Active, tenant.Eos}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := t.DB.ExecContext(ctx, query, args...)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "tenants_phone_key`:
			return ErrDuplicatePhoneNumber
		default:
			return err
		}
	}

	return nil
}

func (t TenantModel) Get(tenant_id string) (*Tenant, error) {
	query := `
	    SELECT tenant_id,house_id,phone,personal_id_type,personal_id,active,eos FROM
		tenants
		WHERE tenant_id = $1`

	var tenant Tenant
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, tenant_id).Scan(
		&tenant.TenantId,
		&tenant.HouseId,
		&tenant.Phone,
		&tenant.PersonalIdType,
		&tenant.Active,
		&tenant.Eos,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &tenant, nil

}

func (t TenantModel) Update(tenant *Tenant) error {

	query := `
	UPDATE tenants 
	SET phone = $1 ,personal_id_type = $2 ,personal_id = $3 ,active = $4, eos = $5
	WHERE tenant_id = $6
	`

	args := []interface{}{tenant.Phone, tenant.PersonalIdType, tenant.PersonalId, tenant.Active, tenant.Eos}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_, err := t.DB.ExecContext(ctx, query, args...)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "tenants_phone_key`:
			return ErrDuplicatePhoneNumber

		default:
			return err

		}

	}

	return nil

}
