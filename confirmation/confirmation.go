package confirmation

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

type UseCase interface {
}

func GenerateConfirmationToken() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return strings.TrimRight(base64.URLEncoding.EncodeToString(bytes), "=")
}
