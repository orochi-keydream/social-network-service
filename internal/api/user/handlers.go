package user

import (
	"net/http"
	"social-network-service/internal/api/common"
	"social-network-service/internal/middleware"
	"social-network-service/internal/model"

	"github.com/gin-gonic/gin"
)

func RegisterUserClosedEndpoints(service common.UserService, jwtService common.JwtService, e *gin.Engine, auth middleware.AuthMiddleware) {
	userRouterOpen := e.Group("")

	userRouterOpen.GET("/user/get/:id", NewGetUserHandler(service))
	userRouterOpen.GET("/user/search", NewSearchUsersHandler(service))

	userRouterClosed := e.Group("")

	userRouterClosed.Use(gin.HandlerFunc(auth))

	userRouterClosed.PUT("/friend/set/:id", NewAddFriendHandler(service, jwtService))
	userRouterClosed.PUT("/friend/delete/:id", NewRemoveFriendHandler(service, jwtService))
}

// @Summary Returns user by ID.
// @Tags User
// @Accept json
// @Param id path string true " "
// @Produce json
// @Success 200 {object} GetUserResponse
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Router /user/get/{id} [get]
func NewGetUserHandler(service common.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		userIdStr := c.Param("id")

		if userIdStr == "" {
			c.Error(model.NewClientError("no user ID provided", nil))
			return
		}

		userId := model.UserId(userIdStr)

		user, err := service.GetUser(c.Request.Context(), userId)

		if err != nil {
			c.Error(err)
			return
		}

		resp, err := mapToGetUserResponse(user)

		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

// @Summary Returns found users.
// @Tags User
// @Accept json
// @Param first_name query string false " "
// @Param second_name query string false " "
// @Produce json
// @Success 200 {object} SearchUsersResponse
// @Failure 400 {object} object
// @Router /user/search [get]
func NewSearchUsersHandler(service common.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		firstName := c.Query("first_name")
		secondName := c.Query("second_name")

		users, err := service.SearchUsers(c.Request.Context(), firstName, secondName)

		if err != nil {
			c.Error(err)
			return
		}

		resp, err := mapToSearchUsersResponse(users)

		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

// @Summary Adds a friend.
// @Tags User
// @Accept json
// @Param id path string true " "
// @Produce json
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Router /friend/set/{id} [put]
// @Security bearer
func NewAddFriendHandler(service common.UserService, jwtService common.JwtService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId, err := jwtService.GetUserId(c)

		if err != nil {
			c.Error(err)
			return
		}

		friendUserIdStr := c.Param("id")

		if friendUserIdStr == "" {
			c.Error(model.NewClientError("user ID not provided", nil))
			return
		}

		friendUserId := model.UserId(friendUserIdStr)

		err = service.AddFriend(c.Request.Context(), userId, friendUserId)

		if err != nil {
			c.Error(err)
			return
		}
	}
}

// @Summary Removes specified friend.
// @Tags User
// @Accept json
// @Param id path string true " "
// @Produce json
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Router /friend/delete/{id} [put]
// @Security bearer
func NewRemoveFriendHandler(service common.UserService, jwtService common.JwtService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId, err := jwtService.GetUserId(c)

		if err != nil {
			c.Error(err)
			return
		}

		friendUserIdStr := c.Param("id")

		if friendUserIdStr == "" {
			c.Error(model.NewClientError("user ID not provided", nil))
			return
		}

		friendUserId := model.UserId(friendUserIdStr)

		err = service.RemoveFriend(c.Request.Context(), userId, friendUserId)

		if err != nil {
			c.Error(err)
			return
		}
	}
}
