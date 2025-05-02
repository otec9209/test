package main

import (
	"errors"
	"fmt"
)

type nabor interface {
	snat(amount float64) error
	vnesti(amount float64)
	balance() float64
}
type chelovek struct {
	Name    string
	Balance float64
}

func (c *chelovek) snat(amount float64) error {
	if amount <= 0 {
		return errors.New("сумма должна быть больше нуля")
	}
	if amount > c.Balance {
		return errors.New("недостаточно средств")
	}
	c.Balance -= amount
	return nil
}
func (c *chelovek) vnesti(amount float64) {
	if amount > 0 {
		c.Balance += amount
	}
}
func (c *chelovek) balance() float64 {
	return c.Balance
}

func main() {

	c := &chelovek{Name: "alexey", Balance: 2000}
	err := c.snat(3000)
	fmt.Println(err)

}
