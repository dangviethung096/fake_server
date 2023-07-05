package handlers

import (
	"fake_server/db"
	"fake_server/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AddAccount(c *gin.Context) {
	var req model.AddAccountRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(err)
		return
	}

	account := db.Account{
		Username: req.Username,
		Password: req.Password,
		Website:  req.Website,
	}
	err = db.CreateAccount(&account)
	if err != nil {
		c.Error(err)
	}

	res := &model.AddAccountResponse{
		Username: req.Username,
		Password: req.Password,
		Created:  account.Created.Format(time.RFC3339),
	}
	c.JSON(http.StatusCreated, res)
}
