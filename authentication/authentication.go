package authentication

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	values := headers.Get("Authorization")
	if len(values) == 0 {
		return "", errors.New("No api key found")
	}
	heads := strings.Split(values, " ")
	if len(heads) != 2 || heads[0] != "ApiKey" {
		return "", errors.New("Malformed ApiKey Header")
	}
	return heads[1], nil
}
