package data

import (
	"database/sql"
)

type Models struct {
	Users    UserModel
	Tokens   TokenModel
	Houses   HouseModel
	Payments PaymentModel
	Tenants  TenantModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Tokens:   TokenModel{DB: db},
		Users:    UserModel{DB: db},
		Houses:   HouseModel{DB: db},
		Payments: PaymentModel{DB: db},
		Tenants:  TenantModel{DB: db},
	}
}
