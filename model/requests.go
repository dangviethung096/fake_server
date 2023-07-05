package model

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AddAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Website  string `json:"website"`
}

type ListAccountRequest struct {
}
