package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"hmgt.hopertz.me/internal/validator"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrEditConflict   = errors.New("edit conflict")
	ErrRecordNotFound = errors.New("record not found")
)

var AnonymousAdmin = &Admin{}

type Admin struct {
	AdminID   string    `json:"admin_id"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextpassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextpassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextpassword
	p.hash = hash

	return nil

}

func (p *password) Matches(plaintextpassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextpassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}

	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}
func ValidateUser(v *validator.Validator, admin *Admin) {

	ValidateEmail(v, admin.Email)
	if admin.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *admin.Password.plaintext)
	}

	if admin.Password.hash == nil {
		panic("missing password hash for admin ")
	}
}

type AdminModel struct {
	DB *sql.DB
}

func (a AdminModel) Insert(admin *Admin) error {
	query := `
	INSERT INTO admins (email, password_hash, activated)
	VALUES ($1, $2, $3 )
	RETURNING admin_id, created_at, version`

	args := []interface{}{admin.Email, admin.Password.hash, admin.Activated}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, args...).Scan(&admin.AdminID, &admin.CreatedAt, &admin.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "admins_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil

}

func (a AdminModel) GetByEmail(email string) (*Admin, error) {
	query := `
SELECT uuid, created_at, email, password_hash, activated, version
FROM admins
WHERE email = $1`
	var admin Admin
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := a.DB.QueryRowContext(ctx, query, email).Scan(
		&admin.AdminID,
		&admin.CreatedAt,
		&admin.Email,
		&admin.Password.hash,
		&admin.Activated,
		&admin.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &admin, nil
}

func (a AdminModel) Update(admin *Admin) error {
	query := `
UPDATE admins
SET email = $1, password_hash = $2, activated = $3, version = version + 1
WHERE admin_id = $4 AND version = $5
RETURNING version`
	args := []interface{}{
		admin.Email,
		admin.Password.hash,
		admin.Activated,
		admin.AdminID,
		admin.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := a.DB.QueryRowContext(ctx, query, args...).Scan(&admin.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "admins_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (a AdminModel) GetForToken(tokenScope, tokenPlaintext string) (*Admin, error) {

	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	// Set up the SQL query.
	query := `
	SELECT admins.admin_id, admins.created_at,admins.email, admins.password_hash, admins.activated, admins.version
	FROM admins
	INNER JOIN tokens
	ON admins.admin_id = tokens.admin_id
	WHERE tokens.hash = $1
	AND tokens.scope = $2
	AND tokens.expiry > $3`

	args := []interface{}{tokenHash[:], tokenScope, time.Now()}
	var admin Admin
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, args...).Scan(
		&admin.AdminID,
		&admin.CreatedAt,
		&admin.Email,
		&admin.Password.hash,
		&admin.Activated,
		&admin.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &admin, nil
}

func (a *Admin) IsAnonymous() bool {
	return a == AnonymousAdmin
}
