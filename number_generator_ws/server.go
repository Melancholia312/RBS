package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"strconv"
	"strings"
	"sync"
	"math/rand"
	"time"
)

var upgrader = websocket.Upgrader{}

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

//uniqueFilterr проверяет полученные числа на уникальность
func uniqueFilter(channelNums chan int,  exitChan chan int, limit int, ws *websocket.Conn,  wg *sync.WaitGroup){

	defer wg.Done()
	uniqueDict := make(map[int]bool) //словарь для проверки уникальности

	for val := range channelNums{
		if _, ok := uniqueDict[val]; !ok {
			uniqueDict[val] = true //true - значит число уже было
			bs := []byte(strconv.Itoa(val))
			err := ws.WriteMessage(websocket.TextMessage, bs) // отправляем числа по вебсокету
			if err != nil {
				log.Println("Error during message writing:", err)
				break
			}
		}
		if len(uniqueDict) == limit{ //при наборе нужного количества уникальных чисел говорим горутинам прекратить работу
			exitChan <- 1
		}
	}
}

//uniqueGenerator генерирует массив со случайными числами
func uniqueGenerator(limit int, max int, goCount int, ws *websocket.Conn) {
	rand.Seed(time.Now().UTC().UnixNano())
	if (limit <= 0 || max <= 0) || (limit > max) || (goCount < 1){ //проверяем валидность данных
		err := ws.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(0)))
		if err != nil {
			log.Println("Error during message writing:", err)
		}
		return
	}
	var wg sync.WaitGroup
	channelNums := make(chan int, limit) //канал для передачи чисел
	exitChan := make(chan int) //канал для передачи флага о прекращении работы
	for i:=0; i<goCount; i++{
		wg.Add(1)
		go randomizer(channelNums, exitChan, max, &wg) //запускаем n-ное количество горутин
	}
	wg.Add(1)
	go uniqueFilter(channelNums, exitChan, limit, ws, &wg) //генерируем уникальные значения
	wg.Wait()
}

//socketHandler отвечает за обработку веб сокета
func socketHandler(w http.ResponseWriter, r *http.Request, ) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	for { //читаем сообщения с клиентской стороны
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		s := fmt.Sprintf("%s", message) //парсим форму
		split := strings.Split(s, ",")
		limit, _ := strconv.Atoi(split[0])
		max, _ := strconv.Atoi(split[1])
		gor, _ :=  strconv.Atoi(split[2])

		go uniqueGenerator(limit, max, gor, conn)
	}
}

//home домашняя страница
func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	//роутинг
	http.HandleFunc("/", home)
	http.HandleFunc("/socket", socketHandler)

	//статика
	fileServer := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))
	//запуск сервера
	log.Fatal(http.ListenAndServe("localhost:4000", nil))
}
