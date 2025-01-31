package taskService

import (
	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(task Task) (Task, error)

	GetAllTasks() ([]Task, error)

	UpdateTaskByID(id uint, task Task) (Task, error)

	DeleteTaskByID(id uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *taskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(task Task) (Task, error) {
	result := r.db.Create(&task)
	if result.Error != nil {
		return Task{}, result.Error
	}
	return task, nil
}

func (r *taskRepository) GetAllTasks() ([]Task, error) {
	var tasks []Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) UpdateTaskByID(id uint, task Task) (Task, error) {
	updateTask := r.db.Model(&task).Where("id = ?", id).Updates(&task)
	if updateTask.Error != nil {
		return Task{}, updateTask.Error
	}
	return task, nil
}

func (r *taskRepository) DeleteTaskByID(id uint) error {
	var task Task
	result := r.db.Delete(&task, id).Error
	return result
}
