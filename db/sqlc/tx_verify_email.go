package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
)

type VerifyEmailTXParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTXResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (s *SQLStore) VerifyEmailTX(ctx context.Context, arg VerifyEmailTXParams) (VerifyEmailTXResult, error) {
	var result VerifyEmailTXResult
	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})

		if err != nil {
			return err
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username:        result.VerifyEmail.Username,
			IsEmailVerified: pgtype.Bool{Bool: true, Valid: true},
		})
		return err
	})
	return result, err
}
