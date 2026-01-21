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
	log.Println("finish get data")

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

	log.Println("finish fill req")

	body, _ := json.Marshal(reqBody)

	log.Println("start request to midtrans")
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

	log.Println(res)

	paymentURL, ok := res["payment_url"].(string)
	if !ok {
		return "", "", fmt.Errorf("unexpected midtrans response: %v", res)
	}

	rawOrderID := res["order_id"].(string)

	// remove last "-<digits>"
	re := regexp.MustCompile(`-\d+$`)
	newOrderID := re.ReplaceAllString(rawOrderID, "")
	if !ok {
		return "", "", fmt.Errorf("unexpected midtrans response: %v", res)
	}

	return paymentURL, newOrderID, nil
}
