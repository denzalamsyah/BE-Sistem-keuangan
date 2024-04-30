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
		return *user, errors.New("email sudah ada")
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
		return nil, errors.New("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return nil, errors.New("email atau password salah")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
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
		_ = s.sessionsRepo.AddSessions(session)
	} else {
		_ = s.sessionsRepo.UpdateSessions(session)
	}

	return &tokenString, nil
}

// membuat token verifikasi
func (s *userService) GenerateResetToken(email string) (string, error) {
	
	// jika email tidak kosong/tersedia di database maka jalankan fungsinya
	if email == "" {
		return "", errors.New("pengguna tidak ditemukan")
	}
	// Mengecek apakah token yang lama ada dan menghapusnya jika ya
	oldResetToken, err := s.userRepo.GetResetTokenByEmail(email)
	if err == nil{
		if err := s.userRepo.DeleteResetToken(oldResetToken); err != nil {
			return "", err
		}
	}
	// membuat token secara acak
	token, err := generateRandomToken()
	if err != nil {
		return "", err
	}
	expirationTime := time.Now().Add(2 * time.Minute)
		// fungsi menambahkan token verifikasi di database
		if err := s.userRepo.CreateResetToken(email, token, expirationTime); err != nil {
			return "", err
		}
	// fungsi mengirim verifikasi ke email
	if err := s.SendVerificationEmail(email, token); err != nil {
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

	if time.Now().After(resetToken.ExpirationTime) {
		return errors.New("token kadaluarsa")
	}

	// cek apakah token yang diparameter sudah sesuai dengan yang didatabase
	if resetToken.TokenHash != token {
		return errors.New("token tidak valid")
	}

	// ceek apakah password baru dan konfirmasi tidak kosong
	if newPassword == "" || confirmPassword == ""{
		return errors.New("kata sandi dan konfirmasi kata sandi tidak boleh kosong")
	}

	// cek apakah terdiri dari 8 karakter
	if len(newPassword) < 8 {
        return errors.New("kata sandi harus terdiri dari minimal 8 karakter")
	}
	// cek apakah password sudah sesuai spesifikasi
	if !s.userRepo.StrongPassword(newPassword){
		return errors.New("password harus terdiri dari simbol, huruf kapital, kecil dan angka")
	}

	// cek apakah password dan konfirmasi sama
	if newPassword != confirmPassword{
		return errors.New("kata sandi dan konfirmasi kata sandi tidak cocok")
	}

	// menghash password yang diambil dari parameter
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
	// mendapatkan user berdasarkan email
	user, err := s.userRepo.GetUserByEmail(email)
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
	htmlBody := `
	<html>
		<head>
			<style>
			h1 { color: #333; }
            p { font-size: 16px; color: #666; }
            h6 { font-size: 16px; color: #000000; }
            a { text-decoration: none; color: #007bff; }
			</style>
		</head>
		<body>
			<h1>Verifikasi Email <h3> Hallo `+ email +`</h3> </h1>
			<h6>Ada yang mau mencoba merubah kata sandi akun Anda, apakah ini Anda?</h6>
			<p>Jika itu Anda, klik <a href="http://localhost:3000/auth/resetpass?token=` + token + `&email=` + email + `">di sini</a> untuk mereset password Anda.</p>
		</body>
	</html>
`
	m := gomail.NewMessage()
	m.SetHeader("From", "Admin SNI <adendadok52@gmail.com>")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Link Verifikasi")
	m.SetBody("text/html", htmlBody)

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
	



