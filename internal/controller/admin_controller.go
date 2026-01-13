package controller

import (
	"net/http"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
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
	// session := sessions.Default(c)
	// id := session.Get("admin_id")
	// if id == nil || id == "noId" {
	// 	msg := "id null"
	// 	res := entity.Response[error]{
	// 		Code:   http.StatusUnauthorized,
	// 		Status: "error",
	// 		Data:   nil,
	// 		Error:  &msg,
	// 	}
	// 	c.JSON(http.StatusUnauthorized, res)
	// 	return
	// }

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
