package test

import (
	"fmt"
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
	token := c.Param("token")
	account, err := database.CreateAccountByToken(token)
	if err != nil {
		c.String(404, err.Error())
		return
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
		return
	}

	if id1 < 1 || id2 < 1 {
		c.String(404, "illegal id")
		return
	}

	err = database.MakeFriendship(id1, id2)
	if err != nil {
		c.String(404, err.Error())
		return
	}

	c.String(200, "success")
}

// HandleTestRemoveFriendship handler for testing database.RemoveFriendship()
func HandleTestRemoveFriendship(c *gin.Context) {
	id1, err := strconv.Atoi(c.Param("id1"))
	if err != nil {
		c.String(404, err.Error())
		return
	}

	id2, err := strconv.Atoi(c.Param("id2"))
	if err != nil {
		c.String(404, err.Error())
		return
	}

	if id1 < 1 || id2 < 1 {
		c.String(404, "illegal id")
		return
	}

	err = database.RemoveFriendship(id1, id2)
	if err != nil {
		c.String(404, err.Error())
		return
	}

	c.String(200, "success")
}

// HandleTestUpdateFacebookFriends handler for testing database.UpdateFacebookFriends()
func HandleTestUpdateFacebookFriends(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(404, err.Error())
		return
	}

	err = database.UpdateFacebookFriends(id)
	if err != nil {
		c.String(404, err.Error())
		return
	}

	c.String(200, "OK")
}

// HandleTestGetFacebookFriends handler for testing database.GetFacebookFriends()
func HandleTestGetFacebookFriends(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(404, err.Error())
		return
	}

	friends, err := database.GetFacebookFriends(id)
	if err != nil {
		c.String(404, err.Error())
		return
	}

	c.String(200, fmt.Sprintf("%v", friends))
}
