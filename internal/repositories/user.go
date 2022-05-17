package repositories

import (
	"database/sql"
	"test_project/internal/models"
	"test_project/pkg/databases"
)

type User interface {
	One(id uint64) (error, *models.User)
}

type UserPostgres struct {
	db *databases.DB
}

func InitUserPostgres(db *databases.DB) *UserPostgres {
	return &UserPostgres{
		db: db,
	}
}

func (u UserPostgres) One(id uint64) (error, *models.User) {
	var user models.User
	rows, err := u.db.Queryx("SELECT id,email,name,description,avatar,created_at,updated_at FROM "+models.UserTable+" WHERE "+models.UserTable+".id=$1 LIMIT 1", id)
	if err != nil && err != sql.ErrNoRows {
		return err, nil
	}
	if rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Description,
			&user.Avatar,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return err, nil
		}
	}
	return nil, &user
}
