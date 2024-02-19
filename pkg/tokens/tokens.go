package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	db "github.com/Hopertz/rmgmt/db/sqlc"
	"github.com/Hopertz/rmgmt/pkg/validator"
	"github.com/google/uuid"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
	ScopePasswordReset  = "password-reset"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	AdminID   uuid.UUID `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func generateToken(id uuid.UUID, ttl time.Duration, scope string) (*Token, error) {

	token := &Token{
		AdminID: id,
		Expiry:  time.Now().Add(ttl),
		Scope:   scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

type TokenModel struct {
	DB *sql.DB
}

func (q *db.Queries) New(admin_id uuid.UUID, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(admin_id, ttl, scope)
	if err != nil {
		return nil, err
	}

	args := db.CreateTokenParams{
		Hash:    token.Hash,
		AdminID: admin_id,
		Expiry:  token.Expiry,
		Scope:   token.Scope,
	}
	err = m.Insert(token)
	return token, err
}
