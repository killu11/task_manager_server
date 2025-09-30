package usecase

import (
	"context"
	"fmt"
	"log"
	"task_manager_server/internal/domain/entites"
	"task_manager_server/internal/domain/repository"
	implrepository "task_manager_server/internal/infrastructure/persistence/repository"
)

type TaskUseCase struct {
	repo repository.TaskRepository
}

func (t *TaskUseCase) Create(
	ctx context.Context,
	title, description string,
	statusID, userID int,
) error {
	task := entites.NewTask(
		title,
		description,
		userID,
		statusID,
	)

	if err := t.repo.Save(ctx, task); err != nil {
		log.Println(err)
		return fmt.Errorf("%w, failed save task", implrepository.RepoErr)
	}
	return nil
}

func (t *TaskUseCase) GetList(ctx context.Context, userID int) ([]*entites.Task, error) {
	taskSlice, err := t.repo.GetAll(ctx, userID)

	if err != nil {
		return nil, err
	}

	return taskSlice, nil
}

func (t *TaskUseCase) GetFilteredList(ctx context.Context, userID, statusID int) ([]*entites.Task, error) {
	filteredTaskSlice, err := t.repo.GetTasksByStatus(ctx, userID, statusID)

	if statusID < 1 || statusID > 3 {
		return nil, fmt.Errorf("Невалдиный статус-код! Коды существующих статусов:\n1. Не начата\n2. В процессе\n3. Завершена")
	}
	if err != nil {
		return nil, err
	}
	return filteredTaskSlice, nil
}

func (t *TaskUseCase) UpdateStatus(ctx context.Context, title string, userID, statusID int) error {
	err := t.repo.UpdateStatusByName(ctx, title, userID, statusID)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskUseCase) DeleteTask(ctx context.Context, title string, userID int) (int, error) {
	return t.repo.DeleteByName(ctx, userID, title)
}

func NewTaskUseCase(repo repository.TaskRepository) *TaskUseCase {
	return &TaskUseCase{
		repo: repo,
	}
}
