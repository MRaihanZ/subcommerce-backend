package service

import (
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/MRaihanZ/subcommerce-backend/internal/utils"
)

func CreatePayment(orderID string, amount int64, userId interface{}) (string, error) {
	paymentURL, err := utils.CreateQrisPaymentLink(orderID, amount, userId)
	if err != nil {
		return "", err
	}

	// payment := &model.Payment{
	// 	ID:         uuid.New().String(),
	// 	OrderID:    orderID,
	// 	Amount:     amount,
	// 	Status:     "PENDING",
	// 	PaymentURL: paymentURL,
	// 	CreatedAt:  time.Now(),
	// 	ExpiredAt:  time.Now().Add(30 * time.Minute),
	// }

	// err = model.CreatePayment(payment)
	// if err != nil {
	// 	return "", err
	// }

	return paymentURL, nil
}

func HandleWebhook(orderID string, transactionStatus string) error {
	statusMap := map[string]int{
		"settlement": 4,
		"capture":    4,
		"pending":    11,
		"expire":     3,
		"cancel":     2,
		"deny":       12,
	}

	status, ok := statusMap[transactionStatus]
	if !ok {
		status = 12
	}

	_, err := model.UpdateStatusOrder(orderID, status)

	return err
}
