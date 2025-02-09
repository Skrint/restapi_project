package taskService

import (
	"gorm.io/gorm"
	"log"
)

type TaskRepository interface {
	CreateTask(task Task) (Task, error)

	GetAllTasks() ([]Task, error)

	GetTasksByUserID(userID uint) ([]Task, error)

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

func (r *taskRepository) GetTasksByUserID(userID uint) ([]Task, error) {
	var tasks []Task
	err := r.db.Where("user_id = ?", userID).Find(&tasks).Error
	if err != nil {
		log.Printf("Ошибка при получении задач для user_id %d: %v", userID, err)
		return nil, err
	}
	return tasks, err
}

func (r *taskRepository) UpdateTaskByID(id uint, task Task) (Task, error) {
	var existingTask Task
	result := r.db.First(&existingTask, id) // Find existing task
	if result.Error != nil {
		return Task{}, result.Error // Task not found
	}

	existingTask.Task = task.Task
	existingTask.IsDone = task.IsDone

	result = r.db.Save(&existingTask) // Save the changes
	if result.Error != nil {
		return Task{}, result.Error
	}

	return existingTask, nil
}

func (r *taskRepository) DeleteTaskByID(id uint) error {
	var task Task
	result := r.db.Delete(&task, id).Error
	return result
}
