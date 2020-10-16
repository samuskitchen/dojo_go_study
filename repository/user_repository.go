package repository

import (
	"context"
	"dojo_go_study/model"
)

type UserRepository interface {
	GetAllUser(ctx context.Context) ([]model.User, error)
	GetOne(ctx context.Context, id uint) (model.User, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, id uint, user model.User) error
	Delete(ctx context.Context, id uint) error
}
