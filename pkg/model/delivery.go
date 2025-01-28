package models

import (
	"fmt"
	"time"
)

type JSONTime time.Time

type Delivery struct {
	OrderID int64    `json:"orderID,omitempty"`
	Address string   `json:"address"`
	Date    JSONTime `json:"date"`
}

const dateFormat = "2006-01-02"

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = str[1 : len(str)-1]
	parsedTime, err := time.Parse(dateFormat, str)
	if err != nil {
		return err
	}

	*t = JSONTime(parsedTime)
	return nil
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(dateFormat))
	return []byte(stamp), nil
}

func (t JSONTime) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t JSONTime) Before(u time.Time) bool {
	return time.Time(t).Before(u)
}

func (t JSONTime) After(u time.Time) bool {
	return time.Time(t).After(u)
}

func (t JSONTime) ToDateString() string {
	return time.Time(t).Format(dateFormat)
}

func (t *JSONTime) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("unsupported Scan, expected string, got %T", value)
	}
	parsedTime, err := time.Parse(dateFormat, str)
	if err != nil {
		return fmt.Errorf("failed to parse time: %v", err)
	}
	*t = JSONTime(parsedTime)
	return nil
}
