package main

import "sync"

// Dish представляет собой структуру для хранения информации о блюде
type Dish struct {
	Name  string  // Название блюда
	Price float64 // Цена блюда
}

// Menu представляет собой структуру для хранения списка блюд
type Orders struct {
	Dishes []Dish // Список блюд
}
type kitchen struct {
	mu     sync.Mutex
	orders []Orders
	isbusy bool
}
