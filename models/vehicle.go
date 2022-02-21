package models

type Vehicle struct {
	Id    int64  `json:"id"`
	Make  string `json:"make"`
	Model string `json:"model"`
}
