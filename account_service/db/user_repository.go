package db

import "core"

type User struct {
	Username string `db:"username"`
	Password string `db:"password"`
}

func (user *User) GetTableName() string {
	return "user"
}

func (user *User) GetPrimaryKey() string {
	return "username"
}

func (usr *User) IsPasswordEqual(password string) bool {
	return usr.Password == password
}

func GetUser(ctx *core.Context, username string) (*User, core.Error) {
	usr := &User{Username: username}
	err := core.SelectById(ctx, usr)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func CreateUser(ctx *core.Context, username string, password string) core.Error {
	user := &User{Username: username, Password: password}
	return core.SaveDataToDB(ctx, user)
}
