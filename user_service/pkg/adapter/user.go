package adapter

import (
	"github.com/pavel/user_service/pkg/logger"
	"github.com/pavel/user_service/pkg/model"
	"github.com/pavel/user_service/pkg/pb/role"
	userGRPC "github.com/pavel/user_service/pkg/pb/user"
)

func UserToGRPC(u *model.User, logger *logger.Logger) *userGRPC.User {

	user := &userGRPC.User{
		Id:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		Description: u.Description,
		Avatar:      u.Avatar,
		CreatedAt:   u.CreatedAt.Unix(),
		UpdatedAt:   u.UpdatedAt.Unix(),
		Role: &role.Role{
			Id:          u.Role.ID,
			Title:       u.Role.Title,
			Description: u.Role.Description,
		},
	}

	logger.Infof("user:one response %v", user)
	return user
}
