package main

import (
	"net/http"
	"log"
	"time"
)

type Client struct {
	writer http.ResponseWriter
	channel chan <- string
}

func handleMessages(messageChan <- chan string, addChan <- chan Client, removeChan <- chan Client) {
	// clients := make(map[http.ResponseWriter] chan <- string)

	for {
		select {
		case message := <- messageChan:
			log.Print("New message: ", message)
		case client := <- addChan:
			log.Print("Client connected: ", client)
		case client := <- removeChan:
			log.Print("Client disconnected: ", client)
		}
	}
}

func handleStream(messageChan chan <- string, addChan chan <- Client, removeChan chan <- Client, writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(200)

	channel := make(chan string)
	client  := Client{writer, channel}

	addChan <- client

	for {
		if _, error := writer.Write([]byte("test\r\n")); error != nil {
			log.Print("Write: ", error)
			break
		}
		writer.(http.Flusher).Flush()
		time.Sleep(time.Second)
	}

	removeChan <- client
}

func main() {
	messagesChan := make(chan string)
	addChan      := make(chan Client)
	removeChan   := make(chan Client)

	go handleMessages(messagesChan, addChan, removeChan)

	http.HandleFunc("/", func (writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "static/index.html")
	})
	http.HandleFunc("/static/", func (writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, request.URL.Path[1:])
	})
	http.HandleFunc("/stream", func (writer http.ResponseWriter, request *http.Request) {
		handleStream(messagesChan, addChan, removeChan, writer, request)
	})

	log.Print("Starting server on :8080")

	if error := http.ListenAndServe(":8080", nil); error != nil {
		log.Fatal("ListenAndServe: ", error)
	}

	log.Print("yeah");
}
