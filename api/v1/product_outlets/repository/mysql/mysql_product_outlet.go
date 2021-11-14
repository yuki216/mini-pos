package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"mini_pos/api/v1/product_outlets/repository"
	"mini_pos/domain"
)

type mysqlProductOutletsRepository struct {
	Conn *sql.DB
}

// NewMysqlProductOutletsRepository will create an object that represent the ProductOutlets.Repository interface
func NewMysqlProductOutletsRepository(Conn *sql.DB) domain.ProductOutletsRepository {
	return &mysqlProductOutletsRepository{Conn}
}

func (m *mysqlProductOutletsRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.ProductOutlets, err error) {
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

	result = make([]domain.ProductOutlets, 0)
	for rows.Next() {
		t := domain.ProductOutlets{}
		productID := int64(0)
		outletID := int64(0)
		supplierID := int64(0)
		err = rows.Scan(
			&t.ID,
			&productID,
			&t.Sku,
			&t.Price,
			&t.QuantityStock,
			&t.QuantityUse,
			&outletID,
			&supplierID,
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

func (m *mysqlProductOutletsRepository) Fetch(ctx context.Context, cursor string, num int64, outlet int64) (res []domain.ProductOutlets, nextCursor string, err error) {
	var outletQ string
	if outlet !=0 {
		outletQ = fmt.Sprintf("outlet_id=%d and",outlet)
	}
	fmt.Println(outletQ)
	query := fmt.Sprintf(`SELECT id,product_id, sku,price, quantity,quantity_use,outlet_id,supplier_id, update_at, create_at
  						FROM outlet_products WHERE %s create_at > ? ORDER BY create_at LIMIT ? `,outletQ)

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

func (m *mysqlProductOutletsRepository) GetByID(ctx context.Context, id int64, outlet int64) (res domain.ProductOutlets, err error) {
	var outletQ string
	if outlet !=0 {
		outletQ = fmt.Sprintf("and outlet_id=%d",outlet)
	}
	query := fmt.Sprintf( `SELECT id,product_id, sku,price, quantity,quantity_use,outlet_id,supplier_id, update_at, create_at
  						FROM outlet_products WHERE ID = ? %s`,outletQ)

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.ProductOutlets{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *mysqlProductOutletsRepository) GetBySku(ctx context.Context, sku string) (res domain.ProductOutlets, err error) {
	query := `SELECT id,product_id, sku,price, quantity,quantity_use,outlet_id,supplier_id, update_at, create_at
  						FROM outlet_products WHERE sku = ?`

	list, err := m.fetch(ctx, query, sku)
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

func (m *mysqlProductOutletsRepository) Store(ctx context.Context, a *domain.RequestProductOutlets) (err error) {
	query := `INSERT  product_outlets 
			  SET sku=? ,  price=?, quantity=?, outlet_id=?, supplier_id=?,
			  update_at=? , create_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx,
		a.Sku,  a.Price, a.QuantityStock,a.OutletID,a.SupplierID, time.Now(), time.Now())
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

func (m *mysqlProductOutletsRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM outlet_products WHERE id = ? "

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

func (m *mysqlProductOutletsRepository) Update(ctx context.Context, a *domain.RequestProductOutlets) (err error) {
	var outletQ string
	if a.OutletID !=0 {
		outletQ = fmt.Sprintf("and outlet_id=%d",a.OutletID)
	}
	query := fmt.Sprintf(`UPDATE outlet_products 
			SET sku=? ,  price=?, quantity=?, outlet_id=?, supplier_id=?,
			  update_at=?
			 WHERE ID = ? %s`,outletQ)

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx,
		a.Sku, a.Price,a.QuantityStock,	a.OutletID, a.SupplierID, time.Now(), a.ID)
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

func (m *mysqlProductOutletsRepository) UpdateStock(ctx context.Context, stock domain.ProductOutletsStock) (err error) {
	var outletQ string
	if stock.OutletID !=0 {
		outletQ = fmt.Sprintf("and outlet_id=%d",stock.OutletID)
	}
	query := fmt.Sprintf(`UPDATE outlet_products 
			SET  quantity_use= quantity_use+?,
			  update_at=?
			 WHERE id = ? %s`,outletQ)

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	st , err := m.GetByID(ctx, stock.ID, int64(stock.OutletID))
	if err != nil {
		return
	}

	if st.QuantityStock >= *st.QuantityUse {
		err = errors.New("error: Quantity out of stock")
		return
	}

	res, err := stmt.ExecContext(ctx,stock.Quantity,time.Now(), stock.ID)
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