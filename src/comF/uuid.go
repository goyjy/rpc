package comF

import (
	"github.com/google/uuid"
	"math/rand"
)

func GetUUID() string {
	uid := uuid.New()
	return uid.String()
}

func GetWorkId() uint32 {
	return rand.Uint32()
}
