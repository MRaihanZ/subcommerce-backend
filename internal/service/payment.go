package service

import (
	"log"
	"regexp"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/MRaihanZ/subcommerce-backend/internal/utils"
)

func CreatePayment(orderID string, amount int64, userId interface{}, state string, orderProductData *entity.GetCheckoutOrderPaymentResponse) (string, string, error) {
	paymentURL, orderId, err := utils.CreateQrisPaymentLink(orderID, amount, userId, state, orderProductData)
	if err != nil {
		return "", "", err
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

	return paymentURL, orderId, nil
}

func HandleWebhook(orderID string, transactionStatus string) error {
	var statusMap map[string]int

	// remove last "-<digits>"
	re := regexp.MustCompile(`-\d+$`)
	newOrderID := re.ReplaceAllString(orderID, "")

	check, err := model.GetExistingOrderSubscription(newOrderID)
	if err != nil {
		log.Println("ERROR IN PAYMENT WEBHOOK, WHEN CHECKING EXISTING ORDER: ", err)
	}

	if *check {
		statusMap = map[string]int{
			"settlement": 13,
			"capture":    13,
			"pending":    11,
			"expire":     3,
			"cancel":     2,
			"deny":       12,
		}
	} else {
		statusMap = map[string]int{
			"settlement": 4,
			"capture":    4,
			"pending":    11,
			"expire":     3,
			"cancel":     2,
			"deny":       12,
		}
	}

	status, ok := statusMap[transactionStatus]
	if !ok {
		status = 12
	} else if status == 13 {
		id, err := model.UpdateReminderScheduleId(newOrderID)
		if err != nil {
			log.Println("ERROR IN PAYMENT WEBHOOK, WHEN UPDATE ID REMINDER SCHEDULE: ", err)
		}
		if id == nil {
			log.Println("ERROR IN PAYMENT WEBHOOK, WHEN UPDATE ID REMINDER SCHEDULE: ", id)
		}
	}

	_, err = model.UpdateStatusOrder(newOrderID, status)

	return err
}
