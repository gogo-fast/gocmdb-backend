package utils

import (
	"github.com/satori/go.uuid"
)

func GenUuid() string {
	return uuid.NewV4().String()
}
