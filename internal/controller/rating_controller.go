package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/errs"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CreateRatingHandler(c *gin.Context) {
	productId := c.Param("product_id")
	numProdId, err := strconv.Atoi(productId)
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

	productVariantId := c.Param("product_variant_id")
	numProdVarId, err := strconv.Atoi(productVariantId)
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

	orderId := c.Param("order_id")
	numOrderId, err := strconv.Atoi(orderId)
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

	var req *entity.RatingRequest
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

	rating, err := model.CreateRating(numProdId, numProdVarId, numOrderId, id, req)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrNoRatingFound):
			code = http.StatusInternalServerError
			msg = err.Error()
		case errors.Is(err, errs.ErrAlreadyRated):
			code = http.StatusInternalServerError
			msg = err.Error()
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

	res := entity.Response[*int]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   rating,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetRatingHandler(c *gin.Context) {
	productId := c.Param("product_id")
	numProdId, err := strconv.Atoi(productId)
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

	rating, err := model.GetRatingByProdId(numProdId)
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

	if rating == nil {
		msg := "no products found"
		res := entity.Response[error]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[*entity.Rating]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   rating,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetRatingCommentsHandler(c *gin.Context) {
	productId := c.Param("product_id")
	numProdId, err := strconv.Atoi(productId)
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

	comments, err := model.GetRatingCommentsByProdId(numProdId)
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

	if comments == nil {
		msg := "no products found"
		res := entity.Response[error]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[[]entity.RatingComments]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   comments,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}
