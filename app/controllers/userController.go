package controllers

import (
	"net/http"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-gonic/gin"
)

type UserAPI interface {
	Login(c *gin.Context)
}

type userAPI struct {
	userService services.UserService
}

func NewUserAPI(userService services.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Login(c *gin.Context) {
	var loginRequest models.UserLogin
	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid decode json"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    *token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, gin.H{"message": "login success"})
}