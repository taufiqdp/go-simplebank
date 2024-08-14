package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)

	if err := fn(q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var TxKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// txName := ctx.Value(TxKey).(string)

		if result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		}); err != nil {
			return err
		}

		if result.FromAccount, err = q.GetAccountForUpdate(ctx, arg.FromAccountID); err != nil {
			return err
		}

		if result.ToAccount, err = q.GetAccountForUpdate(ctx, arg.ToAccountID); err != nil {
			return err
		}

		if result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		}); err != nil {
			return err
		}

		if result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		}); err != nil {
			return err
		}

		account1, err := q.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		account2, err := q.GetAccount(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		if result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      account1.ID,
			Balance: account1.Balance - arg.Amount,
		}); err != nil {
			return err
		}

		if result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      account2.ID,
			Balance: account2.Balance + arg.Amount,
		}); err != nil {
			return err
		}

		return nil
	})

	return result, err
}
