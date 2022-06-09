package repository

import (
	"fmt"
	"github.com/pavel/user_service/pkg/db"
	"github.com/pavel/user_service/pkg/model"
)

type Role interface {
	GetRoles() (error, []*model.Role)
}

type RolePostgres struct {
	db *db.DB
}

func InitRolePostgres(db *db.DB) *RolePostgres {
	return &RolePostgres{
		db: db,
	}
}

func (r *RolePostgres) GetRoles() (error, []*model.Role) {
	sql := fmt.Sprintf("SELECT * FROM %s", model.RoleTable)
	rows, err := r.db.Query(sql)
	if err != nil {
		return err, nil
	}
	var roles []*model.Role
	for rows.Next() {
		role := model.Role{}
		err := rows.Scan(&role.ID, &role.Title, &role.Description)
		if err != nil {
			return err, nil
		}
		roles = append(roles, &role)
	}
	return nil, roles
}
