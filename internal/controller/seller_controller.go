package controller

import (
	"net/http"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetSeller(c *gin.Context) {

}

func GetSellerSummarize(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		msg := "format seller ID salah"
		res := entity.Response[*error]{
			Code:   http.StatusBadRequest,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusBadRequest, res)
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

	res := entity.Response[*entity.SellerSummarize]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   seller,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}
