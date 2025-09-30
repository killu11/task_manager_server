package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"task_manager_server/internal/domain/entites"
)

type TaskRepository struct {
	db *sql.DB
}

func (t *TaskRepository) Save(ctx context.Context, task *entites.Task) error {
	_, err := t.db.ExecContext(
		ctx,
		`INSERT INTO tasks(
                  title,
                  description,
                  user_id,
                  status_id
                  ) 
		VALUES ($1,$2,$3,$4)`,
		task.Title,
		task.Description,
		task.UserID,
		task.StatusID,
	)

	if err != nil {
		return fmt.Errorf("%w: failed save database: %v", RepoErr, err)
	}
	return nil
}

func (t *TaskRepository) GetByName(ctx context.Context, userID int, title string) (*entites.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TaskRepository) GetAll(ctx context.Context, userID int) ([]*entites.Task, error) {
	rows, err := t.db.QueryContext(
		ctx,
		`SELECT t.title, description, s.title, created_at 
			   FROM tasks as t
			   JOIN statuses as s ON s.id = t.status_id 
			   WHERE user_id = $1`,
		userID,
	)

	if err != nil {
		return nil, fmt.Errorf("%w: failed get task-list: %v", RepoErr, err)
	}
	defer rows.Close()

	taskSlice := make([]*entites.Task, 0)
	for rows.Next() {
		task := new(entites.Task)
		err = rows.Scan(
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%w: failed scan tasks: %v", RepoErr, err)
		}
		taskSlice = append(taskSlice, task)
	}

	return taskSlice, nil
}

func (t *TaskRepository) GetTasksByStatus(ctx context.Context, userID, statusID int) ([]*entites.Task, error) {
	taskSlice := make([]*entites.Task, 0)
	rows, err := t.db.QueryContext(
		ctx,
		`SELECT t.title, t.description, s.title, t.created_at 
		FROM tasks as t
		JOIN statuses as s ON s.id = t.status_id
		WHERE user_id = $1 AND status_id = $2`,
		userID, statusID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return taskSlice, nil
		}
		return nil, fmt.Errorf("%w: failed get filetred-list of tasks: %v", RepoErr, err)
	}
	defer rows.Close()
	for rows.Next() {
		task := new(entites.Task)
		err = rows.Scan(
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: failed scan rows: %v", RepoErr, err)
		}
		taskSlice = append(taskSlice, task)
	}
	return taskSlice, nil
}

func (t *TaskRepository) UpdateStatusByName(
	ctx context.Context,
	title string,
	userID, statusID int,
) error {
	res, err := t.db.ExecContext(
		ctx,
		`UPDATE tasks 
		 	   SET status_id = $1, updated_at = NOW()
			   WHERE user_id = $2 AND title = $3`,
		statusID, userID, title,
	)
	if err != nil {
		return fmt.Errorf("%w: failed update task: %v", RepoErr, err)
	}
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("%w, failed to get affected rows:%v", RepoErr, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: task not found: %v", entites.ErrTaskNotFound, err)
	}

	return nil
}

func (t *TaskRepository) DeleteByName(ctx context.Context, userID int, title string) (int, error) {
	var id int
	err := t.db.QueryRowContext(
		ctx,
		`DELETE FROM tasks 
       WHERE user_id = $1 AND title = $2 
       RETURNING id`,
		userID, title,
	).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("%w: failed delete task: %v", RepoErr, err)
	}

	return id, nil
}
func (t *TaskRepository) DeleteAll(ctx context.Context, userID int) error {
	panic("implement me")
}
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}
