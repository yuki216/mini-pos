package mysql

import (
	"context"
	"database/sql"

	"mini_pos/domain"
)

type mysqlCategoryRepo struct {
	DB *sql.DB
}

// NewMysqlCategoryRepository will create an implementation of Category.Repository
func NewMysqlCategoryRepository(db *sql.DB) domain.CategoryRepository {
	return &mysqlCategoryRepo{
		DB: db,
	}
}

func (m *mysqlCategoryRepo) getOne(ctx context.Context, query string, args ...interface{}) (res domain.Category, err error) {
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.Category{}, err
	}
	row := stmt.QueryRowContext(ctx, args...)
	res = domain.Category{}

	err = row.Scan(
		&res.ID,
		&res.Name,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	return
}

func (m *mysqlCategoryRepo) GetByID(ctx context.Context, id int64) (domain.Category, error) {
	query := `SELECT id, name, create_at, update_at FROM categories WHERE id=?`
	return m.getOne(ctx, query, id)
}
