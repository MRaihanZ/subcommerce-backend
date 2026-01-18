package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
)

func CreateQrisPaymentLink(orderID string, amount int64) (string, error) {
	reqBody := entity.MidtransChargeRequest{
		PaymentType: "qris",
	}
	reqBody.TransactionDetails.OrderID = orderID
	reqBody.TransactionDetails.GrossAmt = amount
	reqBody.Expiry.Unit = "minute"
	reqBody.Expiry.Duration = 30

	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest(
		"POST",
		"https://api.midtrans.com/v2/payment-links",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(os.Getenv("MIDTRANS_SERVER_KEY_SANDBOX"), "")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	actions, ok := res["actions"].([]interface{})
	if !ok || len(actions) == 0 {
		return "", fmt.Errorf("midtrans response has no actions: %+v", res)
	}

	action, ok := actions[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid actions format")
	}

	paymentURL, ok := action["url"].(string)
	if !ok {
		return "", fmt.Errorf("payment url not found")
	}

	return paymentURL, nil
}
