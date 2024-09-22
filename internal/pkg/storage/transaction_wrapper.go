package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PGTransactionFn func(tx *sqlx.Tx) error

func (s *Model) HandlePGTransaction(pgTx PGTransactionFn) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		var errTxn error
		if err != nil {
			errTxn = tx.Rollback()
		} else {
			errTxn = tx.Commit()
		}
		if errTxn != nil {
			err = fmt.Errorf("failed executing transaction: %w", err)
		}
	}()

	err = pgTx(tx)
	return err
}
