package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"iam/internal/model"
	commonv1 "shared/pkg/proto/common/v1"
	userv1 "shared/pkg/proto/user/v1"
)

func (a *API) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	if req.UserUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "user_uuid is required")
	}

	user, err := a.service.GetUser(ctx, req.UserUuid)
	if err != nil {
		switch err {
		case model.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &userv1.GetUserResponse{
		User: &commonv1.User{
			Uuid: user.UUID,
			Info: &commonv1.UserInfo{
				Login: user.Login,
				Email: user.Email,
			},
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		},
	}, nil
}