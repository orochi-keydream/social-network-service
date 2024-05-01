package repository

import (
	"context"
	"database/sql"
	"fmt"
	"social-network-service/internal/model"
	"strings"
	"time"

	"github.com/jackc/pgtype"
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
	const query = "insert into users (user_id, first_name, second_name, gender, birthday, biography, city, first_name_tsvector, second_name_tsvector) values ($1, $2, $3, $4, $5, $6, $7, to_tsvector('english', $2), to_tsvector('english', $3))"

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

func (r *UserRepository) AddBulk(ctx context.Context, users []*model.User, tx *sql.Tx) error {
	const query = `
	insert into users
	(
		user_id,
		first_name,
		second_name,
		gender,
		birthday,
		biography,
		city)
	select * from unnest
	(
		$1::text[],
		$2::text[],
		$3::text[],
		$4::integer[],
		$5::date[],
		$6::text[],
		$7::text[]
	)`

	var ec ExecutionContext

	if tx == nil {
		ec = r.db
	} else {
		ec = tx
	}

	count := len(users)

	userIds := make([]string, count)
	firstNames := make([]string, count)
	secondNames := make([]string, count)
	genders := make([]int, count)
	birthdates := make([]time.Time, count)
	biographies := make([]string, count)
	cities := make([]string, count)

	for i := 0; i < count; i++ {
		userIds[i] = string(users[i].UserId)
		firstNames[i] = users[i].FirstName
		secondNames[i] = users[i].SecondName
		genders[i] = int(users[i].Gender)
		birthdates[i] = users[i].Birthdate
		biographies[i] = users[i].Biography
		cities[i] = users[i].City
	}

	pgUserIds := pgtype.TextArray{}
	pgUserIds.Set(userIds)

	pgFirstNames := pgtype.TextArray{}
	pgFirstNames.Set(firstNames)

	pgSecondNames := pgtype.TextArray{}
	pgSecondNames.Set(secondNames)

	pgGenders := pgtype.Int4Array{}
	pgGenders.Set(genders)

	pgBirthdates := pgtype.DateArray{}
	pgBirthdates.Set(birthdates)

	pgBiographies := pgtype.TextArray{}
	pgBiographies.Set(biographies)

	pgCities := pgtype.TextArray{}
	pgCities.Set(cities)

	_, err := ec.ExecContext(ctx, query, pgUserIds, pgFirstNames, pgSecondNames, pgGenders, pgBirthdates, pgBiographies, pgCities)

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

	b.WriteString(" order by user_id limit 20")

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
