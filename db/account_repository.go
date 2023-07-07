package db

import (
	"fake_server/core"
	"time"
)

type Account struct {
	Username string    `db:"username"`
	Password string    `db:"password"`
	Created  time.Time `db:"created"`
	Website  string    `db:"website"`
	Updated  time.Time `db:"updated"`
}

func GetAccount(username string) (*Account, error) {
	var account Account
	row := core.Connection.Session.QueryRow(`SELECT username, password, created, website, updated FROM account WHERE account.username = ?`, username)

	err := row.Scan(&account.Username, &account.Password, &account.Created, &account.Website, &account.Updated)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func ListAccount() ([]Account, error) {
	var accounts []Account
	query := `SELECT username, password, created, website, updated FROM account`
	stmt, err := core.Connection.Session.Prepare(query)
	if err != nil {
		return nil, err
	} else {
		defer stmt.Close()
	}

	results, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer results.Close()

	for results.Next() {
		var account Account
		results.Scan(&account.Username, &account.Password, &account.Created, &account.Website, &account.Updated)
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func CreateAccount(account *Account) error {
	account.Created = time.Now()
	account.Updated = time.Now()
	query := `INSERT INTO account(username, password, created, updated, website) VALUES (?,?,?,?,?)`
	_, err := core.ExecDB(query, account.Username, account.Password, account.Created, account.Updated, account.Website)
	if err != nil {
		return err
	}

	return nil
}

func RemoveAccount(username string) error {
	query := `DELETE FROM account WHERE username = ?`
	_, err := core.ExecDB(query, username)
	if err != nil {
		return err
	}
	return nil
}

func RemoveAllAccount() error {
	query := `DELETE FROM account`
	_, err := core.ExecDB(query)
	if err != nil {
		return err
	}
	return nil
}
