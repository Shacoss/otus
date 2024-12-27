package models

type Billing struct {
	ID      string  `json:"id"`
	UserID  int64   `json:"userID"`
	Account float64 `json:"account"`
}
