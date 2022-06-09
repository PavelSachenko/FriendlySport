package repositories

import (
	"fmt"
	"user_service/internal/models"
	"user_service/pkg/databases"
)

type Role interface {
	GetRoles() (error, []models.Role)
}

type RolePostgres struct {
	db *databases.DB
}

func InitRolePostgres(db *databases.DB) *RolePostgres {
	return &RolePostgres{
		db: db,
	}
}

func (r *RolePostgres) GetRoles() (error, []models.Role) {
	sql := fmt.Sprintf("SELECT * FROM %s", models.RoleTable)
	rows, err := r.db.Query(sql)
	if err != nil {
		return err, nil
	}
	var roles []models.Role
	for rows.Next() {
		role := models.Role{}
		err := rows.Scan(&role.ID, &role.Title, &role.Description)
		if err != nil {
			return err, nil
		}
		roles = append(roles, role)
	}
	return nil, roles
}
