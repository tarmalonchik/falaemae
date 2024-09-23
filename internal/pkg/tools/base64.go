package tools

import (
	"encoding/base64"
)

func ToBase64(message []byte) (out string) {
	base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(message)))
	base64.StdEncoding.Encode(base64Text, message)
	return string(base64Text)
}

func ParseBase64(in string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(in)
}
