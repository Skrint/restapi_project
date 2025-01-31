package taskService

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Name   string `json:"name"`
	IsDone bool   `json:"is_done"`
}
