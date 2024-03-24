package post

import "social-network-service/internal/model"

func mapToReadFeedResponse(posts []*model.Post) *ReadFeedResponse {
	items := []ReadFeedResponseItem{}

	for _, post := range posts {
		item := ReadFeedResponseItem{
			PostId:       string(post.PostId),
			AuthorUserId: string(post.AuthorUserId),
			Text:         post.Text,
		}

		items = append(items, item)
	}

	return &ReadFeedResponse{
		Posts: items,
	}
}
