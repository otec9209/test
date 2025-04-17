package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// список доступных блюд
var menu = []string{"пицца", "суп", "стейк", "салат", "суши"}

// waiter принимает канал заказов
func waiter(orders chan string) {
	for {
		time.Sleep(time.Duration(rand.Intn(6)+5) * time.Second)

		dish := menu[rand.Intn(len(menu))]

		fmt.Println("заказ получен официантом и отправлен на обработку на кухню", dish)

		orders <- dish //отправляем рандомное блюдо по каналу в заказ

	}
}

// Сначала я создал список блюд, которые указаны в задаче через Слайс
// затем я создал функцию которая принимает  заказ через канал
// в функции пишем, что заказ создается каждые 5-10 секунд
// затем создаем случайное заказанное блюдо
// отправляем блюдо по каналу в заказ

func kitchen(orders <-chan string, readmils chan<- string, mu *sync.Mutex, deliveredCount *int) {

	for {
		mu.Lock()
		if *deliveredCount >= 10 {
			mu.Unlock()
			return
		}
		mu.Unlock()

		select {
		case dish := <-orders:
			fmt.Println("начата готовка блюда:", dish)
			time.Sleep(time.Duration(rand.Intn(11)+5) * time.Second)
			fmt.Println("Блюдо готово:", dish)
			readmils <- dish

		default:
			time.Sleep(500 * time.Millisecond)

		}
	}
}

//Создаем функцию кухни, которая принимает блюда от официантов
//затем пишем мютекс на огранчиение блюд до 10
//пишем селект на готовку блюда от 5 до 10 секунд
//передаем по каналу в готовые блюда

func courier(readmils <-chan string, mu *sync.Mutex, deliveredCount *int) {
	for {
		mu.Lock()
		if *deliveredCount >= 10 {
			mu.Unlock()
			return
		}
		mu.Unlock()

		time.Sleep(30 * time.Second)
		var batch []string

	collektlop:
		for {
			select {
			case dish := <-readmils:
				batch = append(batch, dish)
			default:
				break collektlop
			}
		}
		if len(batch) > 0 {
			mu.Lock()
			*deliveredCount += len(batch)
			current := *deliveredCount
			mu.Unlock()
			fmt.Printf("[Курьер] Забрал %d блюд:\n", len(batch))
			for _, d := range batch {
				fmt.Println("  -", d)
			}
			fmt.Printf("[Курьер] Доставка #%d завершена\n", current)
		}
	}
}

func main() {
	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Создание каналов и мьютекса
	orders := make(chan string, 10)
	readmils := make(chan string, 10)
	mu := &sync.Mutex{}
	deliveredCount := 0

	// Запуск горутин
	go waiter(orders)
	go kitchen(orders, readmils, mu, &deliveredCount)
	go courier(readmils, mu, &deliveredCount)

	// Ожидание завершения доставки 10 блюд
	for {
		mu.Lock()
		if deliveredCount >= 10 {
			mu.Unlock()
			fmt.Println("Работа завершена, все заказы выполнены")
			return
		}
		mu.Unlock()
		time.Sleep(1 * time.Second)
	}
}
