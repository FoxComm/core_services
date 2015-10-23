package controllers

import (
	"net/http"

	"github.com/FoxComm/libs/configs"
	"github.com/FoxComm/core_services/user/service"
	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var loginForm LoginForm
	c.Bind(&loginForm)

	email := loginForm.Email
	password := loginForm.Password

	var u = &user.User{}
	u.InitializeWithContext(c)
	if u.First(u, "email = ?", email).Error == nil {
		if u.IsValidPassword(password) {
			c.JSON(http.StatusOK, u)
		}
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

func Logout(c *gin.Context) {
	session, _ := configs.CookieStore.Get(c.Request, "fc-admin-session")
	delete(session.Values, "Session")
	session.Save(c.Request, c.Writer)
}
