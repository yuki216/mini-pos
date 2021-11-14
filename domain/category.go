package domain

import (
	"context"
	"time"
)

// Product ...
type Category struct {
	ID        int64     `json:"id"`
	Name     string    `json:"name" validate:"required"`
	Description   string    `json:"description" `
	UpdatedAt time.Time `json:"update_at"`
	CreatedAt time.Time `json:"create_at"`
}

// CategoryRepository represent the Product's repository contract
type CategoryRepository interface {
	GetByID(ctx context.Context, id int64) (Category, error)
}
