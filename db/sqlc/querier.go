// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateAdmin(ctx context.Context, arg CreateAdminParams) (CreateAdminRow, error)
	CreateHouse(ctx context.Context, arg CreateHouseParams) (uuid.UUID, error)
	CreatePayment(ctx context.Context, arg CreatePaymentParams) error
	CreateTenant(ctx context.Context, arg CreateTenantParams) error
	CreateToken(ctx context.Context, arg CreateTokenParams) error
	DeleteAllToken(ctx context.Context, arg DeleteAllTokenParams) error
	DeleteHouseById(ctx context.Context, id uuid.UUID) error
	DeletePayment(ctx context.Context, id uuid.UUID) error
	GetAdminByEmail(ctx context.Context, email string) (Admin, error)
	GetAllPayments(ctx context.Context) ([]GetAllPaymentsRow, error)
	GetDetailedPaymentById(ctx context.Context, id uuid.UUID) (GetDetailedPaymentByIdRow, error)
	GetHashTokenForAdmin(ctx context.Context, arg GetHashTokenForAdminParams) (GetHashTokenForAdminRow, error)
	GetHouseById(ctx context.Context, id uuid.UUID) (GetHouseByIdRow, error)
	GetHouses(ctx context.Context) ([]GetHousesRow, error)
	GetPaymentById(ctx context.Context, id uuid.UUID) (Payment, error)
	GetTenantById(ctx context.Context, id uuid.UUID) (Tenant, error)
	GetTenantByIdWithHouse(ctx context.Context, id uuid.UUID) (GetTenantByIdWithHouseRow, error)
	GetTenants(ctx context.Context) ([]GetTenantsRow, error)
	UpdateAdmin(ctx context.Context, arg UpdateAdminParams) (uuid.UUID, error)
	UpdateHouseById(ctx context.Context, arg UpdateHouseByIdParams) error
	UpdatePayment(ctx context.Context, arg UpdatePaymentParams) error
	UpdateTenant(ctx context.Context, arg UpdateTenantParams) error
}

var _ Querier = (*Queries)(nil)
