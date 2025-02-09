package userService

import (
	"gorm.io/gorm"
	"restapi_project/internal/taskService"
)

type UserRepository interface {
	CreateUser(user User) (User, error)

	GetAllUsers() ([]User, error)

	UpdateUserByID(id uint, user User) (User, error)

	DeleteUserByID(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user User) (User, error) {
	result := r.db.Create(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func (r *userRepository) GetAllUsers() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) UpdateUserByID(id uint, user User) (User, error) {
	var existingUser User
	result := r.db.First(&existingUser, id) // Find existing user
	if result.Error != nil {
		return User{}, result.Error // user not found
	}

	existingUser.Email = user.Email
	existingUser.Password = user.Password

	result = r.db.Save(&existingUser) // Save the changes
	if result.Error != nil {
		return User{}, result.Error
	}

	return existingUser, nil
}

func (r *userRepository) DeleteUserByID(id uint) error {
	var user User
	if err := r.db.First(&user, id).Error; err != nil {
		return err
	}
	var tasks []taskService.Task
	if err := r.db.Where("user_id = ?", user.ID).Find(&tasks).Error; err != nil {
		return err
	}
	errDeleteUser := r.db.Delete(&user, id).Error
	if errDeleteUser != nil {
		return errDeleteUser
	}
	for _, task := range tasks {
		if err := r.db.Delete(&task).Error; err != nil {
			return err
		}
	}
	return nil
}
