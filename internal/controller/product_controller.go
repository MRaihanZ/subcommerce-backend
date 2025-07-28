package controller

import (
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/gin-gonic/gin"
)

func GetProductsHandler(c *gin.Context) {
	products, err := model.GetAllProductsSummarize()
	if err != nil {
		msg := err.Error()
		res := entity.Response[[]entity.ProductSummarize]{
			Code:   500,
			Status: "error",
			Data:   products,
			Error:  &msg,
		}
		c.JSON(500, res)
		return
	}

	if products == nil {
		msg := "no user found"
		res := entity.Response[[]entity.ProductSummarize]{
			Code:   404,
			Status: "error",
			Data:   products,
			Error:  &msg,
		}
		c.JSON(404, res)
		return
	}

	res := entity.Response[[]entity.ProductSummarize]{
		Code:   200,
		Status: "ok",
		Data:   products,
		Error:  nil,
	}
	c.JSON(200, res)
}

func GetProductHandler(c *gin.Context) {
	id := c.Param("id")
	product, err := model.GetProductById(id)
	if err != nil {
		msg := err.Error()
		res := entity.Response[*entity.JsonProduct]{
			Code:   500,
			Status: "error",
			Data:   product,
			Error:  &msg,
		}
		c.JSON(500, res)
		return
	}

	if product == nil {
		msg := "product not found"
		res := entity.Response[*entity.JsonProduct]{
			Code:   404,
			Status: "error",
			Data:   product,
			Error:  &msg,
		}
		c.JSON(404, res)
		return
	}

	res := entity.Response[*entity.JsonProduct]{
		Code:   200,
		Status: "ok",
		Data:   product,
		Error:  nil,
	}
	c.JSON(200, res)
}
