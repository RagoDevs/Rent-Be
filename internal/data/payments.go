package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Payment struct {
	PaymentId string    `json:"payment_id"`
	TenantId  string    `json:"tenant_id"`
	Period    int       `json:"period"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type PaymentModel struct {
	DB *sql.DB
}

func (p PaymentModel) Insert(payment *Payment) error {
	query := fmt.Sprintf(`INSERT INTO payments (payment_id,tenant_id,period,start_date,end_date) VALUES (%s,$1,$2,$3,$4)`, "uuid_generate_v4()")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []interface{}{payment.TenantId, payment.Period, payment.StartDate, payment.EndDate}

	_, err := p.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}

	return nil

}

func (p PaymentModel) Get(payment_id string) (*Payment, error) {
	query := `SELECT payment_id,tenant_id, period, start_date, end_date FROM payments 
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

func (p PaymentModel) Update(payment Payment) error {
	query := `UPDATE payments
	SET period = $1, end_date = $2
	WHERE payment_id = $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []interface{}{payment.Period, payment.EndDate, payment.PaymentId}

	_, err := p.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}

	return nil
}

func (p PaymentModel) DELETE(payment Payment) error {
	query := `DELETE FROM payments WHERE payment_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, payment.PaymentId)

	if err != nil {
		return err
	}

	return nil
}
