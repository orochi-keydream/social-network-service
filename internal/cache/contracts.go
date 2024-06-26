package cache

type Post struct {
	PostId       string `json:"postId"`
	PublishedAt  string `json:"publishedAt"`
	AuthorUserId string `json:"authorUserId"`
	Text         string `json:"text"`
}
