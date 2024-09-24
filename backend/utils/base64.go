package utils

import (
	"encoding/base64"
	"strings"
)

func convertToBase64StdEncoding(base64URLEncoded string) string {
	base64StdEncoded := base64URLEncoded
	base64StdEncoded = strings.Replace(base64StdEncoded, "-", "+", -1)
	base64StdEncoded = strings.Replace(base64StdEncoded, "_", "/", -1)

	if len(base64StdEncoded)%4 != 0 {
		base64StdEncoded += strings.Repeat("=", 4-len(base64StdEncoded)%4)
	}

	return base64StdEncoded
}

func DecodeToBase64StdEncoding(src string) ([]byte, error) {
	base64StdEncoded := convertToBase64StdEncoding(src)
	base64StdDecoded, err := base64.StdEncoding.DecodeString(base64StdEncoded)

	return []byte(string(base64StdDecoded)), err
}
