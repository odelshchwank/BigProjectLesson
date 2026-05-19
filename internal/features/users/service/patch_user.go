package users_service

import (
	"context"
	"fmt"

	"github.com/odelshchwank/BigProjectLesson/internal/core/domain"
)

func (s *UsersService) PatchUser(
	ctx context.Context,
	id int,
	patch domain.UserPatch,
) (domain.User, error) {
	// 1. get user by id
	user, err := s.usersRepostitory.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user: %w", err)
	}

	// 2. apply patch to user
	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf("apply user patch: %w", err)
	}

	// 3. save patched user in repo
	patchedUser, err := s.usersRepostitory.PatchUser(ctx, id, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("patch user: %w", err)
	}

	return patchedUser, nil
}
