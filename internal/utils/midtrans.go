package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
)

func CreateQrisPaymentLink(orderID string, amount int64, userId interface{}) (string, error) {
	userInfo, _ := model.GetUserPaymentById(userId)

	reqBody := entity.MidtransChargeRequest{
		PaymentType: "qris",
	}
	reqBody.TransactionDetails.OrderID = orderID
	reqBody.TransactionDetails.GrossAmt = amount
	reqBody.Expiry.Unit = "minutes"
	reqBody.Expiry.Duration = 30
	reqBody.CustomerDetails.FirstName = userInfo.Name
	reqBody.CustomerDetails.Email = userInfo.Email

	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest(
		"POST",
		os.Getenv("MIDTRANS_URL"),
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

	paymentURL, ok := res["payment_url"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected midtrans response: %v", res)
	}

	return paymentURL, nil
}
