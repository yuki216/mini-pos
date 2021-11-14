package domain

import (
	"context"
	"time"
)

// Product ...
type Customer struct {
	ID        int64     `json:"id"`
	Name     string    `json:"name" validate:"required"`
	Email   string    `json:"email" `
	Phone   string    `json:"phone" `
	Address   string    `json:"address" `
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
}



// ProductUsecase represent the Product's usecases
type CustomerUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Customer, string, error) // super admin
	GetByID(ctx context.Context, id int64) (Customer, error)// super admin
	Update(ctx context.Context, ar *Customer ) error // super admin
	GetByName(ctx context.Context, name string) (Customer, error)// super admin
	Store(ctx context.Context,prod  *Customer ) error// super admin
	Delete(c context.Context, id int64) (err error) // super admin
}

// ProductRepository represent the Product's repository contract
type CustomerRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Customer, nextCursor string, err error) //super admin
	GetByID(ctx context.Context, id int64) (Customer, error) // super admin
	GetByName(ctx context.Context, name string) (Customer, error) // super admin
	Update(ctx context.Context, ar *Customer) error  // super admin
	Store(ctx context.Context, a *Customer) error //super admin
	Delete(ctx context.Context, id int64) error  // super admin
}
