package notification

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

func NewNotificationService(broker broker.Broker, store Store) *Service {
	return &Service{broker: broker, repository: store, log: *logger.GetLogger()}
}

func (s *Service) ConsumeNotification(queueName string, consumer string) {
	_ = s.broker.Consume(queueName, consumer, func(message amqp.Delivery, headers map[string]interface{}) {
		var notificationMsg models.Notification
		if err := json.Unmarshal(message.Body, &notificationMsg); err != nil {
			s.log.Error(fmt.Sprintf("Failed to decode JSON: %s", err.Error()))
			return
		}
		err := s.repository.CreateNotification(notificationMsg)
		if err != nil {
			s.log.Error(fmt.Sprintf("Failed to create notiication: %s", err.Error()))
		}
	})
}
