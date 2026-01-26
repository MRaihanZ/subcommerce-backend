package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
)

func CreateQrisPaymentLink(orderID string, amount int64, userId interface{}, state string, orderProductData *entity.GetCheckoutOrderPaymentResponse) (string, string, error) {
	userInfo, _ := model.GetUserPaymentById(userId)
	data, _ := model.GetAllCheckoutOrderPayment(userId)

	reqBody := entity.MidtransChargeRequest{}
	reqBody.EnablePayments = append(reqBody.EnablePayments, "other_qris")
	reqBody.TransactionDetails.OrderID = orderID
	reqBody.TransactionDetails.GrossAmt = amount
	reqBody.Expiry.Unit = "minutes"
	reqBody.Expiry.Duration = 30
	reqBody.CustomerRequired = true
	reqBody.CustomerDetails.FirstName = userInfo.Name
	reqBody.CustomerDetails.Email = userInfo.Email
	reqBody.ItemDetails = make([]struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
		Price    int    `json:"price"`
	}, 0, (len(data) + 1))
	for _, info := range data {

		reqBody.ItemDetails = append(reqBody.ItemDetails, struct {
			Id       string `json:"id"`
			Name     string `json:"name"`
			Quantity int    `json:"quantity"`
			Price    int    `json:"price"`
		}{
			Id:       strconv.Itoa(info.PId),
			Name:     info.PName,
			Quantity: info.Quantity,
			Price:    info.UnitPrice,
		})
	}
	reqBody.ItemDetails = append(reqBody.ItemDetails, struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
		Price    int    `json:"price"`
	}{
		Id:       "1",
		Name:     "Fee Aplication",
		Quantity: 1,
		Price:    3000,
	})
	reqBody.QrisDetail.Acquirer = "gopay"
	reqBody.Callbacks.Finish = os.Getenv("WEBSITE_URL") + "/order-list"

	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest(
		"POST",
		os.Getenv("MIDTRANS_URL"),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", "", err
	}

	req.SetBasicAuth(os.Getenv("MIDTRANS_SERVER_KEY_SANDBOX"), "")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	paymentURL, ok := res["payment_url"].(string)
	if !ok {
		return "", "", fmt.Errorf("unexpected midtrans response: %v", res)
	}
	log.Println(res)

	rawOrderID := res["order_id"].(string)

	// remove last "-<digits>"
	re := regexp.MustCompile(`-\d+$`)
	newOrderID := re.ReplaceAllString(rawOrderID, "")
	if !ok {
		return "", "", fmt.Errorf("unexpected midtrans response: %v", res)
	}

	return paymentURL, newOrderID, nil
}

func CancelMidtransPayment(orderId string) error {
	url := fmt.Sprintf("%s/v2/%s/cancel", os.Getenv("MIDTRANS_API_BASE"), orderId)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(os.Getenv("MIDTRANS_SERVER_KEY_SANDBOX"), "")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&body)
		return fmt.Errorf("midtrans cancel failed: %v", body)
	}

	return nil
}

func CreateMidtransDisbursement(
	name string,
	bank string,
	account string,
	amount int64,
	description string,
) error {

	url := fmt.Sprintf("%s/v1/disbursements", os.Getenv("MIDTRANS_DISBURSEMENT_BASE"))

	reqBody := map[string]interface{}{
		"name":        name,
		"bank":        bank,
		"account":     account,
		"amount":      amount,
		"description": description,
	}

	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.SetBasicAuth(os.Getenv("MIDTRANS_SERVER_KEY_SANDBOX"), "")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		return fmt.Errorf("midtrans payout failed: %v", res)
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	log.Println(res)

	return nil
}

func GetMidtransTransactionStatus(orderId string) (string, error) {
	url := fmt.Sprintf("%s/v2/%s/status", os.Getenv("MIDTRANS_API_BASE"), orderId)

	req, err := http.NewRequest("GET", url, nil)
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

	if resp.StatusCode != http.StatusOK {
		var body map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&body)
		return "", fmt.Errorf("midtrans status failed: %v", body)
	}

	var res struct {
		TransactionStatus string `json:"transaction_status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.TransactionStatus, nil
}
