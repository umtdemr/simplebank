package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey     string
	secretKeyByte []byte
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid secret key size: must be at least %v characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey: secretKey, secretKeyByte: []byte(secretKey)}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload.CreateClaims())
	return jwtToken.SignedString(maker.secretKeyByte)
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected siging method: %v", token.Header["alg"])
		}
		return maker.secretKeyByte, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpired
		}
		return nil, ErrInvalidToken
	}
	if !parsedToken.Valid {
		return nil, ErrInvalidToken
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		_, isExpirationExists := claims["exp"] // no need to check if the token expired since it's automatically checked
		if !isExpirationExists {
			return nil, ErrInvalidToken
		}

		username, isUsernameExistsOnMap := claims["username"]
		if !isUsernameExistsOnMap {
			return nil, ErrInvalidToken
		}
		exp, isExpExistsOnMap := claims["expired_at"]
		if !isExpExistsOnMap {
			return nil, ErrInvalidToken
		}
		uuidKey, isUuidExistsOnMap := claims["uuid"]
		if !isUuidExistsOnMap {
			return nil, ErrInvalidToken
		}
		iat, isIatExistsOnMap := claims["issued_at"]
		if !isIatExistsOnMap {
			return nil, ErrInvalidToken
		}
		parsedUuid, err := uuid.Parse(uuidKey.(string))
		if err != nil {
			return nil, ErrInvalidToken
		}

		parsedExpiredAt, err := time.Parse(time.RFC3339, exp.(string))
		if err != nil {
			return nil, ErrInvalidToken
		}

		parsedIssuedAt, err := time.Parse(time.RFC3339, iat.(string))
		if err != nil {
			return nil, ErrInvalidToken
		}

		return &Payload{
			ID:        parsedUuid,
			Username:  username.(string),
			ExpiredAt: parsedExpiredAt,
			IssuedAt:  parsedIssuedAt,
		}, nil
	}
	return nil, ErrInvalidToken
}
