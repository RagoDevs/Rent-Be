package data

import (
	"database/sql"
)

type Models struct {
	Admins   AdminModel
	Tokens   TokenModel
	Houses   HouseModel
	Payments PaymentModel
	Tenants  TenantModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Tokens:   TokenModel{DB: db},
		Admins:   AdminModel{DB: db},
		Houses:   HouseModel{DB: db},
		Payments: PaymentModel{DB: db},
		Tenants:  TenantModel{DB: db},
	}
}
