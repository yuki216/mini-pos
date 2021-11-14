package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"mini_pos/api/v1/Customer/repository"
	"mini_pos/domain"
)

type mysqlCustomerRepository struct {
	Conn *sql.DB
}

// NewMysqlCustomerRepository will create an object that represent the Customer.Repository interface
func NewMysqlCustomerRepository(Conn *sql.DB) domain.CustomerRepository {
	return &mysqlCustomerRepository{Conn}
}

func (m *mysqlCustomerRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Customer, err error) {
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

	result = make([]domain.Customer, 0)
	for rows.Next() {
		t := domain.Customer{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Email,
			&t.Phone,
			&t.Address,
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

func (m *mysqlCustomerRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Customer, nextCursor string, err error) {
	query := `SELECT id, name,email, phone, address,update_at, create_at
  						FROM customers WHERE create_at > ? ORDER BY create_at LIMIT ? `

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

func (m *mysqlCustomerRepository) GetByID(ctx context.Context, id int64) (res domain.Customer, err error) {
	query := `SELECT id, name,email, phone, address,update_at, create_at
  						FROM customers WHERE ID = ? `

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Customer{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *mysqlCustomerRepository) GetByName(ctx context.Context, name string) (res domain.Customer, err error) {
	query := `SELECT id, name,email, phone, address, update_at, create_at
  						FROM customers WHERE name = ?`

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

func (m *mysqlCustomerRepository) Store(ctx context.Context, a *domain.Customer) (err error) {
	query := `INSERT  customers 
			  SET  name=? , address=?, email=?, phone=?,
			  update_at=? , create_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx,
		a.Name,a.Address,a.Email, a.Phone, time.Now(), time.Now())
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

func (m *mysqlCustomerRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM customers WHERE id = ? "

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

func (m *mysqlCustomerRepository) Update(ctx context.Context, a *domain.Customer) (err error) {
	query := `UPDATE customers 
			SET  name=? , address=?, email=?, phone=?,
			  update_at=? 
			 WHERE ID = ? `

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx,
		a.Name,a.Address,a.Email, a.Phone, time.Now(), a.ID)
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
