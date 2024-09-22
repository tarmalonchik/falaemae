package telegram

import (
	"database/sql"

	"github.com/vkidmode/telegram_tree"
)

type SendDocumentRequest struct {
	ChatID   sql.NullInt64
	FileName string
	File     []byte
	FilePath sql.NullString
}

type sendDocResponseBody struct {
	Ok     bool   `json:"ok"`
	Result result `json:"result"`
}

type result struct {
	Document document `json:"document"`
}

type document struct {
	FileID string `json:"file_id"`
}

type NewKeyboardRequest struct {
	NextNodes    []telegram_tree.Node
	Meta         metaInfo
	HideBar      bool
	CallbackBack string
	CallbackSkip string
}

type metaInfo interface {
	GetCallback() string
}
