package controllers

import (
	"log"
	"net/http"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserAPI interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	ResetPassword(c *gin.Context)
	RequestResetToken(ctx *gin.Context)
	

}

type userAPI struct {
	userService services.UserService
}

func NewUserAPI(userService services.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Register(c *gin.Context) {
	var user models.UserRegister

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":   err.Error(),})
		return
	}

	if user.Email == "" || user.Password == "" || user.ConfirmPassword == ""{
		c.JSON(http.StatusBadRequest, gin.H{"messsage": "email atau kata sandi atau konfirmasi kata sandi kosong"})
		return
	}

	if len(user.Password) < 8 {
        c.JSON(http.StatusBadRequest, gin.H{"message": "kata sandi harus terdiri dari minimal 8 karakter"})
        return
    }

	if !u.userService.StrongPassword(user.Password){
		c.JSON(http.StatusBadRequest, gin.H{"message": "password harus terdiri dari simbol, huruf kapital, kecil dan angka"})
		return	
	}

	if user.Password != user.ConfirmPassword{
		c.JSON(http.StatusBadRequest, gin.H{"message": "kata sandi dan konfirmasi kata sandi tidak cocok"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
    }

	var recordUser = models.User{
		Email:    user.Email,
		Password: string(hashedPassword),
		ConfirmPassword:  string(hashedPassword),
	}

	recordUser, err = u.userService.Register(&recordUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":   err.Error(),})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "daftar pengguna berhasil", "data": recordUser})
}
func (u *userAPI) Login(c *gin.Context) {
	var loginRequest models.UserLogin

	
	if err := c.BindJSON(&loginRequest); err != nil {
		if loginRequest.Email == "" || loginRequest.Password == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error":   "email dan password tidak boleh kosong",})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error":   err.Error(),})
		return
	}
	

	user := &models.User{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	}

	token, err := u.userService.Login(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":   err.Error(),})
		return
	}
	
	
	c.JSON(http.StatusOK, gin.H{"message": "login berhasil", "token": *token})
}

func (c *userAPI) ResetPassword(ctx *gin.Context) {
	email := ctx.PostForm("email") 
	token := ctx.PostForm("token")
	newPassword := ctx.PostForm("new_password")
	confirmPassword := ctx.PostForm("confirm_password")

	err := c.userService.VerifyResetToken(email, token, newPassword, confirmPassword) 
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "reset kata sandi berhasil"})
}

func (c *userAPI) RequestResetToken(ctx *gin.Context) {
	email := ctx.PostForm("email")
	user, err := c.userService.GetUserByEmail(email)
	if err != nil {
		return
	}

	_, err = c.userService.GenerateResetToken(user.Email)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Link verifikasi berhasil dikirim, cek email anda"})
}







