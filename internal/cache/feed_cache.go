package cache

import (
	"context"
	"encoding/json"
	"slices"
	"social-network-service/internal/model"

	"github.com/redis/go-redis/v9"
)

type FeedCache struct {
	client *redis.Client
}

func NewFeedCache(client *redis.Client) *FeedCache {
	return &FeedCache{
		client: client,
	}
}

func (c *FeedCache) RecreateFeed(userId model.UserId, posts []*model.Post) error {
	ctx := context.Background()

	if len(posts) == 0 {
		delCmdRes := c.client.Del(ctx, string(userId))

		return delCmdRes.Err()
	}

	byteSlice := make([]interface{}, len(posts))

	for i, post := range posts {
		dto := toPostDto(post)

		bytes, err := json.Marshal(dto)

		if err != nil {
			return err
		}

		byteSlice[i] = bytes
	}

	// TODO: Think how to solve the case where Del() is called, but LPush is not().
	delCmdRes := c.client.Del(ctx, string(userId))

	if delCmdRes.Err() != nil {
		return delCmdRes.Err()
	}

	pushCmdRes := c.client.LPush(ctx, string(userId), byteSlice...)

	if pushCmdRes.Err() != nil {
		return pushCmdRes.Err()
	}

	return nil
}

func (c *FeedCache) GetFeed(userId model.UserId, offset, limit int) ([]*model.Post, error) {
	ctx := context.Background()

	if limit == 0 {
		return []*model.Post{}, nil
	}

	// Redis uses right boundary including specified index.
	// For offset = 0 and limit = 10 the following command should be generated:
	// LRANGE {key} 0 9
	limit--

	rangeRes := c.client.LRange(ctx, string(userId), int64(offset), int64(offset+limit))

	if rangeRes.Err() != nil {
		return nil, rangeRes.Err()
	}

	byteSlice := [][]byte{}

	err := rangeRes.ScanSlice(&byteSlice)

	if err != nil {
		return nil, err
	}

	dtos := make([]Post, len(byteSlice))

	for i, bytes := range byteSlice {
		var dto Post
		err := json.Unmarshal(bytes, &dto)

		if err != nil {
			return nil, err
		}

		dtos[i] = dto
	}

	posts := make([]*model.Post, len(dtos))

	for i, dto := range dtos {
		post, err := toPostModel(dto)

		if err != nil {
			return nil, err
		}

		posts[i] = post
	}

	return posts, nil
}

func (c *FeedCache) AddPost(userId model.UserId, post *model.Post) error {
	ctx := context.Background()

	dto := toPostDto(post)

	bytes, err := json.Marshal(dto)

	if err != nil {
		return err
	}

	pushRes := c.client.LPush(ctx, string(userId), bytes)

	if pushRes.Err() != nil {
		return pushRes.Err()
	}

	trimCmd := c.client.LTrim(ctx, string(userId), 0, 1000)

	if trimCmd.Err() != nil {
		return trimCmd.Err()
	}

	return nil
}

func (c *FeedCache) UpdatePost(userId model.UserId, post *model.Post) error {
	ctx := context.Background()

	rangeCmdRes := c.client.LRange(ctx, string(userId), 0, 1000)

	if rangeCmdRes.Err() != nil {
		return rangeCmdRes.Err()
	}

	byteSlice := [][]byte{}

	err := rangeCmdRes.ScanSlice(&byteSlice)

	if err != nil {
		return err
	}

	dtos := make([]Post, len(byteSlice))

	for i, bytes := range byteSlice {
		var dto Post
		err := json.Unmarshal(bytes, &dto)

		if err != nil {
			return err
		}

		dtos[i] = dto
	}

	idx := slices.IndexFunc(dtos, func(p Post) bool {
		return p.PostId == string(post.PostId)
	})

	if idx == -1 {
		return nil
	}

	dto := toPostDto(post)

	bytes, err := json.Marshal(dto)

	if err != nil {
		return err
	}

	setCmdRes := c.client.LSet(ctx, string(userId), int64(idx), bytes)

	if setCmdRes.Err() != nil {
		return setCmdRes.Err()
	}

	return nil
}
