package model

type LoginReponse struct {
	Token  string `json:"token"`
	Status bool   `json:"status"`
}

type AddAccountResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Created  string `json:"created"`
}

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Created  string `json:"created"`
	Website  string `json:"website"`
}

type ListAccountResponse struct {
	Accounts []Account `json:"accounts"`
}
