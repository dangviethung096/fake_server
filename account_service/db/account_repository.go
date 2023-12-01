package db

import (
	"core"
	"fake_server/service_error"
	"time"
)

type Account struct {
	Username string    `db:"username"`
	Password string    `db:"password"`
	Created  time.Time `db:"created"`
	Website  string    `db:"website"`
	Updated  time.Time `db:"updated"`
}

func (account *Account) GetTableName() string {
	return "account"
}

func (account *Account) GetPrimaryKey() string {
	return "username"
}

func GetAccount(ctx *core.Context, username string) (*Account, core.Error) {
	var account Account
	err := core.SelectByField(ctx, &account, "username", username)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func ListAccount(ctx *core.Context) ([]Account, core.Error) {
	var accounts []Account
	var account Account
	query, params, err := core.GetSelectQuery(&account)
	if err != nil {
		return nil, err
	}

	rows, errQuery := core.DBSession().QueryContext(ctx, query)
	if errQuery != nil {
		return nil, service_error.ERROR_QUERY_FAIL
	}

	for rows.Next() {
		if err := rows.Scan(params...); err != nil {
			return nil, service_error.ERROR_QUERY_FAIL
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func CreateAccount(ctx *core.Context, account *Account) core.Error {
	return core.SaveDataToDB(ctx, account)
}

func RemoveAccount(ctx *core.Context, username string) core.Error {
	account := &Account{Username: username}
	return core.DeleteDataInDB(ctx, account)
}

func RemoveAllAccount(ctx *core.Context) core.Error {
	query := `DELETE FROM account`
	_, err := core.DBSession().ExecContext(ctx, query)
	if err != nil {
		return service_error.ERROR_QUERY_FAIL
	}
	return nil
}
