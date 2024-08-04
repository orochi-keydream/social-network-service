package repository

import (
	"context"
	"database/sql"
	"social-network-service/internal/metric"
	"social-network-service/internal/model"
)

type PostRepository struct {
	cf IConnectionFactory
}

func NewPostRepository(cf IConnectionFactory) *PostRepository {
	return &PostRepository{
		cf: cf,
	}
}

func (r *PostRepository) GetPost(ctx context.Context, postId model.PostId, tx *sql.Tx) (*model.Post, error) {
	const query = "select post_id, published_at, user_id, text from posts where post_id = $1"

	var ec IExecutionContext

	if tx == nil {
		ec = r.cf.GetMaster()
	} else {
		ec = tx
	}

	row := ec.QueryRowContext(ctx, query, postId)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var post model.Post
	err := row.Scan(&post.PostId, &post.PublishedAt, &post.AuthorUserId, &post.Text)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) GetPosts(ctx context.Context, userIds []model.UserId, offset int, limit int, tx *sql.Tx) ([]*model.Post, error) {
	const query = "select post_id, published_at, user_id, text from posts where user_id = any ($1) order by published_at desc offset $2 limit $3"

	var ec IExecutionContext

	if tx == nil {
		ec = r.cf.GetMaster()
	} else {
		ec = tx
	}

	rows, err := ec.QueryContext(ctx, query, userIds, offset, limit)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	posts := []*model.Post{}

	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.PostId, &post.PublishedAt, &post.AuthorUserId, &post.Text)

		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) GetPostsIncludingFriends(ctx context.Context, userId model.UserId, offset, limit int, tx *sql.Tx) ([]*model.Post, error) {
	const query = `
		select
			p.post_id,
			p.published_at,
			p.user_id,
			p.text
		from posts p
		left join user_friends uf on p.user_id = uf.friend_user_id
		where uf.user_id = $1 or p.user_id = $1
		order by p.published_at
		offset $2
		limit $3
		`

	var ec IExecutionContext

	if tx == nil {
		ec = r.cf.GetMaster()
	} else {
		ec = tx
	}

	rows, err := ec.QueryContext(ctx, query, userId, offset, limit)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	posts := []*model.Post{}

	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.PostId, &post.PublishedAt, &post.AuthorUserId, &post.Text)

		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) AddPost(ctx context.Context, post *model.Post, tx *sql.Tx) error {
	const query = "insert into posts (post_id, published_at, user_id, text) values ($1, $2, $3, $4)"

	// TODO: Extract it to an interface.
	metric.IncCreatePostAttempts()

	var ec IExecutionContext

	if tx == nil {
		ec = r.cf.GetMaster()
	} else {
		ec = tx
	}

	_, err := ec.ExecContext(ctx, query, post.PostId, post.PublishedAt, post.AuthorUserId, post.Text)

	if err != nil {
		return err
	}

	// TODO: Extract it to an interface.
	metric.IncCreatePostSuccessful()

	return nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, post *model.Post, tx *sql.Tx) error {
	const query = "update posts set text = $1 where post_id = $2"

	var ec IExecutionContext

	if tx == nil {
		ec = r.cf.GetMaster()
	} else {
		ec = tx
	}

	_, err := ec.ExecContext(ctx, query, post.Text, post.PostId)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostRepository) DeletePost(ctx context.Context, postId model.PostId, tx *sql.Tx) error {
	const query = "delete from posts where post_id = $1"

	var ec IExecutionContext

	if tx == nil {
		ec = r.cf.GetMaster()
	} else {
		ec = tx
	}

	_, err := ec.ExecContext(ctx, query, postId)

	if err != nil {
		return err
	}

	return nil
}
