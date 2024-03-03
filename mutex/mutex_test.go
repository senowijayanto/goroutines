package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// UserBalance represents a user with a balance and provides methods for locking and changing the balance
type UserBalance struct {
	sync.Mutex
	Name    string
	Balance int
}

// Lock acquires the lock for the user's balance
func (ub *UserBalance) Lock() {
	ub.Mutex.Lock()
}

// Unlock releases the lock for the user's balance
func (ub *UserBalance) Unlock() {
	ub.Mutex.Unlock()
}

// ChangeBalance updates the user's balance by the specified amount
func (ub *UserBalance) ChangeBalance(amount int) {
	ub.Balance += amount
}

var lockOrder = sync.Mutex{} // Global mutex to establish a lock order

// Transfer simulates a fund transfer between two users, avoiding potential deadlocks by using a global lock order
func Transfer(user1 *UserBalance, user2 *UserBalance, amount int) {
	lockOrder.Lock()         // Acquire the global lock order
	defer lockOrder.Unlock() // Ensure the global lock order is released even if an error occurs

	user1.Lock() // Acquire lock for user1's balance
	fmt.Println("Lock user1", user1.Name)
	user1.ChangeBalance(-amount) // Withdraw from user1's balance

	time.Sleep(2 * time.Second) // Simulate processing time

	user2.Lock() // Acquire lock for user2's balance
	fmt.Println("Lock user2", user2.Name)
	user2.ChangeBalance(amount) // Deposit to user2's balance

	time.Sleep(2 * time.Second) // Simulate processing time

	user1.Unlock() // Release lock for user1's balance
	user2.Unlock() // Release lock for user2's balance
}

// TestDeadLock tests the Transfer function with two users and checks their final balances
func TestDeadLock(t *testing.T) {
	user1 := &UserBalance{Name: "John", Balance: 1000000}
	user2 := &UserBalance{Name: "Doe", Balance: 1000000}

	go Transfer(user1, user2, 100000) // Initiate a fund transfer from user1 to user2
	go Transfer(user2, user1, 200000) // Initiate a fund transfer from user2 to user1

	time.Sleep(10 * time.Second) // Wait for goroutines to finish

	// Print final balances
	fmt.Println("user1 balance:", user1.Balance)
	fmt.Println("user2 balance:", user2.Balance)
}
