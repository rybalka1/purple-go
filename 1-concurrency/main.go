package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// rand.Seed(time.Now().UnixNano())

	ch := make(chan int)
	squaredCh := make(chan int)

	// Горутина 1: генерация случайных чисел
	go func() {
		slice := make([]int, 10)
		for i := 0; i < 10; i++ {
			slice[i] = rand.Intn(101) // Случайное число от 0 до 100
		}
		for _, num := range slice {
			ch <- num
		}
		close(ch)
	}()

	// Горутина 2: возведение в квадрат
	go func() {
		for num := range ch {
			squaredCh <- num * num
		}
		close(squaredCh)
	}()

	// Основная горутина: получение и вывод квадратов
	var results []int
	for i := 0; i < 10; i++ {
		result := <-squaredCh
		results = append(results, result)
	}
	fmt.Println(results)
}
