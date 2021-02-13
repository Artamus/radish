package radish

import (
	"fmt"
	"strings"
)

func Decode(encoded string) (string, error) {

	without_prefix := strings.TrimPrefix(encoded, "+")
	without_suffix := strings.TrimSuffix(without_prefix, "\r\n")

	if len(without_suffix) == 0 {
		return "", fmt.Errorf("simple string message '%s' is incomplete", encoded)
	}

	return without_suffix, nil
}
