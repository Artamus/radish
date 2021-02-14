package radish

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

// ErrIncompleteRESP is returned when the encoded string fed to Decode is not valid
var ErrIncompleteRESP = fmt.Errorf("incomplete resp string")

var emptyResponse = make([]string, 0)

// Decode returns the decoded body as either string or a slice of strings
func Decode(encoded string) ([]string, error) {
	rdr := bufio.NewReader(strings.NewReader(encoded))
	return decode(rdr)
}

func decode(rdr *bufio.Reader) ([]string, error) {
	firstChar, err := rdr.ReadByte()
	if err != nil {
		return emptyResponse, ErrIncompleteRESP
	}

	switch firstChar {
	case '+':
		return decodeSimpleString(rdr)
	case '$':
		return decodeBulkString(rdr)
	case '*':
		return decodeArray(rdr)
	default:
		return emptyResponse, fmt.Errorf("unknown type of resp message provided")
	}
}

func decodeSimpleString(rdr *bufio.Reader) ([]string, error) {
	contents, err := readToCRLF(rdr)
	if err != nil {
		return emptyResponse, ErrIncompleteRESP
	}

	return []string{contents}, nil
}

func decodeBulkString(rdr *bufio.Reader) ([]string, error) {
	length, err := readToCRLF(rdr)
	contentLength, err := strconv.Atoi(length)
	if err != nil {
		return emptyResponse, ErrIncompleteRESP
	}

	content, err := ioutil.ReadAll(io.LimitReader(rdr, int64(contentLength)))
	if len(content) != contentLength {
		return emptyResponse, ErrIncompleteRESP
	}

	controlBytes, err := ioutil.ReadAll(io.LimitReader(rdr, 2))
	if err != nil || string(controlBytes) != "\r\n" {
		return emptyResponse, ErrIncompleteRESP
	}

	return []string{string(content)}, nil
}

func decodeArray(rdr *bufio.Reader) ([]string, error) {
	content, err := readToCRLF(rdr)
	numItems, err := strconv.Atoi(content)
	if err != nil {
		return emptyResponse, ErrIncompleteRESP
	}

	results := make([]string, 0)
	for i := 0; i < numItems; i++ {
		decoded, err := decode(rdr)

		if err != nil {
			return emptyResponse, ErrIncompleteRESP
		}
		results = append(results, decoded...)
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
