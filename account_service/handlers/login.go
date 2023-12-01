package handlers

import (
	"core"
	"fake_server/db"
	"fake_server/model"
)

func Login(ctx *core.Context, request *model.LoginRequest) (core.HttpResponse, core.HttpError) {
	usr, err := db.GetUser(ctx, request.Username)
	if err != nil {
		return nil, core.NewDefaultHttpError(0, err.Error())
	}

	res := model.LoginReponse{}
	if usr.IsPasswordEqual(request.Password) {
		res.Status = true
		return nil, core.NewDefaultHttpError(0, "Password not match")
	}

	res.Status = false
	return core.NewDefaultHttpResponse(res), nil
}
