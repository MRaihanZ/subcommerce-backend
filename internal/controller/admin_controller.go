package controller

import (
	"errors"
	"net/http"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/errs"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/MRaihanZ/subcommerce-backend/internal/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetAllUsersHandler(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("admin_id")
	if id == nil || id == "noId" {
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

	user, err := model.GetAllUsers()
	if err != nil {
		msg := err.Error()
		res := entity.Response[[]entity.Users]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if user == nil {
		msg := "user not found"
		res := entity.Response[[]entity.Users]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[[]entity.Users]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   user,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetUserByAdminHandler(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("admin_id")
	if id == nil || id == "noId" {
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

	searchQuery := c.Param("name")
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

	user, err := model.GetUsersByName(searchQuery)
	if err != nil {
		msg := err.Error()
		res := entity.Response[[]entity.Users]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if user == nil {
		msg := "user not found"
		res := entity.Response[[]entity.Users]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[[]entity.Users]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   user,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func CreateUserByAdminHandler(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("admin_id")
	if id == nil || id == "noId" {
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

	var req entity.SignUp
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
	user, err := model.CreateUser(req.Name, req.Email, req.Dob, req.Password)
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

	if user == nil {
		msg := "email sudah terdaftar"
		res := entity.Response[error]{
			Code:   http.StatusConflict,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusConflict, res)
		return
	}

	res := entity.Response[*string]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   user,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func UpdateUserByAdminHandler(c *gin.Context) {
	session := sessions.Default(c)
	sessionId := session.Get("admin_id")
	if sessionId == nil || sessionId == "noId" {
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

	id := c.Param("id")
	if id == "" {
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

	var payload entity.UpdateUser
	if err := c.ShouldBind(&payload); err != nil {
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

	file, err := c.FormFile("img")
	if err != nil {
		user, err := model.UpdateUser(id, payload.Name, payload.ImgPath, payload.Email, payload.Password, payload.Dob)
		if err != nil {
			msg := err.Error()
			res := entity.Response[*entity.User]{
				Code:   http.StatusInternalServerError,
				Status: "error",
				Data:   user,
				Error:  &msg,
			}
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		res := entity.Response[*entity.User]{
			Code:   http.StatusOK,
			Status: "ok",
			Data:   user,
			Error:  nil,
		}
		c.JSON(http.StatusOK, res)
	} else {
		uploadService, err := service.UpdateProfile(file, id, "user", c)
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

		if uploadService == nil {
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

		user, err := model.UpdateUser(id, payload.Name, *uploadService, payload.Email, payload.Password, payload.Dob)
		if err != nil {
			msg := err.Error()
			res := entity.Response[*entity.User]{
				Code:   http.StatusInternalServerError,
				Status: "error",
				Data:   user,
				Error:  &msg,
			}
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		res := entity.Response[*entity.User]{
			Code:   http.StatusOK,
			Status: "ok",
			Data:   user,
			Error:  nil,
		}
		c.JSON(http.StatusOK, res)
	}
}

func DeleteUserByAdminHandler(c *gin.Context) {
	session := sessions.Default(c)
	idSession := session.Get("admin_id")
	if idSession == nil || idSession == "noId" {
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

	id := c.Param("id")
	if id == "" {
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

	deleteImg, err := service.DeleteProfile(id, "user")
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

	if deleteImg == nil {
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

	userName, err := model.DeleteUser(id)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
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

	res := entity.Response[*string]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   userName,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetAllSellersHandler(c *gin.Context) {
	session := sessions.Default(c)
	idSession := session.Get("admin_id")
	if idSession == nil || idSession == "noId" {
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

	seller, err := model.GetAllSellers()
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

	res := entity.Response[[]entity.Seller]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   seller,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetSellerByAdminHandler(c *gin.Context) {
	session := sessions.Default(c)
	idSession := session.Get("admin_id")
	if idSession == nil || idSession == "noId" {
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

	searchQuery := c.Param("name")
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

	seller, err := model.GetSellersByName(searchQuery)
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

	res := entity.Response[[]entity.Seller]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   seller,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func CreateSellerByAdminHandler(c *gin.Context) {
	session := sessions.Default(c)
	idSession := session.Get("admin_id")
	if idSession == nil || idSession == "noId" {
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

	id := c.Param("uid")
	if id == "" {
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

	var req entity.CreateSellerRequest
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

	// Create seller in database
	seller, err := model.CreateSeller(id, req)
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

	if seller == nil {
		msg := "error return seller id"
		res := entity.Response[error]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := entity.Response[*string]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   seller,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func UpdateSellerByAdminHandler(c *gin.Context) {
	session := sessions.Default(c)
	idSession := session.Get("admin_id")
	if idSession == nil || idSession == "noId" {
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

	id := c.Param("id")
	if id == "" {
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

	var payload entity.UpdateSeller
	if err := c.ShouldBind(&payload); err != nil {
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

	file, err := c.FormFile("img")
	if err != nil {
		seller, err := model.UpdateSeller(id, payload.Name, payload.ImgPath, payload.Address)
		if err != nil {
			msg := err.Error()
			res := entity.Response[*entity.Seller]{
				Code:   http.StatusInternalServerError,
				Status: "error",
				Data:   seller,
				Error:  &msg,
			}
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		res := entity.Response[*entity.Seller]{
			Code:   http.StatusOK,
			Status: "ok",
			Data:   seller,
			Error:  nil,
		}
		c.JSON(http.StatusOK, res)
	} else {
		uploadService, err := service.UpdateProfile(file, id, "seller", c)
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

		if uploadService == nil {
			msg := "seller tidak ditemukan"
			res := entity.Response[error]{
				Code:   http.StatusInternalServerError,
				Status: "error",
				Data:   nil,
				Error:  &msg,
			}
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		seller, err := model.UpdateSeller(id, payload.Name, *uploadService, payload.Address)
		if err != nil {
			msg := err.Error()
			res := entity.Response[*entity.Seller]{
				Code:   http.StatusInternalServerError,
				Status: "error",
				Data:   seller,
				Error:  &msg,
			}
			c.JSON(http.StatusInternalServerError, res)
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
}

func DeleteSellerByAdminHandler(c *gin.Context) {
	session := sessions.Default(c)
	idSession := session.Get("admin_id")
	if idSession == nil || idSession == "noId" {
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

	id := c.Param("id")
	if id == "" {
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

	deleteImg, err := service.DeleteProfile(id, "seller")
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

	if deleteImg == nil {
		msg := "seller tidak ditemukan"
		res := entity.Response[error]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	sellerName, err := model.DeleteSeller(id)
	if err != nil {
		var code int
		var msg string
		switch {
		case errors.Is(err, errs.ErrSellerNotFound):
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

	session.Set("seller_id", "noId")
	if err := session.Save(); err != nil {
		msg := "Failed to save session | " + err.Error()
		res := entity.Response[*entity.SignIn]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	res := entity.Response[*string]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   sellerName,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}
