package service

import (
	"log"
	"regexp"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/MRaihanZ/subcommerce-backend/internal/utils"
)

func CreatePayment(orderID string, amount int64, userId interface{}, state string, orderProductData *entity.OrderSubscriptionResponse) (string, string, error) {
	paymentURL, orderId, err := utils.CreateQrisPaymentLink(orderID, amount, userId, state, orderProductData)
	if err != nil {
		return "", "", err
	}

	// i save this for later on, even though it's not being used
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

func CancelPayment(orderId string) error {
	return utils.CancelMidtransPayment(orderId)
}

func PayoutToSeller(
	sellerName string,
	bankCode string,
	accountNumber string,
	amount int64,
	description string,
) error {

	return utils.CreateMidtransDisbursement(
		sellerName,
		bankCode,
		accountNumber,
		amount,
		description,
	)
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
			"expire":     2,
			"cancel":     2,
			"deny":       12,
		}
	} else {
		statusMap = map[string]int{
			"settlement": 4,
			"capture":    4,
			"pending":    11,
			"expire":     2,
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

		dataOrderUser, err := model.GetUserIdFromOrders(newOrderID)
		if err != nil {
			log.Println("ERROR IN PAYMENT WEBHOOK, WHEN GET USER, PRODUCT, AND PRODUCT VARIANT ID : ", err)
		}

		nextRemove, err := model.GetNextRemoveReminderSchedule(newOrderID)
		if err != nil {
			log.Println("ERROR IN PAYMENT WEBHOOK, WHEN GET NEXT REMOVE REMINDER SCHEDULES: ", err)
		}

		dataIntervalProduct, err := model.GetIntervalProduct(dataOrderUser.ProductId, dataOrderUser.ProductVariantId)
		if err != nil {
			log.Println("ERROR IN PAYMENT WEBHOOK, WHEN GET INTERVAL PRODUCT : ", err)
		}

		next, warning, remove := utils.CalculateReminderDates(
			*nextRemove,
			dataIntervalProduct.Id,
			dataIntervalProduct.Interval,
		)

		_, err = model.UpdateReminderScheduleInterval(newOrderID, next, warning, remove)
		if err != nil {
			log.Println("ERROR IN PAYMENT WEBHOOK, WHEN UPDATE ID REMINDER SCHEDULE: ", err)
		}
	}

	_, err = model.UpdateStatusOrder(newOrderID, status)

	return err
}

func IsPaymentLinkActive(orderId string) (bool, error) {
	status, err := utils.GetMidtransTransactionStatus(orderId)
	if err != nil {
		return false, err
	}

	return status == "pending", nil
}
