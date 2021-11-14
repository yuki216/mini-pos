package domain

import (
	"context"
	"time"
)

// Product ...
type ProductOutlets struct {
	ID        int64     `json:"id"`
	Product  Product   `json:"product"`
	Sku     string    `json:"sku" validate:"required"`
	Price   float64    `json:"price" `
	QuantityStock   float64    `json:"quantity" `
	QuantityUse   *float64    `json:"quantity_use" `
	Outlet		Outlet    `json:"outlet"`
	Supplier    Supplier  `json:"supplier"`
	UpdatedAt time.Time `json:"update_at"`
	CreatedAt time.Time `json:"create_at"`
}

type RequestProductOutlets struct {
	ID        int64     `json:"id"`
	Sku     string    `json:"sku" validate:"required"`
	Price   float64    `json:"price"  validate:"required"`
	QuantityStock   float64    `json:"quantity"  validate:"required"`
	QuantityUse   *float64    `json:"quantity_use" `
	OutletID		int64    `json:"outlet_id"  validate:"required"`
	SupplierID   int64  `json:"supplier_id"  validate:"required"`
}

type ProductOutletsStock struct {
	ID        int64     `json:"id" validate:"required"`
	ProductID        int64     `json:"product_id" validate:"required"`
	Quantity   int    `json:"quantity_use"  validate:"required"`
	OutletID   int    `json:"outlet_id"`
}

// ProductUsecase represent the Product's usecases
type ProductOutletsUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64, outlet int64) ([]ProductOutlets, string, error) // all access
	GetByID(ctx context.Context, id int64, outlet int64) (ProductOutlets, error)// all access
	Update(ctx context.Context, ar *RequestProductOutlets ) error // super admin
	UpdateStock(c context.Context, stock ProductOutletsStock) (err error)// all access
	GetBySku(ctx context.Context, sku string) (ProductOutlets, error)// all access
	Store(ctx context.Context,prod  *RequestProductOutlets ) error// super admin
	Delete(c context.Context, id int64,  outlet int64) (err error) // super admin
}

// ProductRepository represent the Product's repository contract
type ProductOutletsRepository interface {
	Fetch(ctx context.Context, cursor string, num int64,outlet int64) (res []ProductOutlets, nextCursor string, err error) // super admin dan merchant
	GetByID(ctx context.Context, id int64,  outlet int64) (ProductOutlets, error) // all access
	GetBySku(ctx context.Context, sku string) (ProductOutlets, error) // all acess
	Update(ctx context.Context, ar *RequestProductOutlets) error  // super admin
	UpdateStock(ctx context.Context, stock ProductOutletsStock) (err error) //all access
	Store(ctx context.Context, a *RequestProductOutlets) error //super admin
	Delete(ctx context.Context, id int64) error  // super admin
}
