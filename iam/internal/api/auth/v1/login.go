package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"iam/internal/model"
	authv1 "shared/pkg/proto/auth/v1"
)

func (a *API) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if req.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	sessionUUID, err := a.service.Login(ctx, req.Login, req.Password)
	if err != nil {
		switch err {
		case model.ErrInvalidLogin:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		case model.ErrInvalidPassword:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &authv1.LoginResponse{
		SessionUuid: sessionUUID,
	}, nil
}