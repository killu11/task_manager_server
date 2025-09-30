package entites

import "time"

type Task struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	UserID      int       `json:"user_id,omitempty"`
	StatusID    int       `json:"status_id,omitempty"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewTask(
	title, desc string,
	userID, statusID int,
) *Task {
	return &Task{
		Title:       title,
		Description: desc,
		UserID:      userID,
		StatusID:    statusID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// FormatDates return's create and update format date
func (t *Task) FormatDates() (string, string) {
	return t.CreatedAt.Format(time.DateTime), t.UpdatedAt.Format(time.DateTime)

}
