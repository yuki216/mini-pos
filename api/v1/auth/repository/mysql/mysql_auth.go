package mysql

import (
	"context"
	"database/sql"
	"mini_pos/utils"

	"mini_pos/domain"
)

type mysqlAuthRepo struct {
	DB *sql.DB
}

// NewMysqlAuthorRepository will create an implementation of author.Repository
func NewMysqlAuthRepository(db *sql.DB) domain.AuthRepository {
	return &mysqlAuthRepo{
		DB: db,
	}
}

func (m *mysqlAuthRepo) getUser(ctx context.Context, query string, args ...interface{}) (res domain.User, err error) {
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.User{}, err
	}
	row := stmt.QueryRowContext(ctx, args...)
	res = domain.User{}

	err = row.Scan(
		&res.ID,
		&res.Name,
		&res.UserName,
		&res.Email,
		&res.Password,
		&res.RoleID,
	)

	return
}

func (m *mysqlAuthRepo) GetEmailUser(ctx context.Context, username string) (domain.User, error) {
	query := `SELECT id, fullname as name, username, email, password, role_id FROM users WHERE email=? `
	return m.getUser(ctx, query, username)
}

func (m *mysqlAuthRepo) Register(ctx context.Context, a domain.RegisterUser) (id int, err error) {
	query := `INSERT  users SET email=?, password=?, username=?, fullname=?`
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	password, err := utils.HashPassword(a.Password)
	if err != nil {
		return
	}
	res, err := stmt.ExecContext(ctx, a.Email, password, a.UserName, a.Name)
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	id = int(lastID)
	return
}