package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"

	"mini_pos/domain"
)

type mysqlPaymentRepository struct {
	Conn *sql.DB
}

// NewMysqlPaymentRepository will create an object that represent the Payment.Repository interface
func NewMysqlPaymentRepository(Conn *sql.DB) domain.PaymentRepository {
	return &mysqlPaymentRepository{Conn}
}

func (m *mysqlPaymentRepository) Payment(ctx context.Context, a domain.RequestPayment) (err error) {
	query := `INSERT  payments 
			  SET order_id =?, tax=?, type_payment=?, total_payment=?, update_at=?, create_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx,
		a.Order, a.Tax,a.TypePayment, a.TotalPayment, time.Now(), time.Now())
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

func (m *mysqlPaymentRepository) PaymentByCustomer(ctx context.Context, id int64) (res []domain.PaymentCustomer,err error) {
	var sQuery string
	if id !=0 {
		sQuery = fmt.Sprintf("where orders.customer_id =%d",id)
	}

	query := fmt.Sprintf( `SELECT  p.name as customer, products.name
  				FROM orders join orderItem on orderItem.order_id = orders.id join products on products.id = orderItem.product_id  join customers p on p.id = orders.customer_id  
				%s order by orders.customer_id `,sQuery)

	rows, err := m.Conn.QueryContext(ctx, query)
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

	result := make([]domain.PaymentCustomer, 0)
	for rows.Next() {
		t := domain.PaymentCustomer{}
		err = rows.Scan(
			&t.Customer,
			&t.Product,
		)

		result = append(result, t)
	}
	if len(result) > 0 {
		res = result
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *mysqlPaymentRepository) PaymentByProduct(ctx context.Context, id int64) (res []domain.PaymentProduct,err error) {
	var sQuery string
	if id !=0 {
		sQuery = fmt.Sprintf("where orderItem.product_id =%d",id)
	}

	query := fmt.Sprintf( `SELECT products.name, orderItem.qty quantity, products.unit
  				FROM orders join orderItem on orderItem.order_id = orders.id join products on products.id = orderItem.product_id   
				%s order by orderItem.product_id `,sQuery)

	rows, err := m.Conn.QueryContext(ctx, query)
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

	result := make([]domain.PaymentProduct, 0)
	for rows.Next() {
		t := domain.PaymentProduct{}
		err = rows.Scan(
			&t.Product,
			&t.Quantity,
			&t.Unit,
		)

		result = append(result, t)
	}
	if len(result) > 0 {
		res = result
	} else {
		return res, domain.ErrNotFound
	}

	return
}