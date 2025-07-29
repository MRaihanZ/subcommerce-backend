package controller

import (
	"net/http"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserHandler(c *gin.Context) {
	var req entity.SignUp
	if err := c.BindJSON(&req); err != nil {
		msg := err.Error()
		res := entity.Response[entity.SignUp]{
			Code:   http.StatusBadRequest,
			Status: "error",
			Data:   req,
			Error:  &msg,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}
	user, err := model.CreateUser(req.Name, req.Email, req.Dob, req.Password)
	if err != nil {
		msg := err.Error()
		res := entity.Response[*uuid.UUID]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := entity.Response[*uuid.UUID]{
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
		res := entity.Response[entity.SignIn]{
			Code:   http.StatusBadRequest,
			Status: "error",
			Data:   req,
			Error:  &msg,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	csrfToken := c.GetHeader("X-CSRF-TOKEN")

	user, err := model.GetUserByEmail(req.Email)
	if err != nil {
		msg := err.Error()
		res := entity.Response[*entity.SignIn]{
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
		res := entity.Response[*entity.SignIn]{
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
		res := entity.Response[*entity.SignIn]{
			Code:   http.StatusUnauthorized,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	// save session
	session := sessions.Default(c)
	session.Set("user_id", user.Id)
	session.Set("csrf_token", csrfToken)
	session.Set("is_logged_in", true)
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

func CheckStatus(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("user_id")
	if id == "" {
		res := entity.Response[*entity.SignIn]{
			Code:   http.StatusUnauthorized,
			Status: "ok",
			Data:   nil,
			Error:  nil,
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	stat := entity.Status{Status: "authenticated", Id: id}
	res := entity.Response[entity.Status]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   stat,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)

}
