package radish

import "strings"

func Decode(encoded string) string {

	without_prefix := strings.TrimPrefix(encoded, "+")
	without_suffix := strings.TrimSuffix(without_prefix, "\r\n")

	return without_suffix
}
