package order

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

func NewOrderService(broker broker.Broker, store Store) *Service {
	return &Service{broker: broker, repository: store, log: *logger.GetLogger()}
}

func (s *Service) MakePayment(order models.Order) error {
	err := s.broker.Publish("oder_billing_request", order, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetPaymentResult(queue string) {
	_ = s.broker.Consume(queue, "", func(message amqp.Delivery, headers map[string]interface{}) {
		var orderBillingResponsee models.Order
		if err := json.Unmarshal(message.Body, &orderBillingResponsee); err != nil {
			s.log.Error(fmt.Sprintf("Failed to decode JSON: %s", err.Error()))
			return
		}
		err := s.repository.UpdateOrderByID(orderBillingResponsee)
		if err != nil {
			s.log.Error(fmt.Sprintf("Failed to update order: %s", err.Error()))
		}
	})
}
