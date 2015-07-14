package account

import (
	"bonbon/database"
	"github.com/gin-gonic/gin"
)

// SignUpHandler handler function for signing up accounts
func SignUpHandler(c *gin.Context) {
	conn := database.SQLConnection{}
	conn.Connect()
	// TODO implementation
}

// LoginHandler handler function for user login
func LoginHandler(c *gin.Context) {
	conn := database.SQLConnection{}
	conn.Connect()
	// TODO implementation
}

// LogoutHandler handler function for user logout
func LogoutHandler(c *gin.Context) {
	// TODO implementation
}
