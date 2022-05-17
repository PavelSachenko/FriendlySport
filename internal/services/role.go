package services

import (
	"test_project/internal/models"
	"test_project/internal/repositories"
)

type Role interface {
	All() (error, []models.Role)
}

type RoleService struct {
	repo repositories.Role
}

func InitRoleService(repo repositories.Role) *RoleService {
	return &RoleService{
		repo: repo,
	}
}

func (r *RoleService) All() (error, []models.Role) {
	return r.repo.GetRoles()
}
