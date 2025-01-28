package billing

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"otus/pkg/broker"
	"otus/pkg/logger"
	models "otus/pkg/model"
)

type Service struct {
	broker     broker.Broker
	repository Store
	log        slog.Logger
}

func NewBillingService(broker broker.Broker, store Store) *Service {
	return &Service{broker: broker, repository: store, log: *logger.GetLogger()}
}

func (s *Service) CreatePayment(queue string) {
	_ = s.broker.Consume(queue, "", func(message amqp.Delivery, headers map[string]interface{}) {
		var order models.Order
		if err := json.Unmarshal(message.Body, &order); err != nil {
			s.log.Error(fmt.Sprintf("Failed to decode JSON: %s", err.Error()))
			order.Status = models.FAILED
			order.Message = err.Error()
			s.paymentResult(order, "reject_delivery")
			return
		}
		billing, err := s.repository.GetBillingByUserID(order.UserID)
		if err != nil {
			order.Status = models.FAILED
			order.Message = err.Error()
		} else {
			if order.Price > billing.Account {
				order.Status = models.FAILED
				order.Message = "Not enough money"
			} else {
				billing.Account -= order.Price
				errUpdBilling := s.repository.UpdateBillingByUserID(*billing)
				if errUpdBilling != nil {
					order.Status = models.FAILED
					order.Message = err.Error()
				} else {
					order.Status = models.SUCCESS
				}
			}
		}
		if order.Status == models.SUCCESS {
			s.paymentResult(order, "order_result")
		} else {
			s.paymentResult(order, "reject_delivery")
		}
	})
}

func (s *Service) paymentResult(order models.Order, queue string) {
	_ = s.broker.Publish(queue, order, nil)
}

func (s *Service) CreateBillingAccount(queue string) {
	_ = s.broker.Consume(queue, "", func(message amqp.Delivery, headers map[string]interface{}) {
		var user models.User
		if err := json.Unmarshal(message.Body, &user); err != nil {
			s.log.Error(fmt.Sprintf("Failed to decode JSON: %s", err.Error()))
			return
		}
		_, err := s.repository.CreateBilling(user.ID)
		if err != nil {
			s.log.Error(fmt.Sprintf("Failed to create billing for user %s: %s", user.ID, err.Error()))
		}
	})
}
