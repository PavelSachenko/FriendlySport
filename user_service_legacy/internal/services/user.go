package services

import (
	"user_service/config"
	"user_service/internal/models"
	"user_service/internal/repositories"
)

type User interface {
	GetUser(userId uint64) (error, *models.User)
}

type UserService struct {
	repo repositories.User
	cfg  *config.Config
}

func InitUserService(repo repositories.User, cfg *config.Config) *UserService {
	return &UserService{
		repo: repo,
		cfg:  cfg,
	}
}

func (u UserService) GetUser(userId uint64) (error, *models.User) {
	return u.repo.One(userId)
}
