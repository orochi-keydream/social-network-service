package account

import "social-network-service/internal/api/common"

type LoginRequest struct {
	// ID of the post in UUIDv4 format.
	UserId string `json:"user_id" binding:"required" example:"bbeb7da8-6d75-4419-9d94-91ec52bc506c"`

	// Password to the account.
	Password string `json:"password" binding:"required" example:"123456"`
}

type LoginResponse struct {
	// Authorization token to the account (JWT).
	Token string `json:"token"`
}

type RegisterRequest struct {
	// First name of the user.
	FirstName string `json:"first_name" binding:"required" example:"John"`

	// Second name of the user.
	SecondName string `json:"second_name" binding:"required" example:"Doe"`

	// Gender of the user ("Male" or "Female").
	Gender common.Gender `json:"gender" binding:"required" example:"Male"`

	// Birthday in the format "1990-12-31".
	Birthdate string `json:"birthdate" binding:"required" example:"1990-01-01"`

	// Biography of the user.
	Biography string `json:"biography" binding:"required" example:"Software developer"`

	// City of the user.
	City string `json:"city" binding:"required" example:"New York"`

	// Password to the account.
	Password string `json:"password" binding:"required" example:"123456"`
}

type RegisterResponse struct {
	// ID of the user in UUIDv4 format.
	UserId string `json:"user_id" example:"bbeb7da8-6d75-4419-9d94-91ec52bc506c"`
}
