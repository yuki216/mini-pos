package domain

import (
	"context"
	"time"
)

// Product ...
type Outlet struct {
	ID        int64     `json:"id"`
	Name     string    `json:"name" validate:"required"`
	Address   string    `json:"address" `
	UpdatedAt time.Time `json:"update_at"`
	CreatedAt time.Time `json:"create_at"`
}

// CategoryRepository represent the Product's repository contract
type OutletRepository interface {
	GetByUserID(ctx context.Context, id int64, userId int64) (Outlet, error)
	GetByID(ctx context.Context, id int64) (Outlet, error)
    GetByAll(ctx context.Context, userId int64) ([]Outlet, error)
}
