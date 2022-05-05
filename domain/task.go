package domain

import (
	"time"
)

type Task struct {
	ID          string
	UserId      string
	Name        string
	Description string
	Priority    int
	DueAt       time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TaskService interface {
	Task(id string) (*Task, error)
	Tasks() ([]*Task, error)
	CreateTask(t *Task) error
	DeleteTask(id string) error
	DeleteTasks() error
}
