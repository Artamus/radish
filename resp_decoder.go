package radish

import (
	"fmt"
	"strconv"
	"strings"
)

var IncompleteRESPError = fmt.Errorf("incomplete resp string")

func Decode(encoded string) (string, error) {

	if encoded[0] == '$' {
		return decodeBulkString(encoded)
	}

	return decodeSimpleString(encoded)
}

func decodeSimpleString(encoded string) (string, error) {
	if !strings.HasPrefix(encoded, "+") || !strings.HasSuffix(encoded, "\r\n") {
		return "", IncompleteRESPError
	}

	without_prefix := strings.TrimPrefix(encoded, "+")
	without_suffix := strings.TrimSuffix(without_prefix, "\r\n")

	return without_suffix, nil
}

func decodeBulkString(encoded string) (string, error) {
	headerIndex := strings.Index(encoded, "\r\n")
	numCharacters, _ := strconv.Atoi(encoded[1:headerIndex])

	message := encoded[headerIndex+2 : headerIndex+numCharacters+2]
	return message, nil
}
