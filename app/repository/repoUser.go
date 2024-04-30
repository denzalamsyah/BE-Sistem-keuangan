package repository

import (
	"time"
	"unicode"

	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByEmail(email string) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user *models.User) error
	CreateResetToken(email, tokenHash string, kadaluarsa time.Time) error
	GetResetTokenByEmail(email string) (*models.ResetToken, error)
	DeleteResetToken(token *models.ResetToken) error
	StrongPassword(password string) bool
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

// fmembuat token
func (r *userRepository) CreateResetToken(email, tokenHash string, kadaluarsa time.Time) error {
	resetToken := models.ResetToken{
		Email:     email,
		TokenHash: tokenHash,
		CreatedAt: time.Now(),
		ExpirationTime: kadaluarsa,
	}
	return r.db.Create(&resetToken).Error
}

//mengambil token berdasarkan email
func (r *userRepository) GetResetTokenByEmail(email string) (*models.ResetToken, error) {
	var resetToken models.ResetToken
	if err := r.db.Where("email = ?", email).First(&resetToken).Error; err != nil {
		return nil, err
	}
	return &resetToken, nil
}

// menghapus token ketika berhasil reset password
func (r *userRepository) DeleteResetToken(token *models.ResetToken) error {
    return r.db.Where("email = ? AND token_hash = ?", token.Email, token.TokenHash).Delete(token).Error
}

// fungsi untuk validasi password
func (r *userRepository) StrongPassword(password string) bool {
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}