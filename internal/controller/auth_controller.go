package controller

import (
	"net/http"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserHandler(c *gin.Context) {
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

	csrfToken := c.GetHeader("X-CSRF-TOKEN")

	session := sessions.Default(c)
	session.Set("user_id", user)
	session.Set("seller_id", "noId")
	session.Set("csrf_token", csrfToken)
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
		Data:   user,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func VerifyUserHandler(c *gin.Context) {
	var req entity.SignIn
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

	csrfToken := c.GetHeader("X-CSRF-TOKEN")

	// Session setup
	session := sessions.Default(c)
	session.Set("admin_id", "noId")
	session.Set("seller_id", "noId")

	// Get admin from database
	admin, err := model.GetAdminByEmail(req.Email)
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

	if admin != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
			msg := err.Error()
			res := entity.Response[error]{
				Code:   http.StatusUnauthorized,
				Status: "error",
				Data:   nil,
				Error:  &msg,
			}
			c.JSON(http.StatusUnauthorized, res)
			return
		}
		session.Set("user_id", admin.Id)
		session.Set("admin_id", admin.Id)
		session.Set("csrf_token", csrfToken)

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

		res := entity.Response[*entity.SignIn]{
			Code:   http.StatusOK,
			Status: "ok",
			Data:   admin,
			Error:  nil,
		}
		c.JSON(http.StatusOK, res)
		return
	}

	// Get user from database
	user, err := model.GetUserByEmail(req.Email)
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
		msg := "user not found"
		res := entity.Response[error]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		msg := err.Error()
		res := entity.Response[error]{
			Code:   http.StatusUnauthorized,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	session.Set("user_id", user.Id)

	seller, err := model.GetSellerByEmail(user.Id)
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

	if seller != nil {
		session.Set("seller_id", seller)
	}

	session.Set("csrf_token", csrfToken)

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

	res := entity.Response[*entity.SignIn]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   user,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func CreateSellerHandler(c *gin.Context) {
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

	session.Set("seller_id", seller)
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
		Data:   seller,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func CheckStatusHandler(c *gin.Context) {
	session := sessions.Default(c)
	var stat entity.Status
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
	stat.IsLogin = true

	seller := session.Get("seller_id")
	if seller == nil {
		msg := "seller null"
		res := entity.Response[error]{
			Code:   http.StatusUnauthorized,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	sellerStr, ok := seller.(string)
	if !ok {
		msg := "seller_id is not string"
		res := entity.Response[error]{
			Code:   http.StatusUnauthorized,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	if sellerStr == "noId" {
		stat.IsSeller = false
	} else {
		stat.IsSeller = true
	}

	admin := session.Get("admin_id")
	if admin == nil {
		msg := "admin null"
		res := entity.Response[error]{
			Code:   http.StatusUnauthorized,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	adminStr, ok := admin.(string)
	if !ok {
		msg := "admin_id is not string"
		res := entity.Response[error]{
			Code:   http.StatusUnauthorized,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	if adminStr == "noId" {
		stat.IsAdmin = false
	} else {
		stat.IsAdmin = true
	}

	res := entity.Response[entity.Status]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   stat,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false, // set true if using HTTPS
		SameSite: http.SameSiteLaxMode,
	})
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

	res := entity.Response[string]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   "Logout Success",
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}
