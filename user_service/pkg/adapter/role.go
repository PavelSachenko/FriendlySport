package adapter

import (
	"github.com/pavel/user_service/pkg/logger"
	"github.com/pavel/user_service/pkg/model"
	"github.com/pavel/user_service/pkg/pb/role"
)

func RoleListToGRPC(roles []*model.Role, logger *logger.Logger) []*role.Role {
	var responseRole []*role.Role
	for _, r := range roles {
		newRole := &role.Role{Id: r.ID, Title: r.Title, Description: r.Description}
		responseRole = append(responseRole, newRole)
	}

	logger.Infof("role:all response %v", responseRole)

	return responseRole
}
