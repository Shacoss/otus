package models

import (
	"encoding/json"
	"fmt"
)

type Order struct {
	ID      int64       `json:"id"`
	UserID  int64       `json:"userID"`
	Price   float64     `json:"price"`
	Status  OrderStatus `json:"status"`
	Message string      `json:"message,omitempty"`
}

type OrderStatus int

const (
	FAILED OrderStatus = iota
	SUCCESS
	PROCESSING
)

func (os OrderStatus) String() string {
	return [...]string{"FAILED", "SUCCESS", "PROCESSING"}[os]
}

func ParseOrderStatus(s string) (OrderStatus, error) {
	switch s {
	case "FAILED":
		return FAILED, nil
	case "SUCCESS":
		return SUCCESS, nil
	case "PROCESSING":
		return PROCESSING, nil
	default:
		return -1, fmt.Errorf("unknown status: %s", s)
	}
}

func (os OrderStatus) MarshalJSON() ([]byte, error) {
	status := os.String()
	return json.Marshal(status)
}

func (o *OrderStatus) UnmarshalJSON(data []byte) error {
	var status string
	if err := json.Unmarshal(data, &status); err != nil {
		return fmt.Errorf("invalid status: %w", err)
	}

	parsedStatus, err := ParseOrderStatus(status)
	if err != nil {
		return fmt.Errorf("unsupported status: %s", status)
	}

	*o = parsedStatus
	return nil
}
