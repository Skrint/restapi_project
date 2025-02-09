package userService

import (
	"gorm.io/gorm"
	"restapi_project/internal/taskService"
)

type User struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Tasks    []taskService.Task
}
