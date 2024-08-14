package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Printf("Account %d balace: %d\n", account1.ID, account1.Balance)
	fmt.Printf("Account %d balace: %d\n", account2.ID, account2.Balance)

	n := 5
	amount := int64(10)

	// Run n concurrent transfer transactions
	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), CreateTransferParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result

		}()
	}

	for i := 0; i < n; i++ {
		fmt.Printf("\nIteration: %d\n", i + 1)
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)
		require.NotZero(t, result.Transfer.ID)

		require.Equal(t, account1.ID, result.FromAccount.ID)
		require.Equal(t, account2.ID, result.ToAccount.ID)

		require.Equal(t, -amount, result.FromEntry.Amount)
		require.Equal(t, amount, result.ToEntry.Amount)

		tf, err := store.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)
		require.NotEmpty(t, tf)

		tent, err := store.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)
		require.NotEmpty(t, tent)

		fent, err := store.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)
		require.NotEmpty(t, fent)

		fmt.Println("Transfer: ", tf.Amount)
		fmt.Println("From Entry: ", fent.Amount)
		fmt.Println("To Entry: ", tent.Amount)

		dif1 := account1.Balance - result.FromAccount.Balance
		dif2 := result.ToAccount.Balance - account2.Balance
		
		fmt.Println("Account 1 balace: ", result.FromAccount.Balance)
		fmt.Println("Account 2 balace: ", result.ToAccount.Balance)
		fmt.Println("Dif1: ", dif1)
		fmt.Println("Dif2: ", dif2)

		require.Equal(t, dif1, dif2)
		require.True(t, dif1%amount == 0)
		require.True(t, dif2%amount == 0)
		require.True(t, dif1 > 0 && dif1 <= int64(n) * amount)
		require.True(t, dif2 > 0 && dif2 <= int64(n) * amount)
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}
