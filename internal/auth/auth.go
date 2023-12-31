package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)
type ctxKey int

const AuthKey ctxKey = 1

//go:generate mockgen -source auth.go -destination auth_mock.go -package auth
type Auth interface{
	GenerateToken(claims jwt.RegisteredClaims) (string, error)
	ValidateToken(token string) (jwt.RegisteredClaims, error)
}

type auth struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewAuth(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) (Auth, error) {
	if privateKey == nil || publicKey == nil {
		err := errors.New("private/public key is not present")
		return nil, err
	}
	return &auth{privateKey: privateKey,
		publicKey: publicKey}, nil
}

func (a *auth) GenerateToken(claims jwt.RegisteredClaims) (string, error) {
	//NewWithClaims creates a new Token with the specified signing method and claims.
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Signing our token with our private key.
	tokenStr, err := tkn.SignedString(a.privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token %w", err)
	}

	return tokenStr, nil
}

func (a *auth) ValidateToken(token string) (jwt.RegisteredClaims, error) {
	var c jwt.RegisteredClaims
	// Parse the token with the registered claims.
	tkn, err := jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
		return a.publicKey, nil
	})
	if err != nil {
		return jwt.RegisteredClaims{}, fmt.Errorf("parsing token %w", err)
	}
	// Check if the parsed token is valid.
	if !tkn.Valid {
		return jwt.RegisteredClaims{}, errors.New("invalid token")
	}
	return c, nil
}
