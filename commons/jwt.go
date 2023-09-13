package commons

import (
	"context"
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type JWTCustomClaims struct {
	ID          int    `json:"id"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	gojwt.RegisteredClaims
}

type UserData struct {
	ID          int
	FullName    string
	PhoneNumber string
}

type JWT interface {
	CreateToken(ctx context.Context, data UserData) (token string, err error)
	ValidateToken(token string) (interface{}, error)
}

type jwt struct {
	privateKey []byte
	publicKey  []byte
}

func NewJWT(privateKey []byte, publicKey []byte) JWT {
	return &jwt{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (j *jwt) CreateToken(ctx context.Context, data UserData) (token string, err error) {
	key, err := gojwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	// Set custom claims
	claims := &JWTCustomClaims{
		data.ID,
		data.FullName,
		data.PhoneNumber,
		gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	jwtToken := gojwt.NewWithClaims(gojwt.SigningMethodRS256, claims)

	// Generate encoded token and send it as response.
	token, err = jwtToken.SignedString(key)
	if err != nil {
		return
	}

	return
}

func (j *jwt) ValidateToken(token string) (interface{}, error) {
	key, err := gojwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}

	tok, err := gojwt.Parse(token, func(jwtToken *gojwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*gojwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(gojwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	data := UserData{
		ID:          int(claims["id"].(float64)),
		FullName:    claims["full_name"].(string),
		PhoneNumber: claims["phone_number"].(string),
	}

	return data, nil
}
