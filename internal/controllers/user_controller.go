package controller

import (
	model "github.com/MRaihanZ/subcommerce-backend/internal/models"
	"github.com/gin-gonic/gin"
)

// json template response:
// {
// 	"httpCode": 200,
// 	"status": "OK",
// 	"data": {
// 		"key": "value",
// 		"data": "data"
// 	}
//  "error": null
// }

// for error:
// {
// 	"httpCode": 200,
// 	"status": "OK",
// 	"data": null
//  "error": {
// 		"key": "value",
// 		"data": "data"
// 	}
// }

func GetUsers(c *gin.Context) {
	users, err := model.GetAllUsers()
	if err != nil {
		c.JSON(500, gin.H{"httpCode": "500", "error": err.Error()})
		return
	}
	c.JSON(200, users)
}
