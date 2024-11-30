package utils

import "github.com/google/uuid"

func ValidateUUID(id *string) error {
	_, err := uuid.Parse(*id)

	if err != nil {
		return err
	}

	return nil
}
