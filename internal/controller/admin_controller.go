package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/errs"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/MRaihanZ/subcommerce-backend/internal/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetAllAdminsHandler(c *gin.Context) {
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

	user, err := model.GetAllAdmin()
	if err != nil {
		msg := err.Error()
		res := entity.Response[[]entity.Admin]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if user == nil {
		msg := "admin not found"
		res := entity.Response[[]entity.Admin]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[[]entity.Admin]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   user,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetAdminHandler(c *gin.Context) {
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

	user, err := model.GetAdminByName(searchQuery)
	if err != nil {
		msg := err.Error()
		res := entity.Response[[]entity.Admin]{
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
		res := entity.Response[[]entity.Admin]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[[]entity.Admin]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   user,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func CreateAdminHandler(c *gin.Context) {
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

	var req entity.CreateAdmin
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

	user, err := model.CreateAdmin(req)
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

func UpdateAdminHandler(c *gin.Context) {
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

	var payload entity.UpdateAdmin
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

	user, err := model.UpdateAdmin(id, payload.Name, payload.Email, payload.Password)
	if err != nil {
		msg := err.Error()
		res := entity.Response[*entity.Admin]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := entity.Response[*entity.Admin]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   user,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func DeleteAdminHandler(c *gin.Context) {
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

	userName, err := model.DeleteAdmin(id)
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

func GetAllProductsHandler(c *gin.Context) {
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

	products, err := model.GetAllProducts()
	if err != nil {
		msg := err.Error()
		res := entity.Response[[]*entity.JsonProductAdd]{
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
		res := entity.Response[[]*entity.JsonProductAdd]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   products,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[[]*entity.JsonProductAdd]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   products,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetProductsByAdminHandler(c *gin.Context) {
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

	products, err := model.GetProductsByName(searchQuery)
	if err != nil {
		msg := err.Error()
		res := entity.Response[[]*entity.JsonProductAdd]{
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
		res := entity.Response[[]*entity.JsonProductAdd]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   products,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[[]*entity.JsonProductAdd]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   products,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

// func CreateProductByAdminHandler(c *gin.Context) {
// 	session := sessions.Default(c)
// 	idSession := session.Get("admin_id")
// 	if idSession == nil || idSession == "noId" {
// 		msg := "id null"
// 		res := entity.Response[error]{
// 			Code:   http.StatusUnauthorized,
// 			Status: "error",
// 			Data:   nil,
// 			Error:  &msg,
// 		}
// 		c.JSON(http.StatusUnauthorized, res)
// 		return
// 	}

// 	var req entity.ReqProductAdd
// 	productJSON := c.PostForm("product")
// 	if productJSON == "" {
// 		msg := "missing product data"
// 		res := entity.Response[error]{
// 			Code:   http.StatusBadRequest,
// 			Status: "error",
// 			Data:   nil,
// 			Error:  &msg,
// 		}
// 		c.JSON(http.StatusBadRequest, res)
// 		return
// 	}

// 	if err := json.Unmarshal([]byte(productJSON), &req); err != nil {
// 		msg := "invalid product"
// 		res := entity.Response[error]{
// 			Code:   http.StatusBadRequest,
// 			Status: "error",
// 			Data:   nil,
// 			Error:  &msg,
// 		}
// 		c.JSON(http.StatusBadRequest, res)
// 		return
// 	}

// 	product, err := model.CreateProduct(id, req)
// 	if err != nil {
// 		var code int
// 		var msg string
// 		switch {
// 		case errors.Is(err, errs.ErrNoProductFound):
// 			code = http.StatusBadRequest
// 			msg = err.Error()
// 		default:
// 			log.Println("err product")
// 			log.Println(err)
// 			code = http.StatusInternalServerError
// 			msg = "internal server error"
// 		}
// 		res := entity.Response[error]{
// 			Code:   code,
// 			Status: "error",
// 			Data:   nil,
// 			Error:  &msg,
// 		}
// 		c.JSON(code, res)
// 		return
// 	}

// 	_, err = model.CreateProductVariantsAdd(id, *product, req.PVariants)
// 	if err != nil {
// 		var code int
// 		var msg string
// 		switch {
// 		case errors.Is(err, errs.ErrNoProductVariantFound):
// 			code = http.StatusBadRequest
// 			msg = err.Error()
// 		default:
// 			log.Println("err product variant")
// 			log.Println(err)
// 			code = http.StatusInternalServerError
// 			msg = "internal server error"
// 		}
// 		res := entity.Response[error]{
// 			Code:   code,
// 			Status: "error",
// 			Data:   nil,
// 			Error:  &msg,
// 		}
// 		c.JSON(code, res)
// 		return
// 	}

// 	form, err := c.MultipartForm()
// 	if err != nil {
// 		code := http.StatusBadRequest
// 		msg := "gambar harus di pilih"
// 		res := entity.Response[error]{
// 			Code:   code,
// 			Status: "error",
// 			Data:   nil,
// 			Error:  &msg,
// 		}
// 		c.JSON(code, res)
// 		return
// 	}

// 	files := form.File["images"]
// 	uploadService, err := service.UpdateImgProduct(files, id, *product, c)
// 	if err != nil {
// 		log.Println("err uploadService")
// 		log.Println(err)
// 		msg := err.Error()
// 		res := entity.Response[error]{
// 			Code:   http.StatusInternalServerError,
// 			Status: "error",
// 			Data:   nil,
// 			Error:  &msg,
// 		}
// 		c.JSON(http.StatusInternalServerError, res)
// 		return
// 	}

// 	productImages, err := model.CreateProductImages(id, *product, uploadService)
// 	if err != nil {
// 		var code int
// 		var msg string
// 		switch {
// 		case errors.Is(err, errs.ErrNoProductVariantFound):
// 			code = http.StatusBadRequest
// 			msg = err.Error()
// 		default:
// 			log.Println("err product images")
// 			log.Println(err)
// 			code = http.StatusInternalServerError
// 			msg = "internal server error"
// 		}
// 		res := entity.Response[error]{
// 			Code:   code,
// 			Status: "error",
// 			Data:   nil,
// 			Error:  &msg,
// 		}
// 		c.JSON(code, res)
// 		return
// 	}

// 	res := entity.Response[*int]{
// 		Code:   http.StatusOK,
// 		Status: "ok",
// 		Data:   productImages,
// 		Error:  nil,
// 	}
// 	c.JSON(http.StatusOK, res)
// }

func UpdateActiveProductHandler(c *gin.Context) {
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

	activeStr := c.Param("active")
	active, err := strconv.ParseBool(activeStr)
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

	productVariant, err := model.UpdateActiveProduct(numProdId, active)
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
}

func DeleteProductByAdminHandler(c *gin.Context) {
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

func GetPayoutsHandler(c *gin.Context) {
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

	payouts, err := model.GetPayouts(c.Request.Context())
	if err != nil {
		msg := "Payout tidak ditemukan"
		res := entity.Response[error]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := entity.Response[[]entity.Payout]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   payouts,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func UpdatePayoutStatusHandler(c *gin.Context) {
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

	payoutID := c.Param("id")

	var req entity.UpdatePayoutStatusRequest
	if err := c.BindJSON(&req); err != nil {
		msg := "invalid request body"
		res := entity.Response[error]{
			Code:   http.StatusBadRequest,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	err := model.UpdatePayoutStatus(
		c.Request.Context(),
		payoutID,
		req.Status,
	)
	if err != nil {
		msg := "gagal untuk mengubah payout status"
		res := entity.Response[error]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	payoutEmailData, err := model.GetPayoutEmailData(payoutID)
	if err != nil {
		msg := "gagal untuk mengambil data untuk email"
		res := entity.Response[error]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if payoutEmailData == nil {
		msg := "tidak ada payout data"
		res := entity.Response[error]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	// kirim email ke seller
	today := time.Now()
	emailData := entity.PayoutEmail{
		SellerName:     payoutEmailData.SellerName,
		Id:             payoutID,
		TransferName:   payoutEmailData.Name,
		TransferType:   "gopay",
		TransferAmount: payoutEmailData.Amount,
		CreatedAt:      today.Format("02-01-2006 15:04:05"),
	}

	htmlBody, err := service.RenderPayoutEmail(emailData, "E:/GIU/Devel/go_app/subcommerce-backend/internal/templates/payout_process.html")
	if err != nil {
		log.Fatal(err)
	}

	sender := service.NewBrevoSender()
	_ = sender.SendMail(
		payoutEmailData.Email,
		"Informasi Status Pencairan Dana Penjual Subcommerce",
		htmlBody,
	)

	res := entity.Response[string]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   "payout status diperbarui",
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}
