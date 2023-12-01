package handlers

import (
	"core"
	"fake_server/db"
	"fake_server/model"
	"net/http"
	"time"
)

func RemoveAccount(ctx *core.Context, request *model.RemoveAccountRequest) (core.HttpResponse, core.HttpError) {

	account, err := db.GetAccount(ctx, request.Username)
	if err != nil {
		return nil, core.NewDefaultHttpError(http.StatusBadRequest, err.Error())
	}

	err = db.RemoveAccount(ctx, request.Username)
	if err != nil {
		return nil, core.NewDefaultHttpError(http.StatusBadRequest, err.Error())
	}

	res := &model.RemoveAccountResponse{
		Username: account.Username,
		Password: account.Password,
		Created:  account.Created.Format(time.RFC3339),
		Website:  account.Website,
	}
	return core.NewDefaultHttpResponse(res), nil
}
