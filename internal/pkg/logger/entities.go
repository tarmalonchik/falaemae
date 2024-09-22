package logger

import (
	"encoding/json"
	"io"
)

const (
	chatNotFound = "Bad Request: chat not found"
)

type errResp struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func (e *errResp) ReadFromReader(reader io.Reader) {
	data, _ := io.ReadAll(reader)
	_ = json.Unmarshal(data, e)
}

func (e *errResp) isNotFound() bool {
	if e == nil {
		return false
	}
	return e.Description == chatNotFound
}
