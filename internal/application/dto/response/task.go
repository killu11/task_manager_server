package response

import "task_manager_server/internal/domain/entites"

type TaskListResponse struct {
	Tasks *[]entites.Task
}
