package delivery

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"otus/pkg/broker"
	"otus/pkg/logger"
	models "otus/pkg/model"
	"time"
)

type Service struct {
	broker     broker.Broker
	repository Store
	log        slog.Logger
}

func NewDeliveryService(broker broker.Broker, store Store) *Service {
	return &Service{broker: broker, repository: store, log: *logger.GetLogger()}
}

func (s *Service) CreateDelivery(queue string) {
	_ = s.broker.Consume(queue, "", func(message amqp.Delivery, headers map[string]interface{}) {
		var order models.Order
		if err := json.Unmarshal(message.Body, &order); err != nil {
			s.log.Error(fmt.Sprintf("Failed to decode JSON: %s", err.Error()))
			order.Status = models.FAILED
			order.Message = err.Error()
			s.rejectDelivery(order)
			return
		}
		_, err := s.repository.CreateDelivery(order.Delivery)
		if err != nil {
			order.Status = models.FAILED
			order.Message = fmt.Sprintf("Failed to create delivery: %s", err.Error())
			s.rejectDelivery(order)
			return
		}
		isValidDelivery := isValidDate(order.Delivery.Date)
		if isValidDelivery {
			s.createPayment(order)
		} else {
			order.Status = models.FAILED
			order.Message = "Invalid delivery date"
			s.rejectDelivery(order)
		}
	})
}

func (s *Service) RejectDelivery(queue string) {
	_ = s.broker.Consume(queue, "", func(message amqp.Delivery, headers map[string]interface{}) {
		var order models.Order
		if err := json.Unmarshal(message.Body, &order); err != nil {
			s.log.Error(fmt.Sprintf("Failed to decode JSON: %s", err.Error()))
			return
		}
		s.rejectDelivery(order)
	})

}

func (s *Service) rejectDelivery(order models.Order) {
	err := s.broker.Publish("product_reservation_reject", order, nil)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to reject delivery: %s", err.Error()))
	}
}

func (s *Service) createPayment(order models.Order) {
	err := s.broker.Publish("create_payment", order, nil)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to create payment: %s", err.Error()))
	}
}

func isValidDate(date models.JSONTime) bool {
	now := time.Now()
	oneMonthLater := now.AddDate(0, 1, 0)
	return !date.Before(now) && !date.After(oneMonthLater)
}
