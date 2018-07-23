package util

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
)

func SHA1(w interface{}) (string, error) {
	data, err := json.Marshal(w)
	if err != nil {
		return "", err
	}
	h := sha1.New()
	_, err = h.Write(data)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
