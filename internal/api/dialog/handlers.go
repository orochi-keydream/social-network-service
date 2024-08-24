package dialog

import (
	"net/http"
	"social-network-service/internal/api/common"
	"social-network-service/internal/middleware"
	"social-network-service/internal/model"

	"github.com/gin-gonic/gin"
)

func RegisterDialogEndpoints(service common.DialogService, jwtService common.JwtService, e *gin.Engine, auth middleware.AuthMiddleware) {
	dialogRouter := e.Group("/dialog")

	dialogRouter.Use(gin.HandlerFunc(auth))

	dialogRouter.POST("/:id/send", NewSendMessageHandler(service, jwtService))
	dialogRouter.GET("/:id/list", NewGetMessagesHandler(service, jwtService))

	dialogRouter.GET("/unread/count", NewGetUnreadCountTotal(service, jwtService))
	dialogRouter.GET(":id/unread/count", NewGetUnreadCount(service, jwtService))
	dialogRouter.POST("/unread/read", MarkMessagesAsRead(service, jwtService))
}

// @Summary Returns messages.
// @Tags Dialog
// @Accept json
// @Param id path string true " "
// @Produce json
// @Success 200 {object} GetMessagesResponse
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Router /dialog/{id}/list [get]
// @Security bearer
func NewGetMessagesHandler(service common.DialogService, jwtService common.JwtService) func(*gin.Context) {
	return func(c *gin.Context) {
		fromUsedId, err := jwtService.GetUserId(c)

		if err != nil {
			c.Error(err)
			return
		}

		toUserid := c.Param("id")

		if toUserid == "" {
			c.Error(model.NewClientError("no target user ID passed", nil))
			return
		}

		cmd := model.GetMessagesCommand{
			FromUserId: model.UserId(fromUsedId),
			ToUserId:   model.UserId(toUserid),
		}

		messages, err := service.GetMessages(c.Request.Context(), cmd)

		if err != nil {
			c.Error(err)
			return
		}

		resp := mapToGetMessagesResponse(messages)

		c.JSON(http.StatusOK, resp)
	}
}

// @Summary Sends a message.
// @Tags Dialog
// @Accept json
// @Param id path string true " "
// @Param request body SendMessageRequest true " "
// @Produce json
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Router /dialog/{id}/send [post]
// @Security bearer
func NewSendMessageHandler(service common.DialogService, jwtService common.JwtService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId, err := jwtService.GetUserId(c)

		if err != nil {
			c.Error(err)
			return
		}

		fromUserId := userId
		toUserId := c.Param("id")

		req := SendMessageRequest{}
		err = c.BindJSON(&req)

		if err != nil {
			c.Error(model.NewClientError("Failed to parse JSON", err))
			return
		}

		cmd := model.SendMessageCommand{
			FromUserId: model.UserId(fromUserId),
			ToUserId:   model.UserId(toUserId),
			Text:       req.Text,
		}

		err = service.SendMessage(c.Request.Context(), cmd)

		if err != nil {
			c.Error(err)
		}
	}
}

// @Summary Returns total number of unread messages.
// @Tags Dialog
// @Produce json
// @Success 200 {object} NewGetUnreadCountTotalResponse
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Router /dialog/unread/count [get]
// @Security bearer
func NewGetUnreadCountTotal(service common.DialogService, jwtService common.JwtService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId, err := jwtService.GetUserId(c)

		if err != nil {
			c.Error(err)
			return
		}

		cmd := model.NewGetUnreadCountTotalCommand{
			UserId: userId,
		}

		count, err := service.NewGetUnreadCountTotal(c.Request.Context(), cmd)

		if err != nil {
			c.Error(err)
			return
		}

		resp := NewGetUnreadCountTotalResponse{
			Count: count,
		}

		c.JSON(http.StatusOK, resp)
	}
}

// @Summary Returns number of unread messages in a specified chat.
// @Tags Dialog
// @Param id path string true " "
// @Produce json
// @Success 200 {object} NewGetUnreadCountResponse
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Router /dialog/{id}/unread/count [get]
// @Security bearer
func NewGetUnreadCount(service common.DialogService, jwtService common.JwtService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId, err := jwtService.GetUserId(c)

		if err != nil {
			c.Error(err)
			return
		}

		chatUserId := c.Param("id")

		if chatUserId == "" {
			c.Error(model.NewClientError("no target user ID passed", nil))
			return
		}

		cmd := model.NewGetUnreadCountCommand{
			CurrentUserId: userId,
			ChatUserId:    model.UserId(chatUserId),
		}

		count, err := service.NewGetUnreadCount(c.Request.Context(), cmd)

		if err != nil {
			c.Error(err)
			return
		}

		resp := NewGetUnreadCountResponse{
			Count: count,
		}

		c.JSON(http.StatusOK, resp)
	}
}

// @Summary Returns number of unread messages in a specified chat.
// @Tags Dialog
// @Param request body MarkMessagesAsReadRequest true " "
// @Produce json
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Router /dialog/{id}/unread/read [post]
// @Security bearer
func MarkMessagesAsRead(service common.DialogService, jwtService common.JwtService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId, err := jwtService.GetUserId(c)

		if err != nil {
			c.Error(err)
			return
		}

		req := MarkMessagesAsReadRequest{}
		err = c.BindJSON(&req)

		if err != nil {
			c.Error(model.NewClientError("Failed to parse JSON", err))
			return
		}

		messageIds := make([]model.MessageId, len(req.MessageIds))

		for i, messageId := range req.MessageIds {
			messageIds[i] = model.MessageId(messageId)
		}

		cmd := model.MarkMessagesAsReadCommand{
			UserId:     userId,
			MessageIds: messageIds,
		}

		err = service.MarkMessagesAsRead(c.Request.Context(), cmd)

		if err != nil {
			c.Error(err)
		}
	}
}
