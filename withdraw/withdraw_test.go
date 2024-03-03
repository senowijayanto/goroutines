package main

import (
	"fmt"
	"sync"
	"testing"
)

// Account represents a user's bank account
type Account struct {
	sync.Mutex
	ID      int
	Balance int
}

// NewAccount creates a new account with the given ID and initial balance
func NewAccount(id, initialBalance int) *Account {
	return &Account{ID: id, Balance: initialBalance}
}

// Deposit adds the specified amount to the account balance
func (a *Account) Deposit(amount int) {
	a.Lock()
	defer a.Unlock()
	fmt.Printf("Deposit: +%d to Account %d\n", amount, a.ID)
	a.Balance += amount
}

// Withdraw deducts the specified amount from the account balance
func (a *Account) Withdraw(amount int) {
	a.Lock()
	defer a.Unlock()

	if a.Balance >= amount && amount > 0 {
		fmt.Printf("Withdraw: -%d from Account %d\n", amount, a.ID)
		a.Balance -= amount
	} else {
		fmt.Printf("Withdraw: Insufficient funds or invalid amount in Account %d\n", a.ID)
	}
}

// Transfer transfers the specified amount from one account to another
func Transfer(sender, receiver *Account, amount int, wg *sync.WaitGroup) {
	defer wg.Done()
	sender.Withdraw(amount)
	receiver.Deposit(amount)
}

// TestBankTransactions simulates concurrent bank transactions and checks the final balances
func TestBankTransactions(t *testing.T) {
	// Create accounts
	account1 := NewAccount(1, 1000)
	account2 := NewAccount(2, 2000)

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Simulate concurrent transactions with goroutines
	wg.Add(3)
	go Transfer(account1, account2, 200, &wg)
	go Transfer(account2, account1, 100, &wg)
	go Transfer(account1, account2, 300, &wg)

	// Wait for all goroutines to finish
	wg.Wait()

	// Print final balances
	fmt.Printf("Final balance for Account %d: %d\n", account1.ID, account1.Balance)
	fmt.Printf("Final balance for Account %d: %d\n", account2.ID, account2.Balance)

	// Check if the total balance remains the same (consistency check)
	totalBalance := account1.Balance + account2.Balance
	expectedTotalBalance := 1000 + 2000
	if totalBalance != expectedTotalBalance {
		t.Errorf("Inconsistent total balance. Expected: %d, Actual: %d", expectedTotalBalance, totalBalance)
	}
}
