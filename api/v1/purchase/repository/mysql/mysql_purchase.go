package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"mini_pos/api/v1/Purchase/repository"
	"mini_pos/domain"
)

type mysqlPurchaseRepository struct {
	Conn *sql.DB
}

// NewMysqlPurchaseRepository will create an object that represent the Purchase.Repository interface
func NewMysqlPurchaseRepository(Conn *sql.DB) domain.PurchaseRepository {
	return &mysqlPurchaseRepository{Conn}
}

func (m *mysqlPurchaseRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Purchase, err error) {
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

	result = make([]domain.Purchase, 0)
	for rows.Next() {
		t := domain.Purchase{}
		productID := int64(0)
		categoryID := int64(0)
		outletID := int64(0)
		supplierID := int64(0)
		err = rows.Scan(
			&t.ID,
			&productID,
			&t.Price,
			&t.Quantity,
			&t.TotalPayment,
			&t.Discount,
			&t.Tax,
			&categoryID,
			&supplierID,
			&outletID,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		t.Product = domain.Product{
			ID: productID,
		}
		t.Category = domain.Category{
			ID: categoryID,
		}
		t.Outlet =domain.Outlet{
			ID: outletID,
		}
		t.Supplier =domain.Supplier{
			ID: supplierID,
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlPurchaseRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Purchase, nextCursor string, err error) {

	query := `SELECT id,product_id, cost_price, quantity, total_cost_price, discount, tax, category_id, outlet_id,supplier_id, update_at, create_at
  						FROM purchases WHERE create_at > ? ORDER BY create_at LIMIT ? `

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

func (m *mysqlPurchaseRepository) GetByID(ctx context.Context, id int64) (res domain.Purchase, err error) {

	query :=  `SELECT id,product_id, cost_price, quantity, total_cost_price, discount, tax, category_id, outlet_id,supplier_id, update_at, create_at
  				FROM purchases WHERE ID = ? `

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Purchase{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *mysqlPurchaseRepository) GetBySupplierID(ctx context.Context, id int64) (res []domain.PurchaseBySupplier, err error) {
	var sQuery string
	if id !=0 {
		sQuery = fmt.Sprintf("where supplier_id =%d",id)
	}

	query := fmt.Sprintf( `SELECT s.name as supplier, products.name as product
  				FROM purchases join suppliers s on s.id = purchases.id left join products on products.id = purchases.product_id  
				%s order by supplier_id`,sQuery)

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

	result := make([]domain.PurchaseBySupplier, 0)
	for rows.Next() {
		t := domain.PurchaseBySupplier{}
		err = rows.Scan(
			&t.Supplier,
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

func (m *mysqlPurchaseRepository) GetByProductID(ctx context.Context, id int64) (res []domain.PurchaseByProduct, err error) {
	var sQuery string
	if id !=0 {
		sQuery = fmt.Sprintf("where purchases.product_id =%d",id)
	}

	query := fmt.Sprintf( `SELECT  products.name as product, o.quantity, products.unit
  				FROM purchases join products on products.id = purchases.product_id  left join outlet_products o on o.product_id = purchases.product_id  
				%s order by purchases.product_id`,sQuery)

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

	result := make([]domain.PurchaseByProduct, 0)
	for rows.Next() {
		t := domain.PurchaseByProduct{}
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

func (m *mysqlPurchaseRepository) Store(ctx context.Context, a *domain.RequestPurchase) (err error) {
	query := `INSERT  purchases 
			  SET product_id =?, cost_price=?, quantity=?, total_cost_price=?, discount=?, tax=?, category_id=?, outlet_id=?,supplier_id=?, update_at=?, create_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx,
		a.ProductID, a.Price,a.Quantity, a.TotalPayment, a.Discount, a.Tax, a.CategoryID,a.OutletID,a.SupplierID, time.Now(), time.Now())
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

func (m *mysqlPurchaseRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM purchases WHERE id = ? "

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

func (m *mysqlPurchaseRepository) Update(ctx context.Context, a *domain.RequestPurchase) (err error) {

	query := `UPDATE purchases 
			SET product_id =?, cost_price=?, quantity=?, total_cost_price=?, discount=?, tax=?, category_id=?, outlet_id=?,supplier_id=?, update_at=? 
			 WHERE ID = ? `

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx,
		a.ProductID, a.Price,a.Quantity, a.TotalPayment, a.Discount, a.Tax, a.CategoryID,a.OutletID, a.SupplierID, time.Now(), a.ID)
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


