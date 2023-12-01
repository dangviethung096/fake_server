package handlers

import (
	"core"
	"fake_server/db"
	"fake_server/model"
	"net/http"
)

func RemoveAllAccount(ctx *core.Context, request *model.RemoveAccountRequest) (core.HttpResponse, core.HttpError) {
	err := db.RemoveAllAccount(ctx)
	if err != nil {
		return nil, core.NewDefaultHttpError(0, err.Error())
	}

	return core.NewHttpResponse(http.StatusNoContent, nil), nil
}
