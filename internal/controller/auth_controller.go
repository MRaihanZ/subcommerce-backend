package controller

import (
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserHandler(c *gin.Context) {
	var req entity.SignUp
	if err := c.BindJSON(&req); err != nil {
		msg := err.Error()
		res := entity.Response[entity.SignUp]{
			Code:   400,
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
			Code:   500,
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(500, res)
		return
	}

	res := entity.Response[*uuid.UUID]{
		Code:   200,
		Status: "ok",
		Data:   user,
		Error:  nil,
	}
	c.JSON(200, res)
}

func VerifyUserHandler(c *gin.Context) {
	var req entity.SignIn
	if err := c.BindJSON(&req); err != nil {
		msg := err.Error()
		res := entity.Response[entity.SignIn]{
			Code:   400,
			Status: "error",
			Data:   req,
			Error:  &msg,
		}
		c.JSON(400, res)
		return
	}

	user, err := model.GetUserByEmail(req.Email)
	if err != nil {
		msg := err.Error()
		res := entity.Response[*entity.SignIn]{
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
		res := entity.Response[*entity.SignIn]{
			Code:   404,
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(404, res)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		msg := err.Error()
		res := entity.Response[*entity.SignIn]{
			Code:   401,
			Status: "error",
			Data:   user,
			Error:  &msg,
		}
		c.JSON(401, res)
		return
	}

	res := entity.Response[*entity.SignIn]{
		Code:   200,
		Status: "ok",
		Data:   user,
		Error:  nil,
	}
	c.JSON(200, res)
}
