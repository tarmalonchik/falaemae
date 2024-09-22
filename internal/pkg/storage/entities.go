package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/tarmalonchik/falaemae/internal/pkg/telegram"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID           uuid.UUID      `db:"id"`
	TelegramNick sql.NullString `db:"telegram_user"`
	ChatID       sql.NullInt64  `db:"chat_id"`
	CreatedAt    time.Time      `db:"created_at"`
}

func (u *User) IsAdmin() bool {
	if u == nil {
		return false
	}
	if u.ChatID.Valid {
		return telegram.KeyBoardItem.IsAdmin(u.ChatID.Int64)
	}
	if u.TelegramNick.Valid {
		return telegram.KeyBoardItem.IsAdmin(u.TelegramNick.String)
	}
	return false
}

func (u *User) GetUserString() string {
	var userString string

	if u.TelegramNick.Valid {
		if u.TelegramNick.String != "" {
			userString += fmt.Sprintf("[telegram-nick: @%s]", u.TelegramNick.String)
		}
	}

	if u.ChatID.Valid {
		if u.ChatID.Int64 != 0 {
			userString += fmt.Sprintf("[chatID: %d]", u.ChatID.Int64)
		}
	}
	return userString
}
