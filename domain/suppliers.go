package domain

import (
	"context"
	"time"
)

// Product ...
type Supplier struct {
	ID        int64     `json:"id"`
	Name     string    `json:"name" validate:"required"`
	Address   string    `json:"address" `
	Code   string    `json:"code" `
	UpdatedAt time.Time `json:"update_at"`
	CreatedAt time.Time `json:"create_at"`
}



// ProductUsecase represent the Product's usecases
type SupplierUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Supplier, string, error) // super admin
	GetByID(ctx context.Context, id int64) (Supplier, error)// super admin
	Update(ctx context.Context, ar *Supplier ) error // super admin
	GetByName(ctx context.Context, name string) (Supplier, error)// super admin
	Store(ctx context.Context,prod  *Supplier ) error// super admin
	Delete(c context.Context, id int64) (err error) // super admin
}

// ProductRepository represent the Product's repository contract
type SupplierRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Supplier, nextCursor string, err error) //super admin
	GetByID(ctx context.Context, id int64) (Supplier, error) // super admin
	GetByName(ctx context.Context, name string) (Supplier, error) // super admin
	Update(ctx context.Context, ar *Supplier) error  // super admin
	Store(ctx context.Context, a *Supplier) error //super admin
	Delete(ctx context.Context, id int64) error  // super admin
}
