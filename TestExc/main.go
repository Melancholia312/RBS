package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)


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

	//Создаем словарь и заполняем его, чтобы проверять уникальность числа. true - значит число уже было
	numDict := make(map[int]bool)

	for {
		randNum := rand.Intn(*max) + 1 //Добавляем единицу ради устранения нуля
		if _, ok := numDict[randNum]; !ok {
			numDict[randNum] = true
		}
		if len(numDict) == *limit {
			for name := range numDict {
				fmt.Println(name)
			}
			return
		}
	}

}



