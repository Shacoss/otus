package models

type Notification struct {
	UserID  int64       `json:"userID"`
	OrderID int64       `json:"orderID"`
	Status  OrderStatus `json:"status"`
	Message string      `json:"message,omitempty"`
}
