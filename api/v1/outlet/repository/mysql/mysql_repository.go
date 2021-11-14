package mysql
import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"mini_pos/domain"
)

type mysqlOutletRepo struct {
	DB *sql.DB
}

// NewMysqlOutletRepository will create an implementation of Outlet.Repository
func NewMysqlOutletRepository(db *sql.DB) domain.OutletRepository {
	return &mysqlOutletRepo{
		DB: db,
	}
}

func (m *mysqlOutletRepo) getOne(ctx context.Context, query string, args ...interface{}) (res domain.Outlet, err error) {
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.Outlet{}, err
	}
	row := stmt.QueryRowContext(ctx, args...)
	res = domain.Outlet{}

	err = row.Scan(
		&res.ID,
		&res.Name,
		&res.Address,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	return
}

func (m *mysqlOutletRepo) getAll(ctx context.Context, query string, args ...interface{}) (res []domain.Outlet, err error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
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

	result := make([]domain.Outlet, 0)
	for rows.Next() {
		t := domain.Outlet{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Address,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}
	res = result
	return
}
func (m *mysqlOutletRepo) GetByID(ctx context.Context, id int64) (domain.Outlet, error) {
	query := `SELECT id, name, address, create_at, update_at FROM outlets 
			join user_outlets on user_outlets.outlet_id = outlets.id WHERE outlets.id =?`
	return m.getOne(ctx, query,  id)
}

func (m *mysqlOutletRepo) GetByUserID(ctx context.Context, id int64, userId int64) (domain.Outlet, error) {
	query := `SELECT id, name, address, create_at, update_at FROM outlets 
			join user_outlets on user_outlets.outlet_id = outlets.id WHERE user_id=? and outlets.id =?`
	return m.getOne(ctx, query,userId,  id)
}

func (m *mysqlOutletRepo) GetByAll(ctx context.Context, userId int64) ([]domain.Outlet, error) {
	query := `SELECT id, name, address, create_at, update_at FROM outlets 
			join user_outlets on user_outlets.outlet_id = outlets.id WHERE user_id=?`
	return m.getAll(ctx, query, userId)
}