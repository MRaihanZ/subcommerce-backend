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

	var queryId interface{}
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
		queryId = userId
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
		queryId = sellerId
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

	_, err := model.UpdateStatusOrder(paramOrderId, 14)
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

}
