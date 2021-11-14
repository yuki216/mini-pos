package domain

import (
	"context"
	"time"
)

// Product ...
type Product struct {
	ID        int64     `json:"id"`
	Name     string    `json:"name" validate:"required"`
	Description   string    `json:"description" `
	Color   *string    `json:"color" `
	Size   *string    `json:"size" `
	Unit   *string    `json:"unit" `
	Image   string    `json:"image" `
	Category    Category    `json:"category"`
	UpdatedAt time.Time `json:"update_at"`
	CreatedAt time.Time `json:"create_at"`
}

type RequestProduct struct {
	ID        int64     `json:"id"`
	Name     string    `json:"name" validate:"required"`
	Description   string    `json:"description" `
	Color   *string    `json:"color" `
	Size   *string    `json:"size" `
	Unit   *string    `json:"unit" `
	Image   string    `json:"image" `
	CategoryID    int64    `json:"category_id"`
}


// ProductUsecase represent the Product's usecases
type ProductUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Product, string, error) // all access
	GetByID(ctx context.Context, id int64) (Product, error)// all access
	Update(ctx context.Context, ar *RequestProduct ) error // super admin
	GetByName(ctx context.Context, name string) (Product, error)// all access
	Store(ctx context.Context,prod  *RequestProduct ) error// super admin
	Delete(c context.Context, id int64) (err error) // super admin
}

// ProductRepository represent the Product's repository contract
type ProductRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Product, nextCursor string, err error) // super admin dan merchant
	GetByID(ctx context.Context, id int64) (Product, error) // all access
	GetByName(ctx context.Context, name string) (Product, error) // all acess
	Update(ctx context.Context, ar *RequestProduct) error  // super admin
	Store(ctx context.Context, a *RequestProduct) error //super admin
	Delete(ctx context.Context, id int64) error  // super admin
}
