package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Интервал очистки кольцевого буфера
const bufferDrainInterval time.Duration = 30 * time.Second

// Размер кольцевого буфера
const bufferSize int = 10

// RingIntBuffer - кольцевой буфер целых чисел
type RingIntBuffer struct {
	array []int      // более низкоуровневое хранилище нашего буфера
	pos   int        // текущая позиция кольцевого буфера
	size  int        // общий размер буфера
	m     sync.Mutex // мьютекс для потокобезопасного доступа к буферу
}

// NewRingIntBuffer - создание нового буфера целых чисел
func NewRingIntBuffer(size int) *RingIntBuffer {
	return &RingIntBuffer{make([]int, size), -1, size, sync.Mutex{}}
}

// Push - добавление нового элемента в конец буфера
// При попытке добавления нового элемента в заполненный буфер
// самое старое значение затирается
func (r *RingIntBuffer) Push(el int) {
	r.m.Lock()
	defer r.m.Unlock()
	if r.pos == r.size-1 {
		// Сдвигаем все элементы буфера на одну позицию в сторону начала
		for i := 1; i <= r.size-1; i++ {
			r.array[i-1] = r.array[i]
		}
		r.array[r.pos] = el
	} else {
		r.pos++
		r.array[r.pos] = el
	}
}

// Get - получение всех элементов буфера и его последующая очистка
func (r *RingIntBuffer) Get() []int {
	if r.pos < 0 {
		return nil
	}
	r.m.Lock()
	defer r.m.Unlock()
	var output []int = r.array[:r.pos+1]
	// Виртуальная очистка нашего буфера
	r.pos = -1
	return output
}

// StageInt - Стадия конвейера, обрабатывающая целые числа
type StageInt func(<-chan bool, <-chan int) <-chan int

// PipeLineInt - Пайплайн обработки целых чисел
type PipeLineInt struct {
	stages []StageInt
	done   <-chan bool
}

// NewPipelineInt - Создание пайплайна обработки целых чисел
func NewPipelineInt(done <-chan bool, stages ...StageInt) *PipeLineInt {
	return &PipeLineInt{done: done, stages: stages}
}

// Run - Запуск пайплайна обработки целых чисел
// source - источник данных для конвейера
func (p *PipeLineInt) Run(source <-chan int) <-chan int {
	var c <-chan int = source
	for index := range p.stages {
		c = p.runStageInt(p.stages[index], c)
	}
	return c
}

// runStageInt - запуск отдельной стадии конвейера
func (p *PipeLineInt) runStageInt(stage StageInt, sourceChan <-chan int) <-chan int {
	return stage(p.done, sourceChan)
}

// dataSource - источник данных
func dataSource() (<-chan int, <-chan bool) {
	c := make(chan int)
	done := make(chan bool)
	go func() {
		defer close(done)
		scanner := bufio.NewScanner(os.Stdin)
		for {
			scanner.Scan()
			data := scanner.Text()
			if strings.EqualFold(data, "exit") {
				fmt.Println("Программа завершила работу!")
				return
			}
			i, err := strconv.Atoi(data)
			if err != nil {
				fmt.Println("Программа обрабатывает только целые числа!")
				continue
			}
			c <- i                             // Отправка одного числа
			time.Sleep(100 * time.Millisecond) // Небольшая задержка между вводами
		}
	}()
	return c, done
}

// negativeFilterStageInt - стадия фильтрации отрицательных чисел
func negativeFilterStageInt(done <-chan bool, c <-chan int) <-chan int {
	convertedIntChan := make(chan int)
	go func() {
		for {
			select {
			case data := <-c:
				if data > 0 {
					convertedIntChan <- data // Отправляем только положительные числа
				}
			case <-done:
				close(convertedIntChan) // Закрываем канал при завершении
				return
			}
		}
	}()
	return convertedIntChan
}

// specialFilterStageInt - стадия фильтрации чисел, не кратных 3
func specialFilterStageInt(done <-chan bool, c <-chan int) <-chan int {
	filteredIntChan := make(chan int)
	go func() {
		for {
			select {
			case data := <-c:
				if data != 0 && data%3 == 0 {
					filteredIntChan <- data // Отправляем только числа, кратные 3
				}
			case <-done:
				close(filteredIntChan) // Закрываем канал при завершении
				return
			}
		}
	}()
	return filteredIntChan
}

// bufferStageInt - стадия буферизации
func bufferStageInt(done <-chan bool, c <-chan int) <-chan int {
	bufferedIntChan := make(chan int)
	buffer := NewRingIntBuffer(bufferSize)

	go func() {
		for {
			select {
			case data := <-c:
				buffer.Push(data)
			case <-done:
				return
			}
		}
	}()

	// Вспомогательная горутина для периодической очистки буфера
	go func() {
		for {
			select {
			case <-time.After(bufferDrainInterval):
				bufferData := buffer.Get()
				// Если в кольцевом буфере что-то есть, выводим содержимое построчно
				if bufferData != nil {
					for _, data := range bufferData {
						select {
						case bufferedIntChan <- data:
						case <-done:
							return
						}
					}
				}
			case <-done:
				return
			}
		}
	}()
	return bufferedIntChan
}

// consumer - Потребитель данных от пайплайна2222222222
func consumer(done <-chan bool, c <-chan int) {
	for {
		select {
		case data := <-c:
			fmt.Printf("Обработаны данные: %d\n", data)
		case <-done:
			return
		}
	}
}
func main() {
	// Запускаем наш воображаемый источник данных,
	// он же ответственен за сигнализирование о том,
	// что он завершил работу
	source, done := dataSource()
	// Создаем пайплайн, передаем ему специальный канал,
	// синхронизирующий завершение работы пайплайна,
	// а также передаем ему все стадии
	pipeline := NewPipelineInt(done, negativeFilterStageInt, specialFilterStageInt, bufferStageInt)
	consumer(done, pipeline.Run(source))
}
