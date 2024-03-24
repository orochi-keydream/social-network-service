package post

import (
	"net/http"
	"social-network-service/internal/api/common"
	"social-network-service/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterPostEndpoints(service common.PostService, jwtService common.JwtService, e *gin.Engine) gin.RouterGroup {
	postRouter := e.Group("/post")

	postRouter.POST("/create", NewCreatePostHandler(service, jwtService))
	postRouter.PUT("/update", NewUpdatePostHandler(service, jwtService))
	postRouter.PUT("/delete/:id", NewDeletePostHandler(service, jwtService))
	postRouter.GET("/get/:id", NewGetPostHandler(service))
	postRouter.GET("/feed", NewReadFeedHandler(service, jwtService))

	return *postRouter
}

// @Summary Creates a post.
// @Tags Post
// @Accept json
// @Param request body CreatePostRequest true " "
// @Produce json
// @Success 200 {object} CreatePostResponse
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Router /post/create [post]
// @Security bearer
func NewCreatePostHandler(service common.PostService, jwtService common.JwtService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId, err := jwtService.GetUserId(c)

		if err != nil {
			c.Error(err)
			return
		}

		var req CreatePostRequest
		err = c.BindJSON(&req)

		if err != nil {
			c.Error(model.NewClientError("failed to parse body", err))
			return
		}

		cmd := model.CreatePostCommand{
			AuthorUserId: model.UserId(userId),
			Text:         req.Text,
		}

		postId, err := service.CreatePost(c.Request.Context(), cmd)

		if err != nil {
			c.Error(err)
			return
		}

		resp := CreatePostResponse{
			PostId: string(postId),
		}

		c.JSON(http.StatusOK, resp)
	}
}

// @Summary Returns specified post.
// @Tags Post
// @Accept json
// @Param id path string true " "
// @Produce json
// @Success 200 {object} GetPostResponse
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Router /post/get/{id} [get]
// @Security bearer
func NewGetPostHandler(service common.PostService) func(*gin.Context) {
	return func(c *gin.Context) {
		postIdStr := c.Param("id")

		if postIdStr == "" {
			c.Error(model.NewClientError("post ID was not provided", nil))
			return
		}

		postId := model.PostId(postIdStr)

		post, err := service.GetPost(c.Request.Context(), postId)

		if err != nil {
			c.Error(err)
			return
		}

		resp := GetPostResponse{
			PostId:       string(post.PostId),
			AuthorUserId: string(post.AuthorUserId),
			Text:         post.Text,
		}

		c.JSON(http.StatusOK, resp)
	}
}

// @Summary Returns posts.
// @Tags Post
// @Accept json
// @Param offset query integer false " "
// @Param limit query integer false " "
// @Produce json
// @Success 200 {object} ReadFeedResponse
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Router /post/feed [get]
// @Security bearer
func NewReadFeedHandler(service common.PostService, jwtService common.JwtService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId, err := jwtService.GetUserId(c)

		if err != nil {
			c.Error(err)
			return
		}

		offsetStr := c.DefaultQuery("offset", "0")
		limitStr := c.DefaultQuery("limit", "10")

		offset, err := strconv.Atoi(offsetStr)

		if err != nil {
			c.Error(model.NewClientError("wrong offset value", err))
			return
		}

		if offset < 0 {
			c.Error(model.NewClientError("offset must be positive", nil))
			return
		}

		limit, err := strconv.Atoi(limitStr)

		if err != nil {
			c.Error(model.NewClientError("wrong limit value", err))
			return
		}

		if limit < 0 {
			c.Error(model.NewClientError("limit must be positive", nil))
			return
		}

		cmd := model.ReadPostsCommand{
			UserId: model.UserId(userId),
			Offset: offset,
			Limit:  limit,
		}

		posts, err := service.ReadPosts(c.Request.Context(), cmd)

		if err != nil {
			c.Error(err)
			return
		}

		resp := mapToReadFeedResponse(posts)

		c.JSON(http.StatusOK, resp)
	}
}

// @Summary Updates specified post.
// @Tags Post
// @Accept json
// @Param request body UpdatePostRequest true " "
// @Produce json
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Router /post/update [put]
// @Security bearer
func NewUpdatePostHandler(service common.PostService, jwtService common.JwtService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId, err := jwtService.GetUserId(c)

		if err != nil {
			c.Error(err)
			return
		}

		var req UpdatePostRequest
		err = c.BindJSON(&req)

		if err != nil {
			c.Error(model.NewClientError("failed to parse body", err))
			return
		}

		cmd := model.UpdatePostCommand{
			AuthorUserId: model.UserId(userId),
			PostId:       model.PostId(req.PostId),
			Text:         req.Text,
		}

		err = service.UpdatePost(c.Request.Context(), cmd)

		if err != nil {
			c.Error(err)
			return
		}
	}
}

// @Summary Deletes specified post.
// @Tags Post
// @Accept json
// @Param id path string true " "
// @Produce json
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Router /post/delete/{id} [put]
// @Security bearer
func NewDeletePostHandler(service common.PostService, jwtService common.JwtService) func(*gin.Context) {
	return func(c *gin.Context) {
		// BUG: 500 when no JWT provided and an empty string is given as post ID.
		userId, err := jwtService.GetUserId(c)

		if err != nil {
			c.Error(err)
		}

		postId := c.Param("id")

		if postId == "" {
			c.Error(model.NewClientError("post ID was not specified", nil))
			return
		}

		cmd := model.DeletePostCommand{
			AuthorUserId: model.UserId(userId),
			PostId:       model.PostId(postId),
		}

		err = service.DeletePost(c.Request.Context(), cmd)

		if err != nil {
			c.Error(err)
		}
	}
}
