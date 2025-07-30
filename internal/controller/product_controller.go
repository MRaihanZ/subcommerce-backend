package controller

import (
	"net/http"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/gin-gonic/gin"
)

func GetProductsHandler(c *gin.Context) {
	products, err := model.GetAllProductsSummarize()
	if err != nil {
		msg := err.Error()
		res := entity.Response[[]entity.ProductSummarize]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   products,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if products == nil {
		msg := "no products found"
		res := entity.Response[[]entity.ProductSummarize]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   products,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[[]entity.ProductSummarize]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   products,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetProductsHotHandler(c *gin.Context) {
	products, err := model.GetAllProductsHotSummarize()
	if err != nil {
		msg := err.Error()
		res := entity.Response[[]entity.ProductSummarize]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   products,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if products == nil {
		msg := "no products found"
		res := entity.Response[[]entity.ProductSummarize]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   products,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[[]entity.ProductSummarize]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   products,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetProductsDiscountHandler(c *gin.Context) {
	products, err := model.GetAllProductsDiscountSummarize()
	if err != nil {
		msg := err.Error()
		res := entity.Response[[]entity.ProductSummarize]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   products,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if products == nil {
		msg := "no products found"
		res := entity.Response[[]entity.ProductSummarize]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   products,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[[]entity.ProductSummarize]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   products,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetProductHandler(c *gin.Context) {
	id := c.Param("id")
	product, err := model.GetProductById(id)
	if err != nil {
		msg := err.Error()
		res := entity.Response[*entity.JsonProduct]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   product,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if product == nil {
		msg := "product not found"
		res := entity.Response[*entity.JsonProduct]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   product,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[*entity.JsonProduct]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   product,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}
