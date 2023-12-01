package handlers

import (
	"core"
	"fake_server/db"
	"fake_server/model"
	"time"
)

func AddAccount(ctx *core.Context, request *model.AddAccountRequest) (core.HttpResponse, core.HttpError) {
	account := db.Account{
		Username: request.Username,
		Password: request.Password,
		Website:  request.Website,
	}

	err := db.CreateAccount(ctx, &account)
	if err != nil {
		return nil, core.NewDefaultHttpError(0, err.Error())
	}

	res := &model.AddAccountResponse{
		Username: request.Username,
		Password: request.Password,
		Created:  account.Created.Format(time.RFC3339),
	}

	return core.NewDefaultHttpResponse(res), nil
}
