package account

import (
	"fmt"
	"social-network-service/internal/api/common"
	"social-network-service/internal/model"
	"time"
)

func mapRegisterRequestToCommand(req *RegisterRequest) (*model.RegisterUserCommand, error) {
	birthday, err := time.Parse("2006-01-02", req.Birthdate)

	if err != nil {
		return nil, fmt.Errorf("wrong birthday format: %v", err)
	}

	gender, err := common.MapGenderToModel(req.Gender)

	if err != nil {
		return nil, fmt.Errorf("wrong gender: %v", err)
	}

	registerUser := &model.RegisterUserCommand{
		FirstName:  req.FirstName,
		SecondName: req.SecondName,
		Gender:     gender,
		Birthdate:  birthday,
		Biography:  req.Biography,
		City:       req.City,
		Password:   req.Password,
	}

	return registerUser, nil
}
