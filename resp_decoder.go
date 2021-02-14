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
	contents, err := readToCRLF(rdr)
	contentLength, err := strconv.Atoi(contents)
	if err != nil {
		return "", IncompleteRESPError
	}

	content, err := ioutil.ReadAll(io.LimitReader(rdr, int64(contentLength)))
	if len(content) != contentLength {
		return "", IncompleteRESPError
	}

	controlBytes, err := ioutil.ReadAll(rdr)
	if err != nil || string(controlBytes) != "\r\n" {
		return "", IncompleteRESPError
	}

	return string(content), nil
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
