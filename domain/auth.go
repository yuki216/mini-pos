package domain

import (
	"context"
)

// Auth ...
type UserLogin struct {
	Username             string     `json:"username"`
	Email             string     `json:"email"`
	Password          string     `json:"password"`
}

type UserResponse struct {
	ID             int     `json:"id"`
	Token          string     `json:"token"`
}

type User struct {
	ID                int        `json:"id"`
	Name              string     `json:"name"`
	UserName          string     `json:"username"`
	Email             string     `json:"email"`
	Password          string     `json:"password"`
	RoleID          int     `json:"role_id"`
}

type RegisterUser struct {
	Name              string     `json:"name"`
	UserName          string     `json:"username"`
	Email             string     `json:"email"`
	Password          string     `json:"password"`
	RePassword          string     `json:"rePassword"`
}

type AuthUseCase interface {
	GetEmailUser(ctx context.Context, username string) (User, error)
	Register(ctx context.Context,register  RegisterUser) (res User,err error)
}
// AuthRepository represent the auth repository contract
type AuthRepository interface {
	GetEmailUser(ctx context.Context, username string) (User, error)
	Register(ctx context.Context,register RegisterUser) (id int, err error)
}
