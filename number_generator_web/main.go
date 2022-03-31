package main
import (
	"fmt"
	"math/rand"
	"net/http"
	"html/template"
	"strconv"
	"sync"
	"sort"
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
func uniqueGenerator(channelNums chan int,  exitChan chan int, templateChan chan int, limit int, wg *sync.WaitGroup) {

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
				templateChan <- value //вывод
			}
			return //заканчиваем работу
		}
	}
}


func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		tmpl, _ := template.ParseFiles("index.tmpl") //рендерим шаблон
		Numbers := false
		tmpl.Execute(w, Numbers)
	})

	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request){
		rand.Seed(time.Now().UTC().UnixNano())

		limit, err := strconv.ParseInt(r.FormValue("limit")[0:], 10, 64);
		if err != nil{
			return
		}
		max, err := strconv.ParseInt(r.FormValue("range")[0:], 10, 64);
		if err != nil{
			return
		}
		goCount, err := strconv.ParseInt(r.FormValue("go")[0:], 10, 64);
		if err != nil{
			return
		}

		//валидация
		if (limit <= 0 || max <= 0) || (limit > max) || (goCount < 1){
			tmpl, _ := template.ParseFiles("index.tmpl")
			tmpl.Execute(w, "Invalid input")
			return
		}

		var wg sync.WaitGroup
		channelNums := make(chan int, limit) //канал для передачи чисел
		exitChan := make(chan int) //канал для передачи флага о прекращении работы
		templateChan := make(chan int, limit) //канал для передачи чисел в шаблон

		for i:=0; i<int(goCount); i++{
			wg.Add(1)
			go randomizer(channelNums, exitChan, int(max), &wg) //запускаем n-ное количество горутин
		}

		wg.Add(1)
		go uniqueGenerator(channelNums, exitChan, templateChan, int(limit), &wg) //генерируем уникальные значения

		data := make([]int, 0, int(limit))
		count := 0
		for val := range templateChan{ //читаем из канала и перенаправляем в массив
			data = append(data, val)
			count += 1
			if count == int(limit){
				break
			}
		}

		wg.Wait()

		tmpl, _ := template.ParseFiles("index.tmpl")
		tmpl.Execute(w, data) //рендерим шаблон
	})

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil) //поднимаем сервер на порту 8181
}
