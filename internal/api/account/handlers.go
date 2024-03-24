package account

import (
	"net/http"
	"social-network-service/internal/api/common"
	"social-network-service/internal/model"

	"github.com/gin-gonic/gin"
)

func RegisterAccountEndpoints(service common.AccountService, e *gin.Engine) gin.RouterGroup {
	accountRouter := e.Group("")

	accountRouter.POST("/login", NewLoginHandler(service))
	accountRouter.POST("/user/register", NewRegisterHandler(service))

	return *accountRouter
}

// @Summary Sign in using user ID.
// @Tags Account
// @Accept json
// @Param request body LoginRequest true " "
// @Produce json
// @Success 200 {object} LoginResponse
// @Failure 400 {object} object
// @Router /login [post]
func NewLoginHandler(service common.AccountService) func(*gin.Context) {
	return func(c *gin.Context) {
		req := LoginRequest{}
		err := c.ShouldBindJSON(&req)

		if err != nil {
			c.Error(model.NewClientError("failed to parse body", err))
			return
		}

		userId := model.UserId(req.UserId)

		token, err := service.Login(c.Request.Context(), userId, req.Password)

		if err != nil {
			c.Error(err)
			return
		}

		resp := LoginResponse{
			Token: token,
		}

		c.JSON(http.StatusOK, resp)
	}
}

// @Summary Registers a new user.
// @Tags Account
// @Accept json
// @Param request body RegisterRequest true " "
// @Produce json
// @Success 200 {object} RegisterResponse
// @Failure 400 {object} object
// @Router /user/register [post]
func NewRegisterHandler(service common.AccountService) func(*gin.Context) {
	return func(c *gin.Context) {
		req := RegisterRequest{}
		err := c.ShouldBindJSON(&req)

		if err != nil {
			c.Error(model.NewClientError("failed to parse body", err))
			return
		}

		registerUserCmd, err := mapRegisterRequestToCommand(&req)

		if err != nil {
			c.Error(model.NewClientError("failed to build command", err))
			return
		}

		userId, err := service.RegisterUser(c.Request.Context(), registerUserCmd)

		if err != nil {
			c.Error(err)
			return
		}

		resp := RegisterResponse{
			UserId: string(userId),
		}

		c.JSON(http.StatusOK, resp)
	}
}
