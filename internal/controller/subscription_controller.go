package controller

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

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

func GetAllSubscriptionsBySellerHandler(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("seller_id")
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

	subs, err := model.GetAllSubscriptionsBySeller(id)
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

	res := entity.Response[[]entity.UserSubscriptionBySeller]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   subs,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func CreateOrderSubscriptionHandler(c *gin.Context) {
	var req []entity.OrderSubscriptionResponse
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

	subs, err := model.GetReminderScheduleId(id)
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

	active, err := service.IsPaymentLinkActive(subs.Id)
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

	if active {
		msg := subs.PayLInk
		res := entity.Response[string]{
			Code:   http.StatusOK,
			Status: "ok",
			Data:   msg,
			Error:  nil,
		}
		c.JSON(http.StatusOK, res)
		return
	}

	orderID := uuid.New().String()
	newTotalPrice := int64(req[0].TotalPrice)

	paymentURL, newOrderID, err := service.CreatePayment(orderID, newTotalPrice, id, "subscription", &req[0])
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

	reqAfter = append(reqAfter, req[0].OrderRequest)

	_, err = model.CreateOrder(id, reqAfter, paymentURL, newOrderID, "subscription")
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

	var req entity.CancelSubsRequest
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

	dataUser, err := model.GetUserById(id)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrNoOrderFound):
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

	today := time.Now().Format("02-01-2006 15:04:05")
	emailData := entity.CancelationSubscriptionByUserEmailData{
		SellerName:         req.SellerName,
		UserName:           dataUser.Name,
		ProductName:        req.ProductName,
		ProductVariantName: req.ProductVariantName,
		Timestamp:          today,
		Domain:             os.Getenv("WEBSITE_URL"),
	}

	htmlBody, err := service.RenderSubscriptionCancelationByUserEmail(emailData, "E:/GIU/Devel/go_app/subcommerce-backend/internal/templates/seller_subscription_cancellation.html")
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
