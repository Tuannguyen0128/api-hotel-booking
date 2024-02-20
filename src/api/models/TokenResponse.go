package models

type TokenResponse struct {
	Token      string `json:"token"`
	Token_Type string `json:"token_type"`
	Expire_In  int64  `json:"expire_in"`
}
