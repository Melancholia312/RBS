package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

//Вывод числе на экран
func printNumbers(finalArr *[]int) {
	for _, value := range *finalArr {
		fmt.Println(value)
	}
}

//Создание словаря
func generateUniqueNumbers(max int, limit int) *[]int {
	//Создаем словарь и заполняем его, чтобы проверять уникальность числа. true - значит число уже было
	numDict := make(map[int]bool)

	for {
		randNum := rand.Intn(max) + 1 //Добавляем единицу ради устранения нуля
		if _, ok := numDict[randNum]; !ok {
			numDict[randNum] = true
		}
		if len(numDict) == limit {
			//Перемещаем полученные значения в массив
			keys := make([]int, 0, len(numDict))
			for k := range numDict {
				keys = append(keys, k)
			}
			//Сортировка полученных значений
			sort.Ints(keys)
			return &keys
		}
	}
}

func main() {
	//Меняем сид при каждом запуске
	rand.Seed(time.Now().UTC().UnixNano())

	///Парсинг и валидация флагов
	limit := flag.Int("limit", 0, "# of rand numbers")
	max := flag.Int("max", 1, "the max value")
	flag.Parse()

	if (*limit <= 0 || *max <= 0) || (*limit > *max) {
		fmt.Println("Invalid input")
		return
	}

	printNumbers(generateUniqueNumbers(*max, *limit))

}
