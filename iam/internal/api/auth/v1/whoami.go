package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"iam/internal/model"
	authv1 "shared/pkg/proto/auth/v1"
	commonv1 "shared/pkg/proto/common/v1"
)

func (a *API) Whoami(ctx context.Context, req *authv1.WhoamiRequest) (*authv1.WhoamiResponse, error) {
	if req.SessionUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "session_uuid is required")
	}

	resp, err := a.service.Whoami(ctx, req.SessionUuid)
	if err != nil {
		switch err {
		case model.ErrSessionNotFound:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		case model.ErrSessionExpired:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		case model.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &authv1.WhoamiResponse{
		Session: &commonv1.Session{
			Uuid:      resp.Session.UUID,
			CreatedAt: resp.Session.CreatedAt,
			UpdatedAt: resp.Session.UpdatedAt,
			ExpiresAt: resp.Session.ExpiresAt,
		},
		User: &commonv1.User{
			Uuid: resp.User.UUID,
			Info: &commonv1.UserInfo{
				Login: resp.User.Login,
				Email: resp.User.Email,
			},
			CreatedAt: resp.User.CreatedAt.String(),
			UpdatedAt: resp.User.UpdatedAt.String(),
		},
	}, nil
}