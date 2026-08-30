package v1

import (
	"iam/internal/service/user"
	userv1 "shared/pkg/proto/user/v1"
)

type API struct {
	userv1.UnimplementedUserServiceServer
	service user.Service
}

func NewAPI(service user.Service) *API {
	return &API{
		service: service,
	}
}