package token

import "time"

// Maker is an interface managing tokens
type Maker interface {
	CreateToken(username string, role string, duration time.Duration) (string, *Payload, error)
	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
