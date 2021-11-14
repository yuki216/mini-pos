package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"mini_pos/api/v1/supplier/repository"
	"mini_pos/domain"
)

type mysqlSupplierRepository struct {
	Conn *sql.DB
}

// NewMysqlSupplierRepository will create an object that represent the Supplier.Repository interface
func NewMysqlSupplierRepository(Conn *sql.DB) domain.SupplierRepository {
	return &mysqlSupplierRepository{Conn}
}

func (m *mysqlSupplierRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Supplier, err error) {
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

	result = make([]domain.Supplier, 0)
	for rows.Next() {
		t := domain.Supplier{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Address,
			&t.Code,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlSupplierRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Supplier, nextCursor string, err error) {
	query := `SELECT id, name, address,code, update_at, create_at
  						FROM suppliers WHERE create_at > ? ORDER BY create_at LIMIT ? `

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

func (m *mysqlSupplierRepository) GetByID(ctx context.Context, id int64) (res domain.Supplier, err error) {
	query := `SELECT id, name, address,code, update_at, create_at
  						FROM suppliers WHERE ID = ? `

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Supplier{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *mysqlSupplierRepository) GetByName(ctx context.Context, name string) (res domain.Supplier, err error) {
	query := `SELECT id, name, address,code, update_at, create_at
  						FROM Suppliers WHERE name = ?`

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

func (m *mysqlSupplierRepository) Store(ctx context.Context, a *domain.Supplier) (err error) {
	query := `INSERT  suppliers 
			  SET  name=? , address=?, code=?,
			  update_at=? , create_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx,
		a.Name,a.Address,a.Code, time.Now(), time.Now())
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

func (m *mysqlSupplierRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM suppliers WHERE id = ? "

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

func (m *mysqlSupplierRepository) Update(ctx context.Context, a *domain.Supplier) (err error) {
	query := `UPDATE suppliers 
			SET  name=? , address=?, code=?,
			  update_at=? 
			 WHERE ID = ? `

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx,
		a.Name,a.Address,a.Code, time.Now(), a.ID)
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
