package domain

import (
	"context"
	"time"
)

// Purchase ...
type Payment struct {
	ID        int64     `json:"id"`
	Order   int64    `json:"order_id" `
	Tax   int64    `json:"tax" `
	TypePayment   string    `json:"type_payment" `
	TotalPayment   float64    `json:"total_payment" `
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
}

type RequestPayment struct {
	ID        int64     `json:"id"`
	CustomerID   int64    `json:"customer_id" `
	Order   int64    `json:"order_id" `
	Tax   int64    `json:"tax"`
	TypePayment   string    `json:"type_payment" `
	TotalPayment   float64    `json:"total_payment" `
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
}

// Purchase ...
type PaymentCustomer struct {
	Customer        string     `json:"customer"`
	Product   string    `json:"product" `
}

type PaymentProduct struct {
	Product        string     `json:"product"`
	Quantity   string    `json:"quantity" `
	Unit   string    `json:"unit" `
}

// PurchaseUsecase represent the Purchase's usecases
type PaymentUsecase interface {
	Payment(ctx context.Context, pay RequestPayment ) error // all access
	GetByCustomerID(c context.Context, id int64) ([]PaymentCustomer, error)
	GetByProductID(c context.Context, id int64) ([]PaymentProduct, error)
}

// PurchaseRepository represent the Purchase's repository contract
type PaymentRepository interface {
	Payment(c context.Context, payment RequestPayment) (err error) // all access
	PaymentByCustomer(ctx context.Context, id int64) (res []PaymentCustomer,err error)
	PaymentByProduct(ctx context.Context, id int64) (res []PaymentProduct,err error)
}
