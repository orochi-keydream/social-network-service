package repository

import (
	"context"
	"database/sql"
	"fmt"
	"social-network-service/internal/model"
	"strings"
	"time"
)

type UserDto struct {
	UserId     string    `db:"user_id"`
	FirstName  string    `db:"first_name"`
	SecondName string    `db:"second_name"`
	Gender     int       `db:"gender"`
	Birthday   time.Time `db:"birthday"`
	Biography  string    `db:"biography"`
	City       string    `db:"city"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Add(ctx context.Context, user *model.User, tx *sql.Tx) error {
	const query = "insert into users (user_id, first_name, second_name, gender, birthday, biography, city) values ($1, $2, $3, $4, $5, $6, $7)"

	var ec ExecutionContext

	if tx == nil {
		ec = r.db
	} else {
		ec = tx
	}

	dto := UserDto{
		UserId:     string(user.UserId),
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Gender:     int(user.Gender),
		Birthday:   user.Birthdate,
		Biography:  user.Biography,
		City:       user.City,
	}

	// TODO: Check how we can pass 'user'.
	_, err := ec.ExecContext(ctx, query, dto.UserId, dto.FirstName, dto.SecondName, dto.Gender, dto.Birthday, dto.Biography, dto.City)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Get(ctx context.Context, userId model.UserId, tx *sql.Tx) (*model.User, error) {
	const query = "select user_id, first_name, second_name, gender, birthday, biography, city from users where user_id = $1"

	var ec ExecutionContext

	if tx == nil {
		ec = r.db
	} else {
		ec = tx
	}

	row := ec.QueryRowContext(ctx, query, userId)

	err := row.Err()

	if err != nil {
		return nil, err
	}

	var dto UserDto

	err = row.Scan(&dto.UserId, &dto.FirstName, &dto.SecondName, &dto.Gender, &dto.Birthday, &dto.Biography, &dto.City)

	if err != nil {
		return nil, err
	}

	user := model.User{
		UserId:     model.UserId(dto.UserId),
		FirstName:  dto.FirstName,
		SecondName: dto.SecondName,
		Gender:     model.Gender(dto.Gender),
		Birthdate:  dto.Birthday,
		Biography:  dto.Biography,
		City:       dto.City,
	}

	return &user, nil
}

func (r *UserRepository) SearchUsers(ctx context.Context, firstName string, secondName string, tx *sql.Tx) ([]*model.User, error) {
	var b strings.Builder

	// TODO: Think about how to get rid of 'where 1 = 1'.
	b.Write([]byte("select user_id, first_name, second_name, gender, birthday, biography, city from users where 1 = 1"))

	// TODO: Think about another way to build the query with optional filters,

	params := []any{}
	paramNumber := 0

	// TODO: Think about more efficient way to search text.
	if firstName != "" {
		params = append(params, firstName+"%")
		paramNumber++
		b.WriteString(fmt.Sprintf(" and first_name like $%v", paramNumber))
	}

	if secondName != "" {
		params = append(params, secondName+"%")
		paramNumber++
		b.WriteString(fmt.Sprintf(" and second_name like $%v", paramNumber))
	}

	var ec ExecutionContext

	if tx == nil {
		ec = r.db
	} else {
		ec = tx
	}

	rows, err := ec.QueryContext(ctx, b.String(), params...)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	users := []*model.User{}

	for rows.Next() {
		var dto UserDto

		rows.Scan(&dto.UserId, &dto.FirstName, &dto.SecondName, &dto.Gender, &dto.Birthday, &dto.Biography, &dto.City)

		user := model.User{
			UserId:     model.UserId(dto.UserId),
			FirstName:  dto.FirstName,
			SecondName: dto.SecondName,
			Gender:     model.Gender(dto.Gender),
			Birthdate:  dto.Birthday,
			Biography:  dto.Biography,
			City:       dto.City,
		}

		users = append(users, &user)
	}

	return users, nil
}
