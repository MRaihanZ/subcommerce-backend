package controller

import (
	"errors"
	"net/http"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/errs"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetOrdersHandler(c *gin.Context) {
	// Get all orders
}

func CreateOrderHandler(c *gin.Context) {
	var req []entity.OrderRequest
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

	cart, err := model.CreateOrder(id, req)
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

	res := entity.Response[*string]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   cart,
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
