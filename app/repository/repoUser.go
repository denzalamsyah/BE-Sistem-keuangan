package repository

import (
	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByEmail(email string) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user *models.User) error
	
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	result :=r.db.Where("email = ?", email).First(&user)
	if result.Error != nil{
		if result.Error == gorm.ErrRecordNotFound{
			return user, nil
		}
		return user, result.Error
	}
	return user, nil
}

func (r *userRepository) CreateUser(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
    if err := r.db.Save(&user).Error; err != nil {
        return err
    }
    return nil
}