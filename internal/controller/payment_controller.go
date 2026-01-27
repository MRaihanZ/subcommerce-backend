package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/MRaihanZ/subcommerce-backend/internal/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// func CreatePaymentHandler(c *gin.Context) {
// 	session := sessions.Default(c)
// 	userId := session.Get("user_id")
// 	if userId == nil {
// 		msg := "id null"
// 		res := entity.Response[error]{
// 			Code:   http.StatusUnauthorized,
// 			Status: "error",
// 			Data:   nil,
// 			Error:  &msg,
// 		}
// 		c.JSON(http.StatusUnauthorized, res)
// 		return
// 	}

// 	var req entity.CreatePaymentRequest
// 	if err := c.BindJSON(&req); err != nil {
// 		msg := err.Error()
// 		res := entity.Response[error]{
// 			Code:   http.StatusBadRequest,
// 			Status: "error",
// 			Data:   nil,
// 			Error:  &msg,
// 		}
// 		c.JSON(http.StatusBadRequest, res)
// 		return
// 	}

// 	paymentURL, orderId, err := service.CreatePayment(req.OrderID, req.Amount, userId, "order", nil)
// 	if err != nil {
// 		msg := err.Error()
// 		res := entity.Response[error]{
// 			Code:   http.StatusInternalServerError,
// 			Status: "error",
// 			Data:   nil,
// 			Error:  &msg,
// 		}
// 		c.JSON(http.StatusInternalServerError, res)
// 		return
// 	}

// 	res := entity.Response[interface{}]{
// 		Code:   http.StatusOK,
// 		Status: "ok",
// 		Data: map[string]string{
// 			"payment_url": paymentURL,
// 			"order_id":    orderId,
// 		},
// 		Error: nil,
// 	}
// 	c.JSON(http.StatusOK, res)
// }

func MidtransPayoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	sellerId := session.Get("seller_id")
	if sellerId == nil {
		msg := "id null"
		res := entity.Response[error]{
			Code:   http.StatusUnauthorized,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	var req entity.CreatePayoutRequest
	if err := c.BindJSON(&req); err != nil {
		msg := err.Error()
		res := entity.Response[error]{
			Code:   http.StatusBadRequest,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// err := service.PayoutToSeller(req.Name, "gopay", req.PayId, req.Amount, req.Description)
	// if err != nil {
	// 	msg := err.Error()
	// 	res := entity.Response[error]{
	// 		Code:   http.StatusInternalServerError,
	// 		Status: "error",
	// 		Data:   nil,
	// 		Error:  &msg,
	// 	}
	// 	c.JSON(http.StatusInternalServerError, res)
	// 	return
	// }

	payoutId, err := model.CreatePayout(sellerId, req.Name, req.PayId, req.Description, req.Amount)
	if err != nil {
		msg := err.Error()
		res := entity.Response[error]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	wallet, err := model.GetCurrentWalletAmmount(sellerId)
	if err != nil {
		msg := err.Error()
		res := entity.Response[error]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if wallet == nil {
		msg := "seller tidak ditemukan"
		res := entity.Response[error]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	newWallet := *wallet - req.Amount

	_, err = model.UpdateWalletSellerById(sellerId, newWallet)
	if err != nil {
		msg := err.Error()
		res := entity.Response[error]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	sellerData, err := model.GetSellerNameEmailById(sellerId)
	if err != nil {
		msg := err.Error()
		res := entity.Response[error]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if sellerData == nil {
		msg := "seller tidak ditemukan"
		res := entity.Response[error]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	today := time.Now()
	emailData := entity.PayoutEmail{
		SellerName:     sellerData.Name,
		Id:             *payoutId,
		TransferName:   req.Name,
		TransferType:   "gopay",
		TransferAmount: req.Amount,
		CreatedAt:      today.Format("02-01-2006 15:04:05"),
	}

	htmlBody, err := service.RenderPayoutEmail(emailData, "E:/GIU/Devel/go_app/subcommerce-backend/internal/templates/payout_queue.html")
	if err != nil {
		log.Fatal(err)
	}

	sender := service.NewBrevoSender()
	_ = sender.SendMail(
		sellerData.Email,
		"Informasi Status Pencairan Dana Penjual Subcommerce",
		htmlBody,
	)

	res := entity.Response[string]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   "success",
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

// used in midtrans, don't change response
func MidtransWebhookHandler(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	orderID := payload["order_id"].(string)
	transactionStatus := payload["transaction_status"].(string)

	err := service.HandleWebhook(orderID, transactionStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
