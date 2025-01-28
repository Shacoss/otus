package order

import (
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"net/http"
	"otus/pkg/broker"
	"otus/pkg/logger"
	models "otus/pkg/model"
	"strconv"
)

type Service struct {
	broker     broker.Broker
	repository Store
	productURL string
	log        slog.Logger
}

func NewOrderService(broker broker.Broker, store Store, productURL string) *Service {
	return &Service{broker: broker, repository: store, productURL: fmt.Sprintf("%s/%s", productURL, "/warehouse/product"), log: *logger.GetLogger()}
}

func (s *Service) GetProductByID(productID int64, userID int64) (*models.Product, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%d", s.productURL, productID), nil)
	if err != nil {
		s.log.Error(fmt.Sprintf("Error creating request to warehouse api: %s", err.Error()))
		return nil, err
	}
	req.Header.Set("X-UserID", strconv.FormatInt(userID, 10))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.log.Error(fmt.Sprintf("Error sending request: %v", err))
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		s.log.Error(fmt.Sprintf("Unexpected status code: %d", resp.StatusCode))
		return nil, errors.New("failed to get product information")
	}
	defer resp.Body.Close()
	var product models.Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		s.log.Error(fmt.Sprintf("Error decoding JSON: %v", err))
		return nil, err
	}
	return &product, nil
}

func (s *Service) ReservationProduct(order models.Order) error {
	err := s.broker.Publish("reservation_product", order, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) OrderResult(queue string) {
	_ = s.broker.Consume(queue, "", func(message amqp.Delivery, headers map[string]interface{}) {
		var orderBillingResponse models.Order
		if err := json.Unmarshal(message.Body, &orderBillingResponse); err != nil {
			s.log.Error(fmt.Sprintf("Failed to decode JSON: %s", err.Error()))
			return
		}
		err := s.repository.UpdateOrderByID(orderBillingResponse)
		if err != nil {
			s.log.Error(fmt.Sprintf("Failed to update order: %s", err.Error()))
		}
		notification := models.Notification{UserID: orderBillingResponse.UserID, Status: orderBillingResponse.Status,
			OrderID: orderBillingResponse.ID, Message: orderBillingResponse.Message}
		_ = s.broker.Publish("notification", notification, nil)
	})
}
