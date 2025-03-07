package db

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
	ScopePasswordReset  = "password-reset"
)

var (
	DuplicateEmail = `pq: duplicate key value violates unique constraint "admin_email_key"`
)

type Password struct {
	Plaintext string
	Hash      []byte
}

type TokenLoc struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	AdminID   uuid.UUID `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func generateToken(id uuid.UUID, ttl time.Duration, scope string) (*TokenLoc, error) {

	token := &TokenLoc{
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

func IsValidTokenPlaintext(tokenPlaintext string) (bool, error) {
	if tokenPlaintext == "" {
		return false, errors.New("token must be provided")
	}
	if len(tokenPlaintext) != 26 {
		return false, errors.New("token must be 26 bytes long")
	}

	return true, nil
}

func (s *SQLStore) NewToken(id uuid.UUID, ttl time.Duration, scope string) (*TokenLoc, error) {
	token, err := generateToken(id, ttl, scope)
	if err != nil {
		return nil, err
	}

	args := CreateTokenParams{
		Hash:   token.Hash,
		ID:     id,
		Expiry: token.Expiry,
		Scope:  token.Scope,
	}

	// context with tome of 5 seconds
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = s.CreateToken(c, args)

	return token, err
}

func ReadUUIDParam(c echo.Context) (uuid.UUID, error) {

	id := c.Param("uuid")

	res := isValidUUID(id)
	if !res {
		return uuid.Nil, errors.New("invalid UUId parameter")
	}

	parsedID, err := uuid.Parse(id)

	if err != nil {
		return uuid.Nil, errors.New("parsing uuid failed")
	}
	return parsedID, nil
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func SetPassword(plaintextpassword string) (*Password, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextpassword), 12)

	if err != nil {
		return nil, err
	}

	pwd := &Password{
		Plaintext: plaintextpassword,
		Hash:      hash,
	}

	return pwd, nil

}

func PasswordMatches(pwd Password) (bool, error) {
	err := bcrypt.CompareHashAndPassword(pwd.Hash, []byte(pwd.Plaintext))
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
