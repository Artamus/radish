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

func Decode(encoded string) (string, error) {

	encodedReader := strings.NewReader(encoded)

	firstChar, err := encodedReader.ReadByte()
	if err != nil {
		return "", IncompleteRESPError
	}

	switch firstChar {
	case '+':
		return decodeSimpleString(encodedReader)
	case '$':
		return decodeBulkString(encodedReader)
	default:
		return "", fmt.Errorf("unknown type of resp message provided")
	}
}

func decodeSimpleString(rdr io.Reader) (string, error) {
	bufRdr := bufio.NewReader(rdr)

	contents, err := readToCRLF(bufRdr)
	if err != nil {
		return "", IncompleteRESPError
	}

	return contents, nil
}

func decodeBulkString(rdr io.Reader) (string, error) {
	bufRdr := bufio.NewReader(rdr)
	contents, err := readToCRLF(bufRdr)

	contentLength, err := strconv.Atoi(contents)
	if err != nil {
		return "", IncompleteRESPError
	}

	content, err := ioutil.ReadAll(io.LimitReader(bufRdr, int64(contentLength)))
	if len(content) != contentLength {
		return "", IncompleteRESPError
	}

	controlBytes, err := ioutil.ReadAll(bufRdr)
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
