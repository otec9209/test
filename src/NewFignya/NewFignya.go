package main

import "fmt"

type BankAccount struct {
	Owner        string
	Balance      float64
	Transactions []string
}

func (b *BankAccount) Deposit(amount float64) {
	oldBalance := b.Balance
	b.Balance += amount
	if b.Balance != oldBalance {
		record := fmt.Sprintf("Снятие: -%.2f | Баланс: %.2f", amount, b.Balance)
		b.Transactions = append(b.Transactions, record)
	}
}

func (b *BankAccount) Withdraw(amount float64) {
	if b.Balance < amount {
		fmt.Println("Недостаточно средств")
		return
	}
	oldBalance := b.Balance
	b.Balance -= amount
	if b.Balance != oldBalance {
		record := fmt.Sprintf("Снятие: -%.2f | Баланс: %.2f", amount, b.Balance)
		b.Transactions = append(b.Transactions, record)

	}

}

func (b *BankAccount) PrintBalance() {
	fmt.Printf("Баланс владельца %s: %.2f\n", b.Owner, b.Balance)

}
func (b *BankAccount) PrintHystory() {
	fmt.Println("История операций")
	for _, t := range b.Transactions {
		fmt.Println("/", t)
	}
}

func main() {
	acc := BankAccount{Owner: "Андрей", Balance: 1000}
	acc.PrintBalance()

	acc.Deposit(500)
	acc.PrintBalance()

	acc.Withdraw(2000) // Недостаточно средств
	acc.Withdraw(300)
	acc.PrintBalance()
	acc.PrintHystory()
}
