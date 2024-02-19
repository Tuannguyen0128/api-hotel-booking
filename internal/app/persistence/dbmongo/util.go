package dbmongo

import (
	"strings"

	"api-hotel-booking/internal/app/persistence"
)

func castKnownError(e error) error {
	if e == nil {
		return nil
	}

	if e.Error() == "mongo: no documents in result" {
		return persistence.NotFoundError
	}

	if strings.Contains(e.Error(), "duplicate key error") && strings.Contains(e.Error(), fieldUniqueTag) {
		return persistence.DuplicateUniqueTagError
	}

	if strings.Contains(e.Error(), "connect: connection refused") {
		return persistence.DatabaseConnectError
	}
	return e
}
