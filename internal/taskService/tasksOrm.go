package taskService

import (
	"gorm.io/gorm"
	"restapi_project/internal/web/users"
)

type Task struct {
	gorm.Model
	Task   string     `json:"task"`
	IsDone bool       `json:"is_done"`
	UserID uint       `json:"user_id"`
	User   users.User `gorm:"foreignKey:UserID;references:id;onDelete:CASCADE"`
}
