package main

import (
	"bonbon/communicate"
	"bonbon/database"
	"github.com/gin-gonic/gin"
)

func HandleWebsocket(c *gin.Context) {
	token := c.Param("token")
	account, err := database.CreateAccountByToken(token)
	if err != nil {
		c.String(404, err.Error())
		return
	}

	communicate.ChatHandler(account.ID, c)
}
