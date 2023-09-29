package sirabbitmq

import (
	"github.com/google/uuid"
)

func generateId() string {

	tmpId := uuid.New().String()

	return tmpId

	// offset := strings.Index(tmpId, "-")
	// if offset > -1 {
	// 	tmpId = tmpId[:offset]
	// }
	// return tmpId
}
