package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type UserBalance struct {
	sync.Mutex
	Name    string
	Balance int
}

func (ub *UserBalance) Lock() {
	ub.Mutex.Lock()
}

func (ub *UserBalance) Unlock() {
	ub.Mutex.Unlock()
}

func (ub *UserBalance) ChangeBalance(amount int) {
	ub.Balance += amount
}

var lockOrder = sync.Mutex{}

func Transfer(user1 *UserBalance, user2 *UserBalance, amount int) {
	lockOrder.Lock()
	defer lockOrder.Unlock()

	user1.Lock()
	fmt.Println("Lock user1", user1.Name)
	user1.ChangeBalance(-amount)

	time.Sleep(2 * time.Second)

	user2.Lock()
	fmt.Println("Lock user2", user2.Name)
	user2.ChangeBalance(amount)

	time.Sleep(2 * time.Second)

	user1.Unlock()
	user2.Unlock()
}

func TestDeadLock(t *testing.T) {
	user1 := &UserBalance{Name: "John", Balance: 1000000}

	user2 := &UserBalance{Name: "Doe", Balance: 1000000}

	go Transfer(user1, user2, 100000)
	go Transfer(user2, user1, 200000)

	time.Sleep(10 * time.Second)

	fmt.Println("user1 balance:", user1.Balance)
	fmt.Println("user2 balance:", user2.Balance)
}
