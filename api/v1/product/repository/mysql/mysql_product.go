package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"mini_pos/api/v1/product/repository"
	"mini_pos/domain"
)

type mysqlProductRepository struct {
	Conn *sql.DB
}

// NewMysqlProductRepository will create an object that represent the Product.Repository interface
func NewMysqlProductRepository(Conn *sql.DB) domain.ProductRepository {
	return &mysqlProductRepository{Conn}
}

func (m *mysqlProductRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Product, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.Product, 0)
	for rows.Next() {
		t := domain.Product{}
		categoryID := int64(0)
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.Color,
			&t.Size,
			&t.Unit,
			&t.Image,
			&categoryID,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		t.Category = domain.Category{
			ID: categoryID,
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlProductRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Product, nextCursor string, err error) {

	query := `SELECT id,name, description, color, size,unit,image, category_id, update_at, create_at
  						FROM products WHERE  create_at > ? ORDER BY create_at LIMIT ? `

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return
}

func (m *mysqlProductRepository) GetByID(ctx context.Context, id int64) (res domain.Product, err error) {

	query := `SELECT id,name, description, color, size, unit,image, category_id, update_at, create_at
  						FROM products WHERE ID = ? `

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Product{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *mysqlProductRepository) GetByName(ctx context.Context, name string) (res domain.Product, err error) {
	query := `SELECT id,name, description, color, size, unit,image, category_id, update_at, create_at
  						FROM products WHERE name = ?`

	list, err := m.fetch(ctx, query, name)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}
	return
}

func (m *mysqlProductRepository) Store(ctx context.Context, a *domain.RequestProduct) (err error) {
	query := `INSERT  products 
			  SET name=? , description=?, 
              color=?, size=?, unit=?, image=?,
			  category_id=?, update_at=? , create_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx,
		a.Name,a.Description, a.Color, a.Size, a.Unit,
		a.Image, a.CategoryID,time.Now(), time.Now())
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	a.ID = lastID
	return
}

func (m *mysqlProductRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM products WHERE id = ? "

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}

	return
}

func (m *mysqlProductRepository) Update(ctx context.Context, a *domain.RequestProduct) (err error) {

	query := `UPDATE products 
			SET  name=? , description=?, 
              color=?, size=?, unit=?, image=?,
			  category_id=?,  update_at=?
			 WHERE ID = ? `

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx,
		a.Name,a.Description, a.Color, a.Size,a.Unit,
		a.Image, a.CategoryID, time.Now(), a.ID)
	if err != nil {
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}
