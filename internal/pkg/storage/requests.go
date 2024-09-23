package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

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

func (s *Model) CreateDrive(ctx context.Context, inputUser Drive) (drive Drive, err error) {
	const query = `insert into drives (
					driver_id,
                   	time,
                    direction,
                    car,
                    price,
                    slots, 
                    created_at
				) 
			values ($1, $2, $3, $4, $5, $6, $7) returning *`

	if err = s.db.GetContext(
		ctx,
		&drive,
		query,
		inputUser.DriverID,
		inputUser.Time,
		inputUser.Direction,
		inputUser.Car,
		inputUser.Price,
		inputUser.Slots,
		time.Now().UTC(),
	); err != nil {
		return Drive{}, trace.FuncNameWithError(err)
	}
	return drive, nil
}

func (s *Model) GetLast10Drives(ctx context.Context, driverID uuid.UUID) (drives []Drive, err error) {
	const query = `select * from drives where driver_id = $1 order by time limit 10`

	if err = s.db.SelectContext(
		ctx,
		&drives,
		query,
		driverID,
	); err != nil {
		return nil, trace.FuncNameWithError(err)
	}
	return drives, nil
}

func (s *Model) GetLastCars(ctx context.Context, driverID uuid.UUID) (cars []int64, err error) {
	const query = `select distinct car from drives where driver_id = $1 limit 10`

	if err = s.db.SelectContext(
		ctx,
		&cars,
		query,
		driverID,
	); err != nil {
		return nil, trace.FuncNameWithError(err)
	}
	return cars, nil
}

func (s *Model) CreateOrResolveCarID(ctx context.Context, name string) (carID int64, err error) {
	const queryToCreate = `insert into cars (name) values ($1) on conflict (name) do nothing returning *`
	const queryToGet = `select id from cars where name = $1`

	if _, err = s.db.ExecContext(ctx, queryToCreate, name); err != nil {
		return 0, trace.FuncNameWithErrorMsg(err, "creating car")
	}
	if err = s.db.GetContext(ctx, &carID, queryToGet, name); err != nil {
		return 0, trace.FuncNameWithErrorMsg(err, "getting car")
	}
	return carID, nil
}

func (s *Model) GetCarsByID(ctx context.Context, carIDs []int64) (cars []Car, err error) {
	const query = `select * from cars where id = any($1)`

	if err = s.db.SelectContext(ctx, &cars, query, pq.Array(carIDs)); err != nil {
		return nil, trace.FuncNameWithErrorMsg(err, "getting cars")
	}
	return cars, nil
}

func (s *Model) GetCarByID(ctx context.Context, id int) (car Car, err error) {
	const query = `select * from cars where id = $1`

	if err = s.db.GetContext(ctx, &car, query, id); err != nil {
		return Car{}, trace.FuncNameWithErrorMsg(err, "getting cars")
	}
	return car, nil
}
