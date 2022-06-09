package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"strconv"
	"time"
	"user_service/internal/models"
	"user_service/pkg/databases"
)

type Auth interface {
	SaveToken(userId uint64, token *models.TokenDetails) error
	DeleteToken(uuid string) error
	GetUserIdFromToken(uuid string) (error, uint64)
	IsUserExist(password, email string) (error, bool)
	CreateUser(user *models.User, roleId uint64) (error, uint64)
}

type AuthRepo struct {
	redis    *redis.Client
	postgres *databases.DB
}

func InitAuthRedis(redis *redis.Client, postgres *databases.DB) *AuthRepo {
	return &AuthRepo{
		redis:    redis,
		postgres: postgres,
	}
}

func (a AuthRepo) SaveToken(userId uint64, token *models.TokenDetails) error {
	at := time.Unix(token.AtExpires, 0)
	rt := time.Unix(token.RtExpires, 0)
	now := time.Now()

	err := a.redis.Set(token.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now)).Err()
	if err != nil {
		return err
	}

	err = a.redis.Set(token.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (a AuthRepo) DeleteToken(uuid string) error {
	deleted, err := a.redis.Del(uuid).Result()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return errors.New("deleted: 0")
	}
	return nil
}

func (a AuthRepo) GetUserIdFromToken(uuid string) (error, uint64) {
	userid, err := a.redis.Get(uuid).Result()
	if err != nil {
		return err, 0
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return nil, userID
}

func (a AuthRepo) IsUserExist(password, email string) (error, bool) {
	var id uint64
	err := a.postgres.Get(&id, "SELECT id FROM users WHERE users.password_hash=$1 AND users.email=$2 LIMIT 1", password, email)
	if err != nil || id == 0 {
		return errors.New("unauthorized"), false
	}
	return nil, true
}

func (a AuthRepo) CreateUser(user *models.User, roleId uint64) (error, uint64) {
	err := a.isUserEmailAlreadyExist(user.Email)
	if err != nil {
		return err, 0
	}
	err = a.isRoleAvailableForUser(roleId)
	if err != nil {
		return err, 0
	}

	ctx := context.Background()
	tx, err := a.postgres.BeginTx(ctx, nil)
	insertUserSql := fmt.Sprintf("INSERT INTO %s(email, password_hash) VALUES ($1, $2) RETURNING id", models.UserTable)
	rows, err := tx.QueryContext(ctx, insertUserSql, user.Email, user.PasswordHash)
	if err != nil {
		tx.Rollback()
		return err, 0
	}
	if rows.Next() {
		err = rows.Scan(&user.ID)
		if err != nil {
			tx.Rollback()
			return err, 0
		}
	}
	rows.Close()
	insertUserRoleSql := fmt.Sprintf("INSERT INTO %s(user_id, role_id) VALUES ($1, $2)", models.UserRoleTable)
	rows, err = tx.QueryContext(ctx, insertUserRoleSql, user.ID, roleId)
	if err != nil {
		tx.Rollback()
		return err, 0
	}

	if err != nil || user.ID == 0 {
		return errors.New("not inserted"), 0
	}
	err = tx.Commit()
	if err != nil {
		return err, 0
	}
	return nil, user.ID
}

func (a *AuthRepo) isUserEmailAlreadyExist(email string) error {
	var id uint64
	err := a.postgres.Get(&id, "SELECT id FROM "+models.UserTable+" WHERE "+models.UserTable+".email=$1 LIMIT 1", email)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if id != 0 {
		return errors.New("user with this email address already exists")
	}
	return nil
}

func (a *AuthRepo) isRoleAvailableForUser(roleId uint64) error {
	var roleid uint64
	err := a.postgres.Get(&roleid, "SELECT id FROM "+models.RoleTable+" WHERE "+models.RoleTable+".id=$1 AND title NOT IN ('admin') LIMIT 1", roleId)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if roleid == 0 {
		return errors.New("unavailable role for this user")
	}

	return nil
}