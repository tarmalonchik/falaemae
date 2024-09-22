package storage

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/tarmalonchik/falaemae/internal/pkg/postgresql"
	"github.com/tarmalonchik/falaemae/internal/pkg/trace"
)

type Model struct {
	db *postgresql.SQLXClient
}

func NewModel(db *postgresql.SQLXClient) *Model {
	return &Model{db: db}
}

func (s *Model) UpdateTelegramNick(ctx context.Context, chatID int64, nick *string) error {
	const query = `update users set telegram_user = $1 where chat_id = $2`
	if _, err := s.db.ExecContext(ctx, query, nick, chatID); err != nil {
		return trace.FuncNameWithErrorMsg(err, "update nick")
	}
	return nil
}

func (s *Model) GetUserByChatID(ctx context.Context, chatID int64) (user User, err error) {
	const query = `select * from users where chat_id = $1`
	if err = s.db.GetContext(ctx, &user, query, chatID); err != nil {
		if postgresql.IsNotFound(err) {
			return User{}, ErrUserNotFound
		}
		return User{}, trace.FuncNameWithError(err)
	}
	return user, nil
}

func (s *Model) InitUser(ctx context.Context, inputUser User) (user User, err error) {
	const query = `insert into users (
					telegram_user,
					chat_id,
                   	created_at
				) 
			values ($1, $2, $3) on conflict(chat_id) do nothing returning *`

	if err = s.db.GetContext(
		ctx,
		&user,
		query,
		inputUser.TelegramNick,
		inputUser.ChatID,
		time.Now().UTC(),
	); err != nil {
		return User{}, trace.FuncNameWithError(err)
	}
	return user, nil
}

func (s *Model) GetUser(ctx context.Context, userID uuid.UUID) (user User, err error) {
	const query = `select * from users where id = $1`

	if err = s.db.GetContext(ctx, &user, query, userID); err != nil {
		if postgresql.IsNotFound(err) {
			return User{}, ErrUserNotFound
		}
		return User{}, trace.FuncNameWithErrorMsg(err, "getting user")
	}
	return user, nil
}

func (s *Model) GetUserPayload(ctx context.Context, userID uuid.UUID, orderID *uuid.UUID) (string, error) {
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		return "", trace.FuncNameWithErrorMsg(err, "getting user")
	}

	return "\nДанные пользователя: " + user.GetUserString(), nil
}
