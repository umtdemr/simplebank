package gapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/umtdemr/simplebank/db/sqlc"
	"github.com/umtdemr/simplebank/pb"
	"github.com/umtdemr/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		fmt.Println(err.Error())
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, status.Errorf(codes.AlreadyExists, "username already exists", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	fmt.Println(rsp, rsp.User.Username)
	return rsp, nil
}
