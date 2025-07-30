package controller

import (
	"net/http"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/gin-gonic/gin"
)

func GetUsersHandler(c *gin.Context) {
	users, err := model.GetAllUsers()
	if err != nil {
		msg := err.Error()
		res := entity.Response[[]entity.User]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   users,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if users == nil {
		msg := "no user found"
		res := entity.Response[[]entity.User]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   users,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[[]entity.User]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   users,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetUserHandler(c *gin.Context) {
	id := c.Param("id")
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
