package radish

import (
	"fmt"
	"strings"
)

var IncompleteRESPError = fmt.Errorf("incomplete resp string")

func Decode(encoded string) (string, error) {

	if !strings.HasPrefix(encoded, "+") || !strings.HasSuffix(encoded, "\r\n") {
		return "", IncompleteRESPError
	}

	without_prefix := strings.TrimPrefix(encoded, "+")
	without_suffix := strings.TrimSuffix(without_prefix, "\r\n")

	return without_suffix, nil
}
