package v1

import (
	"iam/internal/service/auth"
	authv1 "shared/pkg/proto/auth/v1"
)

type API struct {
	authv1.UnimplementedAuthServiceServer
	service auth.Service
}

func NewAPI(service auth.Service) *API {
	return &API{
		service: service,
	}
}