package post

type CreatePostRequest struct {
	// Content of the post.
	Text string `json:"text" binding:"required"`
}

type CreatePostResponse struct {
	// ID of the post in UUIDv4 format.
	PostId string `json:"postId"`
}

type GetPostResponse struct {
	// ID of the post in UUIDv4 format.
	PostId string `json:"postId"`

	// User ID of the post author in UUIDv4 format.
	AuthorUserId string `json:"authorUserId"`

	// Content of the post.
	Text string `json:"text"`
}

type ReadFeedResponse struct {
	// List of posts.
	Posts []ReadFeedResponseItem `json:"posts"`
}

type ReadFeedResponseItem struct {
	// ID of the post in UUIDv4 format.
	PostId string `json:"postId"`

	// User ID of the post author in UUIDv4 format.
	AuthorUserId string `json:"authorUserId"`

	// Content of the post.
	Text string `json:"text"`
}

type UpdatePostRequest struct {
	// ID of the post in UUIDv4 format.
	PostId string `json:"postId" binding:"required"`

	// Content of the message.
	Text string `json:"text" binding:"required"`
}
