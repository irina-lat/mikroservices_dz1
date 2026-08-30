package converter

import (
	servicemodel "iam/internal/model"
	repomodel "iam/internal/repository/model"
)

// UserToRepo конвертирует сервисную модель в модель репозитория
func UserToRepo(user *servicemodel.User) *repomodel.User {
	if user == nil {
		return nil
	}

	return &repomodel.User{
		UUID:      user.UUID,
		Login:     user.Login,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// UserToService конвертирует модель репозитория в сервисную модель
func UserToService(user *repomodel.User) *servicemodel.User {
	if user == nil {
		return nil
	}

	return &servicemodel.User{
		UUID:      user.UUID,
		Login:     user.Login,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}