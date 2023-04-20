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
	Sos            time.Time `json:"sos"`
	Eos            time.Time `json:"eos"`
}

type TenantModel struct {
	DB *sql.DB
}

func (t TenantModel) Insert(tenant *Tenant) error {
	query := `INSERT INTO TENANTS (first_name, last_name, house_id, phone, personal_id_type,personal_id,active,sos,eos) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	args := []interface{}{tenant.FirstName, tenant.LastName, tenant.HouseId, tenant.Phone, tenant.PersonalIdType, tenant.PersonalId, tenant.Active, tenant.Sos, tenant.Eos}

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
	    SELECT tenant_id, first_name, last_name, house_id, phone, personal_id_type,personal_id,active,sos,eos FROM
		tenants
		WHERE tenant_id = $1`

	var tenant Tenant
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, tenant_id).Scan(
		&tenant.TenantId,
		&tenant.FirstName,
		&tenant.LastName,
		&tenant.HouseId,
		&tenant.Phone,
		&tenant.PersonalIdType,
		&tenant.PersonalId,
		&tenant.Active,
		&tenant.Sos,
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

func (t TenantModel) GetAll() ([]*Tenant, error) {
	query := `SELECT tenant_id, first_name, last_name, house_id, phone, personal_id_type,personal_id,active,sos,eos FROM
	tenants`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	rows, err := t.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tenants := []*Tenant{}

	for rows.Next() {
		var tenant Tenant

		err := rows.Scan(
			&tenant.TenantId,
			&tenant.FirstName,
			&tenant.LastName,
			&tenant.HouseId,
			&tenant.Phone,
			&tenant.PersonalIdType,
			&tenant.PersonalId,
			&tenant.Active,
			&tenant.Sos,
			&tenant.Eos,
		)

		if err != nil {
			return nil, err
		}

		tenants = append(tenants, &tenant)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tenants, nil

}

func (t TenantModel) Update(tenant *Tenant) error {

	query := `
	UPDATE tenants 
	SET first_name = $1, last_name = $2 ,house_id = $3, phone = $4 ,personal_id_type = $5 ,personal_id = $6 ,active = $7, sos=$8 ,eos = $9
	WHERE tenant_id = $6
	`

	args := []interface{}{tenant.FirstName, tenant.LastName, tenant.HouseId, tenant.Phone, tenant.PersonalIdType, tenant.PersonalId, tenant.Active, tenant.Sos, tenant.Eos}

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
