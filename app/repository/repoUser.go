package repository

import (
	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByEmail(email string) (models.User, error)
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
	return user, nil // TODO: replace this
}