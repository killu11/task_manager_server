package request

type CreateTaskResponse struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	UserID      int    `json:"user_id,omitempty"`
	StatusID    int    `json:"status_id,omitempty"`
}
