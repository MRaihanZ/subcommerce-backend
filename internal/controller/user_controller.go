package controller

import (
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/gin-gonic/gin"
)

func GetUsersHandler(c *gin.Context) {
	users, err := model.GetAllUsers()
	if err != nil {
		msg := err.Error()
		res := entity.Response[[]entity.User]{
			Code:   500,
			Status: "error",
			Data:   users,
			Error:  &msg,
		}
		c.JSON(500, res)
		return
	}

	if users == nil {
		msg := "no user found"
		res := entity.Response[[]entity.User]{
			Code:   404,
			Status: "error",
			Data:   users,
			Error:  &msg,
		}
		c.JSON(404, res)
		return
	}

	res := entity.Response[[]entity.User]{
		Code:   200,
		Status: "ok",
		Data:   users,
		Error:  nil,
	}
	c.JSON(200, res)
}

func GetUserHandler(c *gin.Context) {
	id := c.Param("id")
	user, err := model.GetUserById(id)
	if err != nil {
		msg := err.Error()
		res := entity.Response[*entity.User]{
			Code:   500,
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(500, res)
		return
	}

	if user == nil {
		msg := "user not found"
		res := entity.Response[*entity.User]{
			Code:   404,
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(404, res)
		return
	}

	res := entity.Response[*entity.User]{
		Code:   200,
		Status: "ok",
		Data:   user,
		Error:  nil,
	}
	c.JSON(200, res)
}
