package generated_db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferMoneyTxn(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t) // from account
	account2 := createRandomAccount(t) // to account

	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// Transfer money from account1 to account2
	// Do 5 concurrent transfers of 10 dollars each
	// to test for database consistency, deadlocks, lost updates, etc.

	n := 10
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferMoneyTxnResult)

	go func() {
		for i := 0; i < n; i++ {
			result, err := store.TransferMoneyTxn(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}
	}()

	// Check for errors

	for i := 0; i < n; i++ {

		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer (no concurrency issues as inserting a new row
		// in the transfers table is an atomic operation and doesn't depend on the
		// other transfers)
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// check entries (no concurrency issues as inserting a new row)
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check balances (here is where we can have concurrency issues)
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		require.Equal(t, account1.Balance-int64(i+1)*amount, fromAccount.Balance)
		require.Equal(t, account2.Balance+int64(i+1)*amount, toAccount.Balance)

	}

	// Check account balances at the end
	// account1 should have 50 dollars less and account2 should have 50 dollars more
	requiredBalance1 := account1.Balance - int64(n)*amount
	requiredBalance2 := account2.Balance + int64(n)*amount

	account1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err = store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> Required Balances:", requiredBalance1, requiredBalance2)
	fmt.Println(">> Final balances:", account1.Balance, account2.Balance)

	require.Equal(t, requiredBalance1, account1.Balance)
	require.Equal(t, requiredBalance2, account2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	testStore := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := testStore.TransferMoneyTxn(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
