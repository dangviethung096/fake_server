package handlers

import (
	"fake_server/db"
	"fake_server/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ListAccount(c *gin.Context) {
	var req model.ListAccountRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(err)
		return
	}
	accounts, err := db.ListAccount()
	if err != nil {
		c.Error(err)
		return
	}

	res := &model.ListAccountResponse{
		Accounts: []model.Account{},
	}
	for _, account := range accounts {
		res.Accounts = append(res.Accounts, model.Account{
			Username: account.Username,
			Password: account.Password,
			Created:  account.Created.Format(time.RFC3339),
			Website:  account.Website,
		})
	}

	c.JSON(http.StatusOK, res)
}
