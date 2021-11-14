package mysql


import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"mini_pos/api/v1/order/repository"
	"mini_pos/domain"
)

type mysqlOrderRepository struct {
	Conn *sql.DB
}

// NewMysqlOrderRepository will create an object that represent the Order.Repository interface
func NewMysqlOrderRepository(Conn *sql.DB) domain.OrderRepository {
	return &mysqlOrderRepository{Conn}
}

func (m *mysqlOrderRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Order, err error) {
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

	result = make([]domain.Order, 0)
	for rows.Next() {
		t := domain.Order{}
		customerID := int64(0)
		err = rows.Scan(
			&t.ID,
			&customerID,
			&t.IsCheckout,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		t.Customer = domain.Customer{
			ID: customerID,
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlOrderRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Order, nextCursor string, err error) {

	query := `SELECT id,customer_id, isCheckout, update_at, create_at
  						FROM orders WHERE %s create_at > ? ORDER BY create_at LIMIT ? `

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

func (m *mysqlOrderRepository) GetByID(ctx context.Context, id int64) (res domain.Order, err error) {

	query := `SELECT id,customer_id, isCheckout, update_at, create_at
  						FROM orders WHERE ID = ? and isCheckout=0`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Order{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *mysqlOrderRepository) GetStatus(ctx context.Context) (res int64, err error) {

	query := `SELECT id,customer_id, isCheckout, update_at, create_at
  						FROM orders WHERE isCheckout=0`

	list, err := m.fetch(ctx, query)
	if err != nil {
		return 0, err
	}

	if len(list) > 0 {
		res = list[0].ID
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *mysqlOrderRepository) Store(ctx context.Context, a *domain.Order) (id int64, err error) {
	query := `INSERT  orders 
			  SET customer_id=?, isCheckout=0,
			  update_at=? , create_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	fmt.Println(err,stmt)
	if err != nil {
		return
	}
	res, err := stmt.ExecContext(ctx,a.Customer.ID, time.Now(), time.Now())
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	a.ID = lastID
	id = lastID
	return
}

func (m *mysqlOrderRepository) Update(ctx context.Context, a *domain.Order) (err error) {

	query := `UPDATE orders 
			SET customer_id=?, isCheckout=?, update_at=?
			 WHERE ID = ? `

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx,a.Customer.ID, a.IsCheckout, time.Now(), a.ID)
	if err != nil {
		return
	}
	fmt.Println(res,err)
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
