package handlers

import (
	"fake_server/db"
	"fake_server/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RemoveAllAccount(c *gin.Context) {
	var req model.RemoveAllAccount
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(err)
		return
	}

	err = db.RemoveAllAccount()
	if err != nil {
		c.Error(err)
	}

	c.AbortWithStatus(http.StatusNoContent)
}
