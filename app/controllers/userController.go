package controllers

import (
	"net/http"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-gonic/gin"
)

type UserAPI interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
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

	if user.Email == "" || user.Password == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or password  is empty"})
		return
	}

	var recordUser = models.User{
		Email:    user.Email,
		Password: user.Password,
	}

	recordUser, err := u.userService.Register(&recordUser)
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
	

	cookies := &http.Cookie{
		Name: "session_token",
		Value: *token,
		Path: "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(c.Writer, cookies)
	
	c.JSON(http.StatusOK, gin.H{"message": "login success", "token":*token})
}