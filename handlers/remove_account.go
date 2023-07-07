package handlers

import (
	"fake_server/db"
	"fake_server/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RemoveAccount(c *gin.Context) {
	var req model.RemoveAccountRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(err)
		return
	}

	account, err := db.GetAccount(req.Username)
	if err != nil {
		c.Error(err)
	}

	err = db.RemoveAccount(req.Username)
	if err != nil {
		c.Error(err)
	}

	res := &model.RemoveAccountResponse{
		Username: account.Username,
		Password: account.Password,
		Created:  account.Created.Format(time.RFC3339),
		Website:  account.Website,
	}
	c.JSON(http.StatusOK, res)
}
