package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/errs"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/MRaihanZ/subcommerce-backend/internal/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetProductsSellerHandler(c *gin.Context) {
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

	product, err := model.GetProductsBySellerId(id)
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

	if product == nil {
		msg := "product not found"
		res := entity.Response[error]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[[]*entity.JsonProductAdd]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   product,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetProductsHandler(c *gin.Context) {
	var err error
	//search wajib
	//min max optional
	searchQuery := c.Query("search")
	if searchQuery == "" {
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

	minPriceQuery := c.Query("min")
	var numMinPriceQuery int
	if minPriceQuery != "" {
		numMinPriceQuery, err = strconv.Atoi(minPriceQuery)
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
	}

	maxPriceQuery := c.Query("max")
	var numMaxPriceQuery int
	if maxPriceQuery != "" {
		numMaxPriceQuery, err = strconv.Atoi(maxPriceQuery)
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
	}

	products, err := model.GetAllProductsSummarize(searchQuery, numMinPriceQuery, numMaxPriceQuery)
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
	var err error
	minPriceQuery := c.Query("min")
	var numMinPriceQuery int
	if minPriceQuery != "" {
		numMinPriceQuery, err = strconv.Atoi(minPriceQuery)
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
	}

	maxPriceQuery := c.Query("max")
	var numMaxPriceQuery int
	if maxPriceQuery != "" {
		numMaxPriceQuery, err = strconv.Atoi(maxPriceQuery)
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
	}

	products, err := model.GetAllProductsHotSummarize(numMinPriceQuery, numMaxPriceQuery)
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
	var err error
	minPriceQuery := c.Query("min")
	var numMinPriceQuery int
	if minPriceQuery != "" {
		numMinPriceQuery, err = strconv.Atoi(minPriceQuery)
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
	}

	maxPriceQuery := c.Query("max")
	var numMaxPriceQuery int
	if maxPriceQuery != "" {
		numMaxPriceQuery, err = strconv.Atoi(maxPriceQuery)
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
	}

	products, err := model.GetAllProductsDiscountSummarize(numMinPriceQuery, numMaxPriceQuery)
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

func CreateProductHandler(c *gin.Context) {
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

	var req entity.ReqProductAdd
	productJSON := c.PostForm("product")
	if productJSON == "" {
		msg := "missing product data"
		res := entity.Response[error]{
			Code:   http.StatusBadRequest,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if err := json.Unmarshal([]byte(productJSON), &req); err != nil {
		msg := "invalid product"
		res := entity.Response[error]{
			Code:   http.StatusBadRequest,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	product, err := model.CreateProduct(id, req)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrNoProductFound):
			code = http.StatusBadRequest
			msg = err.Error()
		default:
			log.Println("err product")
			log.Println(err)
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

	_, err = model.CreateProductVariantsAdd(id, *product, req.PVariants)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrNoProductVariantFound):
			code = http.StatusBadRequest
			msg = err.Error()
		default:
			log.Println("err product variant")
			log.Println(err)
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

	form, err := c.MultipartForm()
	if err != nil {
		code := http.StatusBadRequest
		msg := "gambar harus di pilih"
		res := entity.Response[error]{
			Code:   code,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(code, res)
		return
	}

	files := form.File["images"]
	uploadService, err := service.UpdateImgProduct(files, id, *product, c)
	if err != nil {
		log.Println("err uploadService")
		log.Println(err)
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

	productImages, err := model.CreateProductImages(id, *product, uploadService)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrNoProductVariantFound):
			code = http.StatusBadRequest
			msg = err.Error()
		default:
			log.Println("err product images")
			log.Println(err)
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
		Data:   productImages,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func CreateProductVariantsHandler(c *gin.Context) {
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

	var req []entity.ReqProductVariantUpdate
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

	productVariant, err := model.CreateProductVariants(id, numProdId, numProdVarId, req)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrNoProductVariantFound):
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
		Data:   productVariant,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func UpdateProductHandler(c *gin.Context) {
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

	stateAction := c.Query("state")
	switch stateAction {
	case "product":
		var req entity.ReqProductUpdate
		productJSON := c.PostForm("product")
		if productJSON == "" {
			msg := "missing product data"
			res := entity.Response[error]{
				Code:   http.StatusBadRequest,
				Status: "error",
				Data:   nil,
				Error:  &msg,
			}
			c.JSON(http.StatusBadRequest, res)
			return
		}

		if err := json.Unmarshal([]byte(productJSON), &req); err != nil {
			msg := "invalid product"
			res := entity.Response[error]{
				Code:   http.StatusBadRequest,
				Status: "error",
				Data:   nil,
				Error:  &msg,
			}
			c.JSON(http.StatusBadRequest, res)
			return
		}

		form, err := c.MultipartForm()
		if err == nil {
			files := form.File["images"]
			if len(files) != 0 {
				result, err := service.DeleteImgProduct(id, numProdId)
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

				if result == nil {
					msg := "user tidak ditemukan"
					res := entity.Response[error]{
						Code:   http.StatusInternalServerError,
						Status: "error",
						Data:   nil,
						Error:  &msg,
					}
					c.JSON(http.StatusInternalServerError, res)
					return
				}

				uploadService, err := service.UpdateImgProduct(files, id, numProdId, c)
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

				_, err = model.DeleteProductImage(numProdId)
				if err != nil {
					var code int
					var msg string
					switch {
					case errors.Is(err, errs.ErrNoProductImagesFound):
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

				_, err = model.UpdateProductImage(numProdId, uploadService)
				if err != nil {
					var code int
					var msg string
					switch {
					case errors.Is(err, errs.ErrNoProductImagesFound):
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
			}
		}

		if !req.HasVariant {
			productVariant, err := model.UpdateProductDefault(id, numProdId, numProdVarId, req)
			if err != nil {
				var code int
				var msg string
				switch {
				case errors.Is(err, errs.ErrNoProductFound):
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
				Data:   productVariant,
				Error:  nil,
			}
			c.JSON(http.StatusOK, res)
			return
		} else {
			productVariant, err := model.UpdateProduct(id, numProdId, req)
			if err != nil {
				var code int
				var msg string
				switch {
				case errors.Is(err, errs.ErrNoProductFound):
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
				Data:   productVariant,
				Error:  nil,
			}
			c.JSON(http.StatusOK, res)
			return
		}
	case "variant":
		var req entity.ReqProductVariantUpdate
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

		productVariant, err := model.UpdateProductVariant(id, numProdId, numProdVarId, req)
		if err != nil {
			var code int
			var msg string
			switch {
			case errors.Is(err, errs.ErrNoProductFound):
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
			Data:   productVariant,
			Error:  nil,
		}
		c.JSON(http.StatusOK, res)
		return
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
}

func DeleteProductHandler(c *gin.Context) {
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

	stateAction := c.Query("state")
	switch stateAction {
	case "product":
		product, err := model.DeleteProduct(numProdId)
		if err != nil {
			var code int
			var msg string
			switch {
			case errors.Is(err, errs.ErrNoProductFound):
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
			Data:   product,
			Error:  nil,
		}
		c.JSON(http.StatusOK, res)
		return
	case "variant":
		product, err := model.DeleteProductVariant(numProdId, numProdVarId)
		if err != nil {
			var code int
			var msg string
			switch {
			case errors.Is(err, errs.ErrNoProductVariantFound):
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
			Data:   product,
			Error:  nil,
		}
		c.JSON(http.StatusOK, res)
		return
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
}
