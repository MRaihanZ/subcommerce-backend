package controller

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/errs"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/MRaihanZ/subcommerce-backend/internal/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAllSubscriptionsByUserHandler(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("user_id")
	if id == nil {
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

	subs, err := model.GetAllSubscriptionsByUser(id)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrNoSubscriptionFound):
			code = http.StatusNotFound
			msg = err.Error()
		default:
			code = http.StatusInternalServerError
			msg = "internal server error"
		}
		res := entity.Response[error]{
			Code:   code,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(code, res)
		return
	}

	res := entity.Response[[]entity.UserSubscription]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   subs,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func CreateOrderSubscriptionHandler(c *gin.Context) {
	var req entity.OrderSubscriptionResponse
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

	session := sessions.Default(c)
	id := session.Get("user_id")
	if id == nil {
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

	orderID := uuid.New().String()
	newTotalPrice := int64(req.OrderRequest.TotalPrice)

	paymentURL, orderID, err := service.CreatePayment(orderID, newTotalPrice, id, "subscription", &req.OrderPaymentRequest)
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

	var reqAfter []entity.OrderRequest

	reqAfter = append(reqAfter, req.OrderRequest)

	_, err = model.CreateOrder(id, reqAfter, paymentURL, orderID)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrNoOrderFound):
			code = http.StatusInternalServerError
			msg = err.Error()
		case errors.Is(err, errs.ErrNoCheckoutFound):
			code = http.StatusBadRequest
			msg = err.Error()
		default:
			code = http.StatusInternalServerError
			msg = "internal server error"
		}

		res := entity.Response[error]{
			Code:   code,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(code, res)
		return
	}

	res := entity.Response[string]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   paymentURL,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func DeleteSubscription(c *gin.Context) {
	// get session
	session := sessions.Default(c)
	userId := session.Get("user_id")
	sellerId := session.Get("seller_id")

	// cek session if not login
	if userId == nil || sellerId == nil {
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

	stateQuery := c.Query("state")
	switch stateQuery {
	case "user":
		if userId == "noId" {
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
	case "seller":
		if sellerId == "noId" {
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
	default:
		msg := "wrong query"
		res := entity.Response[error]{
			Code:   http.StatusBadRequest,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	paramOrderId := c.Param("order_id")
	if paramOrderId == "" {
		msg := "wrong query"
		res := entity.Response[error]{
			Code:   http.StatusBadRequest,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	var req entity.CancelationSubscriptionByUserEmailData
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

	seller, err := model.GetEmailSellerByOrderId(paramOrderId)
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

	if seller == nil {
		msg := "seller not found"
		res := entity.Response[error]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	_, err = model.UpdateStatusOrder(paramOrderId, 14)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrNoOrderFound):
			code = http.StatusInternalServerError
			msg = err.Error()
		default:
			code = http.StatusInternalServerError
			msg = "internal server error"
		}
		res := entity.Response[error]{
			Code:   code,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(code, res)
		return
	}

	_, err = model.DeleteReminderSchedule(paramOrderId)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrNoSubscriptionFound):
			code = http.StatusInternalServerError
			msg = err.Error()
		default:
			code = http.StatusInternalServerError
			msg = "internal server error"
		}
		res := entity.Response[error]{
			Code:   code,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(code, res)
		return
	}

	req.Domain = os.Getenv("WEBSITE_URL")
	htmlBody, err := service.RenderSubscriptionCancelationByUserEmail(req, "E:/GIU/Devel/go_app/subcommerce-backend/internal/templates/seller_subscription_cancellation.html")
	if err != nil {
		log.Fatal(err)
	}

	//
	sender := service.NewBrevoSender()
	_ = sender.SendMail(
		*seller,
		"Informasi Pemberhentian Langganan User Subcommerce",
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
