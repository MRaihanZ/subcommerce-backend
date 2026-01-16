package controller

import (
	"net/http"
	"strconv"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetAllConversationsHandler(c *gin.Context) {
	// get session
	session := sessions.Default(c)
	userId := session.Get("user_id")
	sellerId := session.Get("seller_id")

	// cek session if not login
	if userId == nil || sellerId == nil {
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

	var queryId interface{}
	stateQuery := c.Query("state")
	switch stateQuery {
	case "user":
		if userId == "noId" {
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
		queryId = userId
	case "seller":
		if sellerId == "noId" {
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
		queryId = sellerId
	default:
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

	conversation, err := model.GetAllConversations(queryId, stateQuery)
	if err != nil {
		msg := err.Error()
		res := entity.Response[*error]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if conversation == nil {
		msg := "Conversation room not found"
		res := entity.Response[*error]{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := entity.Response[[]entity.Conversation]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   conversation,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetConversationHandler(c *gin.Context) {
	// get session
	session := sessions.Default(c)
	userId := session.Get("user_id")
	sellerId := session.Get("seller_id")

	// cek session if not login
	if userId == nil || sellerId == nil {
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

	paramId := c.Param("id")
	if paramId == "" {
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

	// variable for user or seller id
	var userA, userB interface{}
	var err error
	var user *entity.User
	var seller *entity.Seller
	var conv entity.ConversationRoom

	// cek user or seller id if wrong query state, and set userA and userB accordingly
	stateQuery := c.Query("state")
	switch stateQuery {
	case "user":
		if userId == "noId" {
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

		// to get seller name and image
		seller, err = model.GetSellerById(paramId)
		if err != nil {
			msg := err.Error()
			res := entity.Response[*error]{
				Code:   http.StatusInternalServerError,
				Status: "error",
				Data:   nil,
				Error:  &msg,
			}
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		if seller == nil {
			msg := "seller not found"
			res := entity.Response[*error]{
				Code:   http.StatusNotFound,
				Status: "error",
				Data:   nil,
				Error:  &msg,
			}
			c.JSON(http.StatusNotFound, res)
			return
		}

		conv.Name = seller.Name
		conv.Img = seller.Img

		userA = userId
		userB = paramId
	case "seller":
		if sellerId == "noId" {
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

		// to get user name and image
		user, err = model.GetUserById(paramId)
		if err != nil {
			msg := err.Error()
			res := entity.Response[*error]{
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
			res := entity.Response[*entity.User]{
				Code:   http.StatusNotFound,
				Status: "error",
				Data:   user,
				Error:  &msg,
			}
			c.JSON(http.StatusNotFound, res)
			return
		}

		conv.Name = user.Name
		conv.Img = user.Img

		userA = paramId
		userB = sellerId
	default:
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

	// get conversation id
	convId, err := model.GetConversationsById(userA, userB)
	if err != nil {
		msg := err.Error()
		res := entity.Response[*string]{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   convId,
			Error:  &msg,
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// if conversation id not found, create new one
	if convId == nil {
		convId, err = model.CreateConversation(userA, userB)
		if err != nil {
			msg := err.Error()
			res := entity.Response[*string]{
				Code:   http.StatusInternalServerError,
				Status: "error",
				Data:   convId,
				Error:  &msg,
			}
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		if convId == nil {
			msg := "message room sudah terdaftar"
			res := entity.Response[error]{
				Code:   http.StatusConflict,
				Status: "error",
				Data:   nil,
				Error:  &msg,
			}
			c.JSON(http.StatusConflict, res)
			return
		}
	}

	conv.Id = *convId

	res := entity.Response[entity.ConversationRoom]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   conv,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func GetMessagesHandler(c *gin.Context) {
	// get session
	session := sessions.Default(c)
	userId := session.Get("user_id")
	sellerId := session.Get("seller_id")

	// cek session if not login
	if userId == nil || sellerId == nil {
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

	paramId := c.Param("id")
	if paramId == "" {
		msg := "conversation id is required"
		res := entity.Response[error]{
			Code:   http.StatusBadRequest,
			Status: "error",
			Data:   nil,
			Error:  &msg,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	stateQuery := c.Query("state")
	switch stateQuery {
	case "user":
		if userId == "noId" {
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
	case "seller":
		if sellerId == "noId" {
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
	default:
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

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	before := c.Query("before") // sent_at cursor
	beforeID := c.Query("msg_id")

	messages, err := model.GetMessagesByConvID(
		paramId,
		before,
		beforeID,
		limit,
	)
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

	res := entity.Response[[]entity.Message]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   messages,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}

func CreateMessageHandler(c *gin.Context) {
	// get session
	session := sessions.Default(c)
	userId := session.Get("user_id")
	sellerId := session.Get("seller_id")

	// cek session if not login
	if userId == nil || sellerId == nil {
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

	var isUser bool
	var senderId interface{}
	stateQuery := c.Query("state")
	switch stateQuery {
	case "user":
		if userId == "noId" {
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
		isUser = true
		senderId = userId
	case "seller":
		if sellerId == "noId" {
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
		isUser = false
		senderId = sellerId
	default:
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

	var req entity.MessageReceive
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

	message, err := model.CreateMessage(req.ConvId, senderId, isUser, req.Content)
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

	res := entity.Response[*string]{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   message,
		Error:  nil,
	}
	c.JSON(http.StatusOK, res)
}
