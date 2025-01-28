package models

import "fmt"

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
