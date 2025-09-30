package repository

import (
	"context"
	"task_manager_server/internal/domain/entites"
)

type TaskRepository interface {
	Save(ctx context.Context, task *entites.Task) error
	GetByName(ctx context.Context, userID int, title string) (*entites.Task, error)
	GetAll(ctx context.Context, userID int) ([]*entites.Task, error)
	GetTasksByStatus(ctx context.Context, userID, statusID int) ([]*entites.Task, error)
	UpdateStatusByName(ctx context.Context, title string, userID, statusID int) error
	DeleteByName(ctx context.Context, userID int, title string) (int, error)
	DeleteAll(ctx context.Context, userID int) error
}
