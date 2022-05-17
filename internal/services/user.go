package services

import (
	"test_project/config"
	"test_project/internal/models"
	"test_project/internal/repositories"
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
