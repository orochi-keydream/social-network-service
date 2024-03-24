package user

import "social-network-service/internal/api/common"

type GetUserResponse struct {
	// ID of the user in UUIDv4 format.
	UserId string `json:"user_id" example:"bbeb7da8-6d75-4419-9d94-91ec52bc506c"`

	// First name of the user.
	FirstName string `json:"first_name" example:"John"`

	// Second name of the user.
	SecondName string `json:"second_name" example:"Doe"`

	// Gender of the user ("Male" or "Female").
	Gender common.Gender `json:"gender" example:"Male"`

	// Birthday in the format "1990-12-31".
	Birthdate string `json:"birthdate" example:"1990-01-01"`

	// Biography of the user.
	Biography string `json:"biography" example:"Software developer"`

	// City of the user.
	City string `json:"city" example:"New York"`
}

type SearchUsersResponse struct {
	// List of found users.
	Users []SearchUsersResponseItem `json:"users"`
}

type SearchUsersResponseItem struct {
	// ID of the user in UUIDv4 format.
	UserId string `json:"user_id" example:"bbeb7da8-6d75-4419-9d94-91ec52bc506c"`

	// First name of the user.
	FirstName string `json:"first_name" example:"John"`

	// Second name of the user.
	SecondName string `json:"second_name" example:"Doe"`

	// Gender of the user ("Male" or "Female").
	Gender common.Gender `json:"gender" example:"Male"`

	// Birthday in the format "1990-12-31".
	Birthdate string `json:"birthdate" example:"1990-01-01"`

	// Biography of the user.
	Biography string `json:"biography" example:"Software developer"`

	// City of the user.
	City string `json:"city" example:"New York"`
}
