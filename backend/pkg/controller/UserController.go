package controller

import "C"
import (
	"backend/pkg/models"
	"backend/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	s service.UserService
}

func NewUserController(s service.UserService) *UserController {
	return &UserController{s}
}

// @Summary      Sign up user
// @Description  Sign up user
// @Tags         auth
// @Produce      json
// @Param user body models.User required "user model"
// @Success      200  {object}  models.User
// @Failure      400  {object}  string
// @Router       /auth/sign-up/ [POST]
func (u UserController) Signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := u.s.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary      Sign in user
// @Description  Sign in user
// @Tags         auth
// @Produce      json
// @Param user body models.User required "user model"
// @Success      200  {object}  models.User
// @Failure      400  {object}  string
// @Router       /auth/sign-in/ [POST]
func (u UserController) Signin(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := u.s.GenerateToken(user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
