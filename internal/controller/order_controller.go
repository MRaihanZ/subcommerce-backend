package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/errs"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/MRaihanZ/subcommerce-backend/internal/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetOrdersHandler(c *gin.Context) {
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

	orders, err := model.GetAllOrder(id)
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

	res := entity.Response[[]entity.OrderGetResponse]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   orders,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetOrderSellerHandler(c *gin.Context) {
	stateAction := c.Query("state")
	orderCreate := c.Query("order_create")
	orderID := c.Query("order_id")

	if stateAction != "" && stateAction != "next" && stateAction != "previous" {
		msg := "invalid state"
		res := entity.Response[error]{
			Code:   http.StatusBadRequest,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if stateAction == "next" || stateAction == "previous" {
		if orderCreate == "" || orderID == "" {
			msg := "order_create and order_id are required for next/previous state"
			res := entity.Response[error]{
				Code:   http.StatusBadRequest,
				Status: "error",
				Data:   nil,
				Error:  &msg,
			}
			c.JSON(http.StatusBadRequest, res)
			return
		}
	}

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

	orders, err := model.GetOrderBySellerID(id, stateAction, orderCreate, orderID)
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

	res := entity.Response[[]entity.OrderGetResponseBySeller]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   orders,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func CreateOrderHandler(c *gin.Context) {
	var req []entity.OrderRequest
	// check req bind json
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

	var paymentURL, orderID string
	var err error

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

	for _, v := range req {
		exist, err := model.GetExistingOrderSubscriptionByUserId(id, v.PId, v.PVId)
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

		if *exist {
			msg := "Error: Salah satu produk yang dipilih sudah berlangganan"
			res := entity.Response[error]{
				Code:   http.StatusInternalServerError,
				Status: "error",
				Data:   nil,
				Error:  &msg,
			}
			c.JSON(http.StatusInternalServerError, res)
			return
		}
	}

	orderID = uuid.New().String()
	newTotalPrice := int64(req[0].TotalPrice)

	paymentURL, orderID, err = service.CreatePayment(orderID, newTotalPrice, id, "order", nil)
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

	_, err = model.CreateOrder(id, req, paymentURL, orderID)
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

func DeleteOrderHandler(c *gin.Context) {
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

	_, err := model.DeleteOrder(id)
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

	res := entity.Response[string]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   "success",
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func UpdateStatusOrderHandler(c *gin.Context) {
	orderId := c.Param("order_id")
	statusId := c.Param("status_id")
	numStatusId, err := strconv.Atoi(statusId)
	if err != nil {
		msg := "wrong query value"
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

	status, err := model.UpdateStatusOrder(orderId, numStatusId)
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

	res := entity.Response[*string]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   status,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetCheckoutOrdersHandler(c *gin.Context) {
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

	checkouts, err := model.GetAllCheckoutOrder(id)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrNoCheckoutFound):
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

	res := entity.Response[[]entity.GetCheckoutOrderResponse]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   checkouts,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func CreateCheckoutOrderHandler(c *gin.Context) {
	stateAction := c.Query("state")
	if stateAction == "" {
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

	var req []entity.OrderCheckout
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

	for _, v := range req {
		exist, err := model.GetExistingOrderSubscriptionByUserId(id, v.PId, v.PVId)
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

		if *exist {
			msg := "Error: Salah satu produk yang dipilih sudah berlangganan"
			res := entity.Response[error]{
				Code:   http.StatusInternalServerError,
				Status: "error",
				Data:   nil,
				Error:  &msg,
			}
			c.JSON(http.StatusInternalServerError, res)
			return
		}
	}

	cart, err := model.CreateCheckoutOrder(id, req, stateAction)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrNoCheckoutFound):
			code = http.StatusInternalServerError
			msg = err.Error()
		case errors.Is(err, errs.ErrCheckoutRequestZero):
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

	res := entity.Response[*string]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   cart,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetOrderPaymentsHandler(c *gin.Context) {
	payments, err := model.GetAllOrderPayment()
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrNoPaymentFound):
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

	res := entity.Response[[]entity.GetOrderPaymentResponse]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   payments,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func UpdateOrderHandler(c *gin.Context) {
	// Update order by ID
}
