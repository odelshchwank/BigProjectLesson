package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/odelshchwank/BigProjectLesson/internal/core/domain"
	core_errors "github.com/odelshchwank/BigProjectLesson/internal/core/errors"
	core_postgres_pool "github.com/odelshchwank/BigProjectLesson/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	userID *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT *
	FROM todoapp.tasks
	WHERE author_user_id = $1
	ORDER BY id ASC
	LIMIT $2
	OFFSET $3;
	`

	rows, err := r.pool.Query(
		ctx,
		query,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"select tasks: %w",
			err,
		)
	}
	defer rows.Close()

	var taskModels []TaskModel

	for rows.Next() {
		var taskModel TaskModel

		err := rows.Scan(
			&taskModel.ID,
			&taskModel.Version,
			&taskModel.Title,
			&taskModel.Description,
			&taskModel.Completed,
			&taskModel.CreatedAt,
			&taskModel.CompletedAt,
			&taskModel.AuthorUserID,
		)

		if err != nil {
			if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
				return []domain.Task{}, fmt.Errorf(
					"%v: user with id=`%d`: %w",
					err,
					taskModel.AuthorUserID,
					core_errors.ErrNotFound,
				)
			}

			return nil, fmt.Errorf(
				"scan users: %w",
				err,
			)
		}

		taskModels = append(taskModels, taskModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"next rows: %w",
			err,
		)
	}

	taskDomains := taskDomainsFromModels(taskModels)

	return taskDomains, nil
}
