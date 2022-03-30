package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

//randomizer генерирует случайные числа
func randomizer(channelNums chan int, exitChan chan int, max int, wg *sync.WaitGroup){

	defer wg.Done()
	for {
		randNum := rand.Intn(max) + 1 //генерация числа. +1 для избежания нуля
		select {
			case channelNums <- randNum: //отправляем числа
			case <- exitChan: //при получении числа в этот канал заканчиваем работу горутины
				return
		}

	}
}

//uniqueGenerator проверяет полученные числа на уникальность и выводит их на экран
func uniqueGenerator(channelNums chan int,  exitChan chan int, limit int, wg *sync.WaitGroup) {

	defer wg.Done()
	uniqueDict := make(map[int]bool) //словарь для проверки уникальности

	for val := range channelNums{
		if _, ok := uniqueDict[val]; !ok {
			uniqueDict[val] = true //true - значит число уже было
		}
		if len(uniqueDict) == limit{ //при наборе нужного количества уникальных чисел говорим горутинам прекратить работу
			exitChan <- 1
			close(exitChan)
			keys := make([]int, 0, len(uniqueDict))
			for k := range uniqueDict {
				keys = append(keys, k)
			}
			//Сортировка полученных значений
			sort.Ints(keys)
			for _, value := range keys{
				fmt.Println(value) //вывод на экран
			}
			return //заканчиваем работу
		}
	}
}


func main() {
	//Меняем сид при каждом запуске
	rand.Seed(time.Now().UTC().UnixNano())

	///Парсинг и валидация флагов
	limit := flag.Int("limit", 0, "# of rand numbers")
	max := flag.Int("max", 1, "the max value")
	goCount := flag.Int("go", 1, "the max value")
	flag.Parse()

	if (*limit <= 0 || *max <= 0) || (*limit > *max) || (*goCount < 1) {
		fmt.Println("Invalid input")
		return
	}

	var wg sync.WaitGroup
	channelNums := make(chan int, *limit) //канал для передачи чисел
	exitChan := make(chan int) //канал для передачи флага о прекращении работы

	for i:=0; i<*goCount; i++{
		wg.Add(1)
		go randomizer(channelNums, exitChan, *max, &wg) //запускаем n-ное количество горутин
	}

	wg.Add(1)
	go uniqueGenerator(channelNums, exitChan, *limit, &wg)

	wg.Wait() //ждем окончание работы горутин

}
