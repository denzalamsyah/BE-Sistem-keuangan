package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type UserService interface {
	Register(user *models.User) (models.User, error)
	Login(user *models.User) (token *string, err error)
	GenerateResetToken(email string) (string, error)
	VerifyResetToken(email, token, newPassword, confirmPassword string) error
	GetUserByEmail(email string) (models.User, error)	
	StrongPassword(password string) bool
}

type userService struct {
	userRepo     repository.UserRepository
	sessionsRepo repository.SessionRepository
	
}

func NewUserService(userRepository repository.UserRepository, sessionsRepo repository.SessionRepository) UserService {
	return &userService{userRepository, sessionsRepo}
}

func (s *userService) GetUserByEmail(email string) (models.User, error) {
	return s.userRepo.GetUserByEmail(email)
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

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
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

// membuat token verifikasi
func (s *userService) GenerateResetToken(email string) (string, error) {
	// mendapatkan user berdasarkan email di database

	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "",err
	}

	// membuat token secara acak
	token, err := generateRandomToken()
	if err != nil {
		return "", err
	}

	// jika email tidak kosong/tersedia di database maka jalankan fungsinya
	if user.Email != "" {
		// fungsi menambahkan token verifikasi di database
		if err := s.userRepo.CreateResetToken(user.Email, token); err != nil {
			return "", err
		}
	}
	
	// fungsi mengirim verifikasi ke email
	if err := s.SendVerificationEmail(user.Email, token); err != nil {
		return "", err
	}

	return token, nil
}

// fungsi untuk reset passsword dan verifikasi token yang ada di database
func (s *userService) VerifyResetToken(email, token, newPassword, confirmPassword string) error {

	// cari token berdasarkan email
	resetToken, err := s.userRepo.GetResetTokenByEmail(email)
	if err != nil {
		return err
	}

	// cek apakaha token sudah kadaluarsa
	if time.Now().After(resetToken.CreatedAt.Add(time.Hour * 24)) {
		return errors.New("token has expired")
	}

	// cek apakah token yang diparameter sudah sesuai dengan yang didatabase
	if resetToken.TokenHash != token {
		return errors.New("invalid token")
	}
	
	// mendapatkan user berdasarkan email
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}

	// ceek apakah password baru dan konfirmasi tidak kosong
	if newPassword == "" || confirmPassword == ""{
		return errors.New("password and password confirmation cannot be empty")
	}

	// cek apakah terdiri dari 8 karakter
	if len(newPassword) < 8 {
        return errors.New("password must be at least 8 characters long")
    }

	// cek apakah password sudah sesuai spesifikasi
	if !s.userRepo.StrongPassword(newPassword){
		return errors.New("the password must consist of symbols, capital, small and numbers")
	}

	// cek apakah password dan konfirmasi sama
	if newPassword != confirmPassword{
		return errors.New("password and password confirmation do not match")
	}

	// menghash password yang diambil dari parameter
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

	// mengubah password yang didatabase 
    user.Password = string(hashedPassword)
	user.ConfirmPassword = string(hashedPassword)

	if err := s.userRepo.UpdateUser(&user); err != nil {
		return err
	}

	// menghapus token di database setelah berhasil reset password
	return s.userRepo.DeleteResetToken(resetToken)
}

// fungsi untuk mengirim verifikasi email
func (s *userService) SendVerificationEmail(email, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "adendadok52@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Link Verifikasi")
	m.SetBody("text/html", "Klik <b>LINK</b><i>di bawah ini</i>!")
	m.SetBody("text/html", "http://localhost:3000/auth/resetpass?token=" + token + "&email="+email+"\r\n")
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.gmail.com", 587, "adendadok52@gmail.com", "gmwlpcgxbzctetic")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// fungsi untuk mmebuat token secara random
func generateRandomToken() (string, error) {
	// Generate 32 bytes of random data
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Convert the random bytes to a hexadecimal string
	token := hex.EncodeToString(randomBytes)
	return token, nil
}

// fungsi untuk validasi password
func (s *userService) StrongPassword(password string) bool{
	return s.userRepo.StrongPassword(password)
}
	



