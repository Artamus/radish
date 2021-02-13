package radish

import (
	"fmt"
	"strings"
)

func Decode(encoded string) (string, error) {

	if !strings.HasPrefix(encoded, "+") || !strings.HasSuffix(encoded, "\r\n") {
		return "", fmt.Errorf("incomplete resp string")
	}

	without_prefix := strings.TrimPrefix(encoded, "+")
	without_suffix := strings.TrimSuffix(without_prefix, "\r\n")

	return without_suffix, nil
}
