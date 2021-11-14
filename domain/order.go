package domain

import (
	"context"
	"time"
)

// Purchase ...
type Order struct {
	ID        int64     `json:"id"`
	Customer   Customer    `json:"customer" `
	IsCheckout   int    `json:"isCheckout" `
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
}


// PurchaseUsecase represent the Purchase's usecases
//type PurchaseUsecase interface {
//	Fetch(ctx context.Context, cursor string, num int64) ([]Purchase, string, error) // all access
//	GetByID(ctx context.Context, id int64) (Purchase, error)// all access
//	Update(ctx context.Context, ar *RequestPurchase ) error // super admin
//	//GetBySku(ctx context.Context, sku string) (Purchase, error)// all access
//	Store(ctx context.Context,prod  *RequestPurchase ) error// super admin
//	Delete(c context.Context, id int64) (err error) // super admin
//}

// PurchaseRepository represent the Purchase's repository contract
type OrderRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Order, nextCursor string, err error) // all access
	GetByID(ctx context.Context, id int64) (Order, error) // all access
	GetStatus(ctx context.Context) (id int64, err error)
	Store(ctx context.Context, a *Order)  (id int64, err error)  // all access
	Update(ctx context.Context, a *Order) error // all access
}
