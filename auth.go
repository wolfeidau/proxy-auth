package auth

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

const tokenLength = 48

func generateState() (string, error) {
	token := make([]byte, tokenLength)
	if _, err := io.ReadFull(rand.Reader, token); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(token), nil
}
