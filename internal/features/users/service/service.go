package users_service

import (
	"context"

	"github.com/odelshchwank/BigProjectLesson/internal/core/domain"
)

type UsersService struct {
	usersRepostitory UsersRepostitory
}

type UsersRepostitory interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)
}

func NewUsersService(
	usersRepostitory UsersRepostitory,
) *UsersService {
	return &UsersService{
		usersRepostitory: usersRepostitory,
	}
}
