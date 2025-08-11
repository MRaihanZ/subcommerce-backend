package controller

import (
	"net/http"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetUserHandler(c *gin.Context) {
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

	user, err := model.GetUserById(id)
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

	if user == nil {
		msg := "user not found"
		res := entity.Response[*entity.User]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
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
