package warehouse

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

func NewWarehouseService(broker broker.Broker, store Store) *Service {
	return &Service{broker: broker, repository: store, log: *logger.GetLogger()}
}

func (s *Service) ReservationProduct(queue string) {
	_ = s.broker.Consume(queue, "", func(message amqp.Delivery, headers map[string]interface{}) {
		var order models.Order
		if err := json.Unmarshal(message.Body, &order); err != nil {
			s.log.Error(fmt.Sprintf("Failed to decode JSON: %s", err.Error()))
			order.Status = models.FAILED
			order.Message = err.Error()
			s.orderReject(order)
			return
		}
		product, err := s.repository.GetProductByID(order.Product.ID)
		if err != nil {
			order.Status = models.FAILED
			order.Message = err.Error()
			s.orderReject(order)
			return
		}
		if order.Product.Quantity > product.Quantity {
			order.Status = models.FAILED
			order.Message = "There is not enough quantity of the product in stock."
			s.orderReject(order)
			return
		}
		product.Quantity -= order.Product.Quantity
		err = s.repository.UpdateProductQuantity(*product)
		if err != nil {
			order.Status = models.FAILED
			order.Message = err.Error()
			s.orderReject(order)
			return
		}
		s.createDelivery(order)
	})
}

func (s *Service) ReservationRejectProduct(queue string) {
	_ = s.broker.Consume(queue, "", func(message amqp.Delivery, headers map[string]interface{}) {
		var order models.Order
		if err := json.Unmarshal(message.Body, &order); err != nil {
			s.log.Error(fmt.Sprintf("Failed to decode JSON: %s", err.Error()))
			return
		}
		product, err := s.repository.GetProductByID(order.Product.ID)
		product.Quantity += order.Product.Quantity
		err = s.repository.UpdateProductQuantity(*product)
		if err != nil {
			order.Status = models.FAILED
			order.Message = err.Error()
			s.orderReject(order)
			return
		}
		s.orderReject(order)
	})
}

func (s *Service) orderReject(order models.Order) {
	err := s.broker.Publish("order_result", order, nil)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to reject order: %s", err.Error()))
	}
}

func (s *Service) createDelivery(order models.Order) {
	err := s.broker.Publish("create_delivery", order, nil)
	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to create delivery: %s", err.Error()))
	}
}
