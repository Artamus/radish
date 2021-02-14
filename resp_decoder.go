package radish

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

var IncompleteRESPError = fmt.Errorf("incomplete resp string")

// Decode returns the decoded body as either string or a slice of strings
func Decode(encoded string) (interface{}, error) {
	rdr := bufio.NewReader(strings.NewReader(encoded))
	return decode(rdr)
}

func decode(rdr *bufio.Reader) (interface{}, error) {
	firstChar, err := rdr.ReadByte()
	if err != nil {
		return "", IncompleteRESPError
	}

	switch firstChar {
	case '+':
		return decodeSimpleString(rdr)
	case '$':
		return decodeBulkString(rdr)
	case '*':
		return decodeArray(rdr)
	default:
		return "", fmt.Errorf("unknown type of resp message provided")
	}
}

func decodeSimpleString(rdr *bufio.Reader) (string, error) {
	contents, err := readToCRLF(rdr)
	if err != nil {
		return "", IncompleteRESPError
	}

	return contents, nil
}

func decodeBulkString(rdr *bufio.Reader) (string, error) {
	length, err := readToCRLF(rdr)
	contentLength, err := strconv.Atoi(length)
	if err != nil {
		return "", IncompleteRESPError
	}

	content, err := ioutil.ReadAll(io.LimitReader(rdr, int64(contentLength)))
	if len(content) != contentLength {
		return "", IncompleteRESPError
	}

	controlBytes, err := ioutil.ReadAll(io.LimitReader(rdr, 2))
	if err != nil || string(controlBytes) != "\r\n" {
		return "", IncompleteRESPError
	}

	return string(content), nil
}

func decodeArray(rdr *bufio.Reader) ([]interface{}, error) {
	content, err := readToCRLF(rdr)
	numItems, err := strconv.Atoi(content)
	if err != nil {
		return make([]interface{}, 0), IncompleteRESPError
	}

	results := make([]interface{}, numItems)
	for i := 0; i < numItems; i++ {
		decoded, err := decode(rdr)

		if err != nil {
			return make([]interface{}, 0), IncompleteRESPError
		}
		results[i] = decoded
	}

	return results, nil
}

func readToCRLF(rdr *bufio.Reader) (string, error) {
	content, err := rdr.ReadString('\r')
	if err != nil {
		return "", fmt.Errorf("")
	}

	controlByte, err := rdr.ReadByte()
	if err != nil || controlByte != '\n' {
		return "", fmt.Errorf("")
	}

	return content[:len(content)-1], nil
}
