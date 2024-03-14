package services

import (
	"errors"
	"time"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
	"github.com/golang-jwt/jwt"
)

type UserService interface {
	Register(user *models.User) (models.User, error)
	Login(user *models.User) (token *string, err error)
	ResetPassword(email, newPassword string) error

	
}

type userService struct {
	userRepo     repository.UserRepository
	sessionsRepo repository.SessionRepository
	
}

func NewUserService(userRepository repository.UserRepository, sessionsRepo repository.SessionRepository) UserService {
	return &userService{userRepository, sessionsRepo}
}

func (s *userService) Register(user *models.User) (models.User, error) {
	dbUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return *user, err
	}

	if dbUser.Email != "" || dbUser.ID != 0 {
		return *user, errors.New("email already exists")
	}

	user.CreatedAt = time.Now()

	newUser, err := s.userRepo.CreateUser(*user)
	if err != nil {
		return *user, err
	}

	return newUser, nil
}
func (s *userService) Login(user *models.User) (token *string, err error) {
	dbUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	if dbUser.Email == "" || dbUser.ID == 0 {
		return nil, errors.New("user not found")
	}

	if user.Password != dbUser.Password {
		return nil, errors.New("wrong email or password")
	}

	expirationTime := time.Now().Add(20 * time.Minute)
	claims := &models.Claims{
		Email: dbUser.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := t.SignedString(models.JwtKey)
	if err != nil {
		return nil, err
	}

	session := models.Session{
		Token:  tokenString,
		Email:  user.Email,
		Expiry: expirationTime,
	}

	_, err = s.sessionsRepo.SessionAvailEmail(session.Email)
	if err != nil {
		err = s.sessionsRepo.AddSessions(session)
	} else {
		err = s.sessionsRepo.UpdateSessions(session)
	}

	return &tokenString, nil
}

func (s *userService) ResetPassword(email string, newPassword string) error {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	user.Password = newPassword

	if err := s.userRepo.UpdateUser(&user); err != nil {
		return err
	}
	return nil
}



// func (s *userService) VerifikasiEmail(email string) error {
//     _, err := s.userRepo.GetUserByEmail(email)
//     if err != nil {
//         return err
//     }

//     // Generate verification token (you can use UUID or any other method)
//     verificationToken := "your_verification_token"

//     // Create verification link
//     verificationLink := "https://yourdomain.com/reset-password?token=" + verificationToken

//     if err := s.emailService.SendVerificationEmail(email, verificationLink); err != nil {
//         return err
//     }

//     return nil
// }
