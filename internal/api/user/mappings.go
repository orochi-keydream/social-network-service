package user

import (
	"fmt"
	"social-network-service/internal/api/common"
	"social-network-service/internal/model"
)

func mapToGetUserResponse(user *model.User) (*GetUserResponse, error) {
	gender, err := common.MapGenderFromModel(user.Gender)

	if err != nil {
		return nil, fmt.Errorf("failed to map gender: %w", err)
	}

	resp := &GetUserResponse{
		UserId:     string(user.UserId),
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Gender:     gender,
		Birthdate:  user.Birthdate.Format("2006-01-02"),
		Biography:  user.Biography,
		City:       user.City,
	}

	return resp, nil
}

func mapToSearchUsersResponse(users []*model.User) (*SearchUsersResponse, error) {
	items := []SearchUsersResponseItem{}

	for _, user := range users {

		gender, err := common.MapGenderFromModel(user.Gender)

		if err != nil {
			return nil, fmt.Errorf("failed to map gender: %w", err)
		}

		item := SearchUsersResponseItem{
			UserId:     string(user.UserId),
			FirstName:  user.FirstName,
			SecondName: user.SecondName,
			Gender:     gender,
			Birthdate:  user.Birthdate.Format("2006-01-02"),
			Biography:  user.Biography,
			City:       user.City,
		}

		items = append(items, item)
	}

	resp := &SearchUsersResponse{
		Users: items,
	}

	return resp, nil
}
