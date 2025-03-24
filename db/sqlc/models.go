// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"password_hash"`
	Activated    bool      `json:"activated"`
	Version      uuid.UUID `json:"version"`
}

type House struct {
	ID        uuid.UUID `json:"id"`
	Location  string    `json:"location"`
	Block     string    `json:"block"`
	Partition int16     `json:"partition"`
	Occupied  bool      `json:"occupied"`
	Version   uuid.UUID `json:"version"`
}

type Payment struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Amount    int32     `json:"amount"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Version   uuid.UUID `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Tenant struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Phone          string    `json:"phone"`
	HouseID        uuid.UUID `json:"house_id"`
	PersonalIDType string    `json:"personal_id_type"`
	PersonalID     string    `json:"personal_id"`
	Photo          string    `json:"photo"`
	Active         bool      `json:"active"`
	Sos            time.Time `json:"sos"`
	Eos            time.Time `json:"eos"`
	Version        uuid.UUID `json:"version"`
}

type Token struct {
	Hash   []byte    `json:"hash"`
	ID     uuid.UUID `json:"id"`
	Expiry time.Time `json:"expiry"`
	Scope  string    `json:"scope"`
}
