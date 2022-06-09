package service

import (
	"github.com/pavel/user_service/config"
	"github.com/pavel/user_service/pkg/model"
	"github.com/pavel/user_service/pkg/repository"
)

type User interface {
	GetUser(userId uint64) (error, *model.User)
}

type UserService struct {
	repo repository.User
	cfg  *config.Config
}

func InitUserService(repo repository.User, cfg *config.Config) UserService {
	return UserService{
		repo: repo,
		cfg:  cfg,
	}
}

func (u UserService) GetUser(userId uint64) (error, *model.User) {
	return u.repo.One(userId)
}
