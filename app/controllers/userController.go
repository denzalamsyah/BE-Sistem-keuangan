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
		c.JSON(http.StatusBadRequest, gin.H{"messsage": "email or password or password confirmation is empty"})
		return
	}

	if len(user.Password) < 8 {
        c.JSON(http.StatusBadRequest, gin.H{"message": "password must be at least 8 characters long"})
        return
    }

	if !u.userService.StrongPassword(user.Password){
		c.JSON(http.StatusBadRequest, gin.H{"message": "the password must consist of symbols, capital, small and numbers"})
		return	
	}

	if user.Password != user.ConfirmPassword{
		c.JSON(http.StatusBadRequest, gin.H{"message": "password and password confirmation do not match"})
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

	c.JSON(http.StatusCreated, gin.H{"message": "success register", "data": recordUser})
}
func (u *userAPI) Login(c *gin.Context) {
	var loginRequest models.UserLogin

	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":   err.Error(),})
		
		return
	}


	if loginRequest.Email == "" || loginRequest.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or password is empty"})
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
	
	
	c.JSON(http.StatusOK, gin.H{"message": "login success", "token": *token})
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

	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

func (c *userAPI) RequestResetToken(ctx *gin.Context) {
	email := ctx.PostForm("email")
	user, err := c.userService.GetUserByEmail(email)
	if err != nil {
		log.Printf("User not found for email %s", email)
		return
	}

	_, err = c.userService.GenerateResetToken(user.Email)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Token generated and sent"})
}







