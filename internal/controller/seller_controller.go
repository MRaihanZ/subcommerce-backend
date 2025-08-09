package controller

import (
	"net/http"
	"strconv"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetSeller(c *gin.Context) {
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

	seller, err := model.GetSellerById(id)
	if err != nil {
		msg := err.Error()
		res := entity.Response[*error]{
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
		res := entity.Response[*error]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[*entity.Seller]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   seller,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetSellerSummarize(c *gin.Context) {
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

	seller, err := model.GetSellerByIdSummarize(numProdId)
	if err != nil {
		msg := err.Error()
		res := entity.Response[*error]{
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
		res := entity.Response[*error]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[*entity.SellerSummarize]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   seller,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}
