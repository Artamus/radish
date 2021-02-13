package radish

import "strings"

func Decode(encoded string) (string, error) {

	without_prefix := strings.TrimPrefix(encoded, "+")
	without_suffix := strings.TrimSuffix(without_prefix, "\r\n")

	return without_suffix, nil
}
