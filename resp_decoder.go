package radish

import (
	"fmt"
	"strconv"
	"strings"
)

var IncompleteRESPError = fmt.Errorf("incomplete resp string")

func Decode(encoded string) (string, error) {

	encodedReader := strings.NewReader(encoded)

	firstChar, err := encodedReader.ReadByte()
	if err != nil {
		return "", IncompleteRESPError
	}

	switch firstChar {
	case '+':
		return decodeSimpleString(encoded)
	case '$':
		return decodeBulkString(encoded)
	default:
		return "", fmt.Errorf("unknown type of resp message provided")
	}
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
	if headerIndex == -1 {
		return "", IncompleteRESPError
	}
	numCharacters, _ := strconv.Atoi(encoded[1:headerIndex])

	if headerIndex+2+numCharacters >= len(encoded) {
		return "", IncompleteRESPError
	}

	if !strings.HasSuffix(encoded, "\r\n") {
		return "", IncompleteRESPError
	}

	message := encoded[headerIndex+2 : headerIndex+numCharacters+2]
	return message, nil
}
