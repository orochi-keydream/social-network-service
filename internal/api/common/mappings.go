package common

import (
	"fmt"
	"social-network-service/internal/model"
)

func MapGenderToModel(gender Gender) (model.Gender, error) {
	switch gender {
	case GenderMale:
		return model.GenderMale, nil
	case GenderFemale:
		return model.GenderFemale, nil
	default:
		return 0, fmt.Errorf("unsupported gender %v", gender)
	}
}

func MapGenderFromModel(gender model.Gender) (Gender, error) {
	switch gender {
	case model.GenderMale:
		return GenderMale, nil
	case model.GenderFemale:
		return GenderFemale, nil
	default:
		return "", fmt.Errorf("unsupported gender %v", gender)
	}
}
