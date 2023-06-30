package bodyhasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func CalculateHash(data []byte, hashKey string) (string, error) {
	h := hmac.New(sha256.New, []byte(hashKey))
	_, err := h.Write(data)
	if err != nil {
		return "", fmt.Errorf("didn't come up with %w", err)
	}
	hs := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return hs, nil
}
