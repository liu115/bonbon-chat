package test

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"bonbon/database"
	"bonbon/communicate"
)

// HandleTestWebsocket handler for testing websocket api
func HandleTestWebsocket(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err == nil {
		communicate.ChatHandler(id, c)
	} else {
		c.String(404, "not found")
	}
}

// HandleTestCreateAccountByToken handler for testing database.CreateAccountByToken()
func HandleTestCreateAccountByToken(c *gin.Context) {
	token := c.Param(":token")
	account, err := database.CreateAccountByToken(token)
	if err != nil {
		c.String(404, err.Error())
	}

	c.String(200, strconv.Itoa(account.ID))
}

// HandleTestMakeFriendship handler for testing database.MakeFriendship()
func HandleTestMakeFriendship(c *gin.Context) {
	id1, err := strconv.Atoi(c.Param("id1"))
	if err != nil {
		c.String(404, err.Error())
	}

	id2, err := strconv.Atoi(c.Param("id2"))
	if err != nil {
		c.String(404, err.Error())
	}

	if id1 < 1 || id2 < 1 {
		c.String(404, "illegal id")
	}

	err = database.MakeFriendship(id1, id2)
	if err != nil {
		c.String(404, err.Error())
	}

	c.String(200, "success")
}

// HandleTestRemoveFriendship handler for testing database.RemoveFriendship()
func HandleTestRemoveFriendship(c *gin.Context) {
	id1, err := strconv.Atoi(c.Param("id1"))
	if err != nil {
		c.String(404, err.Error())
	}

	id2, err := strconv.Atoi(c.Param("id2"))
	if err != nil {
		c.String(404, err.Error())
	}

	if id1 < 1 || id2 < 1 {
		c.String(404, "illegal id")
	}

	err = database.RemoveFriendship(id1, id2)
	if err != nil {
		c.String(404, err.Error())
	}

	c.String(200, "success")
}
