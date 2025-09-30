package repository

import (
	"context"
	"task_manager_server/internal/domain/entites"
)

type UserRepository interface {
	Save(
		ctx context.Context,
		user *entites.User,
	) (int, error)

	GetUserIDAndHash(
		ctx context.Context,
		username string,
	) (int, []byte, error)
}
