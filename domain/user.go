package domain

import (
	"context"
	"time"

	"github.com/capstone-orang-jago/parkin-be/dto"
	"github.com/gin-gonic/gin"
)

type (
	User struct {
		ID         string
		Name       string
		Telephone  string
		Email      string
		Password   string
		Role       string
		Saldo      int64
		IsVerified bool
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}

	UserRepository interface {
		CreateUser(ctx context.Context, user *User) error
		GetAllUsers(ctx context.Context) ([]*User, error)
		GetUserByID(ctx context.Context, id string) (*User, error)
		UpdateUser(ctx context.Context, user *User, id string) error
		DeleteUser(ctx context.Context, id string) error
	}

	UserUsecase interface {
		// GetUserById(ctx context.Context, id string) (res *User, accessToken *string, err error)
		GetUserByID(ctx context.Context, id string) (res *User, err error)
		GetAllUsers(ctx context.Context) ([]*User, error)
		CreateUser(ctx context.Context, req dto.UserRequestDTO) (res *User, err error)
		UpdateUser(ctx context.Context, req dto.UserRequestDTO, id string) (res *User, err error)
		DeleteUser(ctx context.Context, id string) error
	}

	UserHandler interface {
		GetUserByID(c *gin.Context)
		GetAllUsers(c *gin.Context)
		CreateUser(c *gin.Context)
		UpdateUser(c *gin.Context)
		DeleteUser(c *gin.Context)
	}
)
