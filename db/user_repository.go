package db

import "fake_server/core"

type User struct {
	Username string `db:"username"`
	password string `db:"password"`
}

func (usr *User) IsPasswordEqual(password string) bool {
	return usr.password == password
}

func GetUser(username string) (*User, error) {
	var usr User
	query := `SELECT user.username, user.password FROM user WHERE user.username = ?`
	row := core.Connection.Session.QueryRow(query, username)

	err := row.Scan(&usr.Username, &usr.password)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}

func CreateUser(username string, password string) error {
	query := `INSERT INTO user(username, password) VALUES (?, ?)`
	_, err := core.ExecDB(query, username, password)
	if err != nil {
		return err
	}

	return nil
}
