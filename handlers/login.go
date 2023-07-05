package handlers

import (
	"fake_server/db"
	"fake_server/model"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req model.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(err)
		return
	}

	usr, err := db.GetUser(req.Username)
	if err != nil {
		c.Error(err)
		return
	}

	res := model.LoginReponse{}
	if usr.IsPasswordEqual(req.Password) {
		res.Status = true
		c.JSON(200, res)
		return
	}

	res.Status = false
	c.JSON(200, res)
}
