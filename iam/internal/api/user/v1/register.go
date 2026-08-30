package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"iam/internal/model"
	userv1 "shared/pkg/proto/user/v1"
)

func (a *API) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	if req.Info == nil {
		return nil, status.Error(codes.InvalidArgument, "user info is required")
	}

	login := req.Info.Info.Login
	email := req.Info.Info.Email
	password := req.Info.Password

	if login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	user, err := a.service.Register(ctx, login, email, password)
	if err != nil {
		switch err {
		case model.ErrUserAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &userv1.RegisterResponse{
		UserUuid: user.UUID,
	}, nil
}