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

func GetCartsHandler(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("user_id")
	if id == nil {
		msg := "id null"
		res := entity.Response[*entity.AddProductCart]{
			Code:   http.StatusUnauthorized,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	carts, err := model.GetAllCartProduct(id)
	if err != nil {
		msg := err.Error()
		res := entity.Response[[]entity.CartProduct]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   carts,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if carts == nil {
		msg := "tidak ada produk di keranjang"
		res := entity.Response[[]entity.CartProduct]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   carts,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[[]entity.CartProduct]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   carts,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func CreateCartHandler(c *gin.Context) {
	var req entity.AddProductCart
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

	cart, err := model.CreateCart(id, req.PId, req.PvId, req.Quantity)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrQuantityLessThanZero):
			code = http.StatusBadRequest
			msg = err.Error()
		case errors.Is(err, errs.ErrProductNotFound):
			code = http.StatusBadRequest
			msg = err.Error()
		case errors.Is(err, errs.ErrQuantityLessThanMinOrder):
			code = http.StatusBadRequest
			msg = err.Error()
		case errors.Is(err, errs.ErrNotEnoughStock):
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

	res := entity.Response[*int]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   cart,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func UpdateCartHandler(c *gin.Context) {
	var req entity.UpdateProductCart
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

	carts, err := model.UpdateCart(id, req.PId, req.PvId, req.Quantity)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrQuantityLessThanZero):
			code = http.StatusBadRequest
			msg = err.Error()
		case errors.Is(err, errs.ErrProductNotFound):
			code = http.StatusBadRequest
			msg = err.Error()
		case errors.Is(err, errs.ErrQuantityLessThanMinOrder):
			code = http.StatusBadRequest
			msg = err.Error()
		case errors.Is(err, errs.ErrNotEnoughStock):
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

	res := entity.Response[*int]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   carts,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func DeleteCartsHandler(c *gin.Context) {
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
	carts, err := model.DeleteCarts(id)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrCartsEmpty):
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

	res := entity.Response[*int64]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   carts,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func DeleteCartHandler(c *gin.Context) {
	prodId := c.Query("product_id")
	prodVarId := c.Query("product_variant_id")
	if prodId == "" || prodVarId == "" {
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

	numProdId, err := strconv.Atoi(prodId)
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

	numProdVarId, err := strconv.Atoi(prodVarId)
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

	cart, err := model.DeleteCart(id, numProdId, numProdVarId)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrNoCartProduct):
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

	res := entity.Response[*entity.DeleteProductCart]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   cart,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}
