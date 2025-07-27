package controller

import (
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUserHandler(c *gin.Context) {
	var req entity.SignUp
	if err := c.BindJSON(&req); err != nil {
		msg := err.Error()
		res := entity.Response[entity.SignUp]{
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
		res := entity.Response[*uuid.UUID]{
			Code:   "500",
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(500, res)
		return
	}

	res := entity.Response[*uuid.UUID]{
		Code:   "200",
		Status: "ok",
		Data:   user,
		Error:  nil,
	}
	c.JSON(200, res)
}
