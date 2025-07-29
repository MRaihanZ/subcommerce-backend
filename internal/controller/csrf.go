package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func CreateToken(c *gin.Context) {
	token := csrf.GetToken(c)
	c.JSON(200, gin.H{"code": http.StatusOK, "csrf_token": token, "error": nil})
}

func GetToken(c *gin.Context) {
	session := sessions.Default(c)
	token := session.Get("csrf_token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "csrf_token": nil, "error": "no csrf token in session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "csrf_token": token, "error": nil})
}
