package users_postgres_repository

import core_postgres_pool "github.com/odelshchwank/BigProjectLesson/internal/core/repository/postgres/pool"

type UsersRepostitory struct {
	pool core_postgres_pool.Pool
}

func NewUsersRepository(
	pool core_postgres_pool.Pool,
) *UsersRepostitory {
	return &UsersRepostitory{
		pool: pool,
	}
}
