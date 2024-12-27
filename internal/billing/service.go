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

func (s *Service) RequestBilling(queueOrderRequest string, queueOrderResponse string, queueNotification string) {
	_ = s.broker.Consume(queueOrderRequest, "", func(message amqp.Delivery, headers map[string]interface{}) {
		var orderMsg models.Order
		if err := json.Unmarshal(message.Body, &orderMsg); err != nil {
			s.log.Error(fmt.Sprintf("Failed to decode JSON: %s", err.Error()))
			return
		}
		billing, err := s.repository.GetBillingByUserID(orderMsg.UserID)
		if err != nil {
			orderMsg.Status = models.FAILED
			orderMsg.Message = err.Error()
		} else {
			if orderMsg.Price > billing.Account {
				orderMsg.Status = models.FAILED
				orderMsg.Message = "Not enough money"
			} else {
				billing.Account -= orderMsg.Price
				errUpdBilling := s.repository.UpdateBillingByUserID(*billing)
				if errUpdBilling != nil {
					orderMsg.Status = models.FAILED
					orderMsg.Message = err.Error()
				} else {
					orderMsg.Status = models.SUCCESS
				}
			}
		}
		notification := models.Notification{UserID: orderMsg.UserID, Status: orderMsg.Status, OrderID: orderMsg.ID}
		_ = s.broker.Publish(queueOrderResponse, orderMsg, nil)
		_ = s.broker.Publish(queueNotification, notification, nil)
	})
}

func (s *Service) CreateBilling(queue string) {
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
