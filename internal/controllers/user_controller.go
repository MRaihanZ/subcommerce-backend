package controller

import (
	"github.com/MRaihanZ/subcommerce-backend/internal/entities"
	model "github.com/MRaihanZ/subcommerce-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUsersHandler(c *gin.Context) {
	users, err := model.GetAllUsers()
	if err != nil {
		msg := err.Error()
		res := entities.Response[[]entities.User]{
			Code:   "500",
			Status: "error",
			Data:   users,
			Error:  &msg,
		}
		c.JSON(500, res)
		return
	}

	if users == nil {
		msg := "no user found"
		res := entities.Response[[]entities.User]{
			Code:   "404",
			Status: "error",
			Data:   users,
			Error:  &msg,
		}
		c.JSON(404, res)
		return
	}

	res := entities.Response[[]entities.User]{
		Code:   "200",
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
		res := entities.Response[*entities.User]{
			Code:   "500",
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(500, res)
		return
	}

	if user == nil {
		msg := "user not found"
		res := entities.Response[*entities.User]{
			Code:   "404",
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(404, res)
		return
	}

	res := entities.Response[*entities.User]{
		Code:   "200",
		Status: "ok",
		Data:   user,
		Error:  nil,
	}
	c.JSON(200, res)
}

func CreateUserHandler(c *gin.Context) {
	var req entities.SignUp
	if err := c.BindJSON(&req); err != nil {
		msg := err.Error()
		res := entities.Response[entities.SignUp]{
			Code:   "400",
			Status: "error",
			Data:   req,
			Error:  &msg,
		}
		c.JSON(400, res)
		return
	}
	user, err := model.CreateUser(req.Name, req.Email, req.Dob, req.Password)
	if err != nil {
		msg := err.Error()
		res := entities.Response[*uuid.UUID]{
			Code:   "500",
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(500, res)
		return
	}

	res := entities.Response[*uuid.UUID]{
		Code:   "200",
		Status: "ok",
		Data:   user,
		Error:  nil,
	}
	c.JSON(200, res)
}
