package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"mini_pos/constants"
	"time"

	"github.com/sirupsen/logrus"

	"mini_pos/api/v1/order_item/repository"
	"mini_pos/domain"
)

type mysqlOrderItemRepository struct {
	Conn *sql.DB
}

// NewMysqlOrderItemRepository will create an object that represent the OrderItem.Repository interface
func NewMysqlOrderItemRepository(Conn *sql.DB) domain.OrderItemRepository {
	return &mysqlOrderItemRepository{Conn}
}

func (m *mysqlOrderItemRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.OrderItem, err error) {
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

	result = make([]domain.OrderItem, 0)
	for rows.Next() {
		t := domain.OrderItem{}
		orderID := int64(0)
		productID := int64(0)
		err = rows.Scan(
			&t.ID,
			&orderID,
			&productID,
			&t.Qty,
			&t.Discount,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		t.Order = domain.Order{
			ID: orderID,
		}
		t.Product = domain.Product{
			ID: productID,
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlOrderItemRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.OrderItem, nextCursor string, err error) {

	query := `SELECT id,order_id, product_id, qty, discount, update_at, create_at
  						FROM orderItem WHERE  create_at > ? ORDER BY create_at LIMIT ? `

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

func (m *mysqlOrderItemRepository) FetchByOrderID(ctx context.Context, id int64) (res []domain.OrderItem, err error) {

	query := `SELECT id,order_id, product_id, qty, discount, update_at, create_at
  						FROM orderItem WHERE  order_id = ? `


	res, err = m.fetch(ctx, query, id)
	if err != nil {
		return
	}

	if len(res) == int(0) {
		err = constants.ErrDataNotFound
		return
	}

	return
}

func (m *mysqlOrderItemRepository) GetByID(ctx context.Context, id int64) (res domain.OrderItem, err error) {

	query := `SELECT id,order_id, product_id, qty, discount, update_at, create_at
  						FROM orderItem WHERE ID = ? `

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.OrderItem{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *mysqlOrderItemRepository) GetByOrderIDAndProdID(ctx context.Context, orderId int64, prodId int64) (res domain.OrderItem, err error) {

	query := `SELECT id,order_id, product_id, qty, discount, update_at, create_at
  						FROM orderItem WHERE order_id = ? and product_id=?`

	list, err := m.fetch(ctx, query, orderId, prodId)
	if err != nil {
		return domain.OrderItem{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}


func (m *mysqlOrderItemRepository) Store(ctx context.Context, a *domain.RequestOrderItem) (err error) {
	query := `INSERT  orderItem 
			  SET order_id=?, product_id=?, qty=?, discount=?, 
			   update_at=? , create_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx,a.OrderID,a.ProductID,a.Qty, a.Discount, time.Now(), time.Now())
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

func (m *mysqlOrderItemRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM orderItem WHERE id = ? "

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

func (m *mysqlOrderItemRepository) Update(ctx context.Context, a *domain.RequestOrderItem) (err error) {

	query := `UPDATE orderItem 
			SET  qty=?, update_at=? 
			 WHERE ID = ? `

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	res, err := stmt.ExecContext(ctx,a.Qty, time.Now(), a.ID)
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
