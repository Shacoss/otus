package models

import (
	"encoding/json"
	"fmt"
)

type Order struct {
	ID       int64       `json:"id"`
	UserID   int64       `json:"userID"`
	Price    float64     `json:"price"`
	Product  Product     `json:"product"`
	Delivery Delivery    `json:"delivery"`
	Status   OrderStatus `json:"status"`
	Message  string      `json:"message,omitempty"`
}

func (os OrderStatus) MarshalJSON() ([]byte, error) {
	status := os.String()
	return json.Marshal(status)
}

func (os *OrderStatus) UnmarshalJSON(data []byte) error {
	var status string
	if err := json.Unmarshal(data, &status); err != nil {
		return fmt.Errorf("invalid status: %w", err)
	}
	parsedStatus, err := ParseOrderStatus(status)
	if err != nil {
		return fmt.Errorf("unsupported status: %s", status)
	}
	*os = parsedStatus
	return nil
}
