package handlers

import (
	"core"
	"fake_server/db"
	"fake_server/model"
	"time"
)

func ListAccount(ctx *core.Context, request *model.ListAccountRequest) (core.HttpResponse, core.HttpError) {
	accounts, err := db.ListAccount(ctx)
	if err != nil {
		return nil, core.NewDefaultHttpError(0, err.Error())
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

	return core.NewDefaultHttpResponse(res), nil
}
