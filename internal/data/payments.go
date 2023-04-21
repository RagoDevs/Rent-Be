package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Payment struct {
	PaymentId string    `json:"payment_id"`
	TenantId  string    `json:"tenant_id"`
	Period    int       `json:"period"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Renewed   bool      `json:"renewed"`
}

type PaymentModel struct {
	DB *sql.DB
}

func (p PaymentModel) Insert(payment *Payment) error {
	query := `INSERT INTO payments (tenant_id,period,start_date,end_date, renewed) VALUES ($1,$2,$3,$4,$5)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []interface{}{payment.TenantId, payment.Period, payment.StartDate, payment.EndDate, payment.Renewed}

	_, err := p.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}

	return nil

}

func (p PaymentModel) Get(payment_id string) (*Payment, error) {
	query := `SELECT payment_id,tenant_id, period, start_date, end_date, renewed FROM payments 
	WHERE payment_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var payment Payment

	err := p.DB.QueryRowContext(ctx, query, payment_id).Scan(
		&payment.PaymentId,
		&payment.TenantId,
		&payment.Period,
		&payment.StartDate,
		&payment.EndDate,
		&payment.Renewed,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound

		default:
			return nil, err

		}
	}

	return &payment, nil

}


func (p PaymentModel) GetUnrenewed(tenant_id string) (*Payment, error) {
	query := `SELECT payment_id,tenant_id, period, start_date, end_date, renewed FROM payments 
	WHERE renewed = false and tenant_id = $1 `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var payment Payment

	err := p.DB.QueryRowContext(ctx, query, tenant_id).Scan(
		&payment.PaymentId,
		&payment.TenantId,
		&payment.Period,
		&payment.StartDate,
		&payment.EndDate,
		&payment.Renewed,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound

		default:
			return nil, err

		}
	}

	return &payment, nil

}

func (p PaymentModel) GetAll() ([]*Payment, error) {
	query := `SELECT payment_id,tenant_id, period, start_date, end_date, renewed FROM payments`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	rows, err := p.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	payments := []*Payment{}

	for rows.Next() {
		var payment Payment

		err := rows.Scan(
			&payment.PaymentId,
			&payment.TenantId,
			&payment.Period,
			&payment.StartDate,
			&payment.EndDate,
			&payment.Renewed,
		)

		if err != nil {
			return nil, err
		}

		payments = append(payments, &payment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil

}

func (p PaymentModel) Update(payment Payment) error {
	query := `UPDATE payments
	SET period = $1, start_date = $2, end_date = $3, renewed = $4
	WHERE payment_id = $5`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []interface{}{payment.Period, payment.StartDate, payment.EndDate, payment.Renewed, payment.PaymentId}

	_, err := p.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}

	return nil
}

func (p PaymentModel) DELETE(payment_id string) error {
	query := `DELETE FROM payments WHERE payment_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, payment_id)

	if err != nil {
		return err
	}

	return nil
}
