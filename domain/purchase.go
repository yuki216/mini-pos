package domain

import (
	"context"
	"time"
)

// Purchase ...
type Purchase struct {
	ID        int64     `json:"id"`
	Product     Product    `json:"product"`
	Price   float64    `json:"cost_price" `
	Quantity   float64    `json:"quantity" `
	TotalPayment   float64    `json:"total_cost_price" `
	Discount   float64    `json:"discount" `
	Tax   int    `json:"tax" `
	Category    Category    `json:"category"`
	Supplier    Supplier  `json:"supplier"`
	Outlet		Outlet    `json:"outlet"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
}

type PurchaseBySupplier struct {
	Supplier        string     `json:"supplier"`
	Product     string    `json:"product"`
}

type PurchaseByProduct struct {
	Product     string    `json:"product"`
	Quantity        string     `json:"quantity"`
	Unit        string     `json:"unit"`
}

type RequestPurchase struct {
	ID        int64     `json:"id"`
	ProductID     int64    `json:"product_id"`
	Price   float64    `json:"cost_price" `
	Quantity   float64    `json:"quantity" `
	TotalPayment   float64    `json:"total_cost_price" `
	Discount   float64    `json:"discount" `
	Tax   int    `json:"tax" `
	CategoryID    int64    `json:"category_id"`
	SupplierID    int64  `json:"supplier_id"`
	OutletID		int64    `json:"outlet_id"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
}

// PurchaseUsecase represent the Purchase's usecases
type PurchaseUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Purchase, string, error) // all access
	GetByID(ctx context.Context, id int64) (Purchase, error)// all access
	Update(ctx context.Context, ar *RequestPurchase ) error // super admin
	GetBySupplierID(ctx context.Context, id int64) ([]PurchaseBySupplier, error) // super admin
	GetByProductID(ctx context.Context, id int64) ( []PurchaseByProduct, error)
	Store(ctx context.Context,prod  *RequestPurchase ) error// super admin
	Delete(c context.Context, id int64) (err error) // super admin
}

// PurchaseRepository represent the Purchase's repository contract
type PurchaseRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Purchase, nextCursor string, err error) //  super admin
	GetByID(ctx context.Context, id int64) (Purchase, error) //  super admin
	GetBySupplierID(ctx context.Context, id int64) (res []PurchaseBySupplier, err error) // super admin
	GetByProductID(ctx context.Context, id int64) (res []PurchaseByProduct, err error)
	Update(ctx context.Context, ar *RequestPurchase) error  // super admin
	Store(ctx context.Context, a *RequestPurchase) error //super admin
	Delete(ctx context.Context, id int64) error  // super admin
}
