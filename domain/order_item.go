package domain

import (
	"context"
	"time"
)

// Purchase ...
type OrderItem struct {
	ID        int64     `json:"id"`
	Order   Order    `json:"order" `
	Product   Product    `json:"product" `
	Qty   float64    `json:"qty" `
	Discount   *int64    `json:"discount" `
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
}

type RequestOrderItem struct {
	ID        int64     `json:"id"`
	OrderID   int64    `json:"order_id" `
	ProductID   int64    `json:"product_id" `
	Qty   float64    `json:"qty" `
	Discount   *int64    `json:"discount" `
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
}

// PurchaseUsecase represent the Purchase's usecases
type OrderItemUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]OrderItem, string, error) // all access
	FetchByOrderID(ctx context.Context, id int64) (res []OrderItem, err error)
	GetByID(ctx context.Context, id int64) (OrderItem, error)// all access
	Update(ctx context.Context, ar *RequestOrderItem ) error // all access
	Store(ctx context.Context,prod  *RequestOrderItem ) error//  all access
	Delete(c context.Context, id int64) (err error) // all access

}

// PurchaseRepository represent the Purchase's repository contract
type OrderItemRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []OrderItem, nextCursor string, err error) // all access
	FetchByOrderID(ctx context.Context, id int64) (res []OrderItem, err error)
	GetByID(ctx context.Context, id int64) (OrderItem, error) // all access
	GetByOrderIDAndProdID(ctx context.Context, orderId int64, prodId int64) (res OrderItem, err error)
	Store(ctx context.Context, a *RequestOrderItem) error // all access
	Update(ctx context.Context, a *RequestOrderItem) error // all access
	Delete(c context.Context, id int64) (err error) // all access
}
