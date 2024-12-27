package auth

import (
	"log/slog"
	"otus/pkg/broker"
	"otus/pkg/logger"
	models "otus/pkg/model"
)

type Service struct {
	broker       broker.Broker
	log          slog.Logger
	billingQueue string
}

func NewAuthService(broker broker.Broker, billingQueue string) *Service {
	return &Service{broker: broker, log: *logger.GetLogger(), billingQueue: billingQueue}
}

func (s *Service) CreateBilling(user models.User) error {
	err := s.broker.Publish(s.billingQueue, user, nil)
	if err != nil {
		return err
	}
	return nil
}
