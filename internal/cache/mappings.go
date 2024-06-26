package cache

import (
	"social-network-service/internal/model"
	"time"
)

func toPostDto(post *model.Post) Post {
	return Post{
		PostId:       string(post.PostId),
		PublishedAt:  post.PublishedAt.Format(time.RFC3339),
		AuthorUserId: string(post.AuthorUserId),
		Text:         post.Text,
	}
}

func toPostModel(dto Post) (*model.Post, error) {
	publishedAt, err := time.Parse(time.RFC3339, dto.PublishedAt)

	if err != nil {
		return nil, err
	}

	return &model.Post{
		PostId:       model.PostId(dto.PostId),
		PublishedAt:  publishedAt,
		AuthorUserId: model.UserId(dto.AuthorUserId),
		Text:         dto.Text,
	}, nil
}
