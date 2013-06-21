package main

import (
	"net/http"
	"log"
)

type Client struct {
	channel chan <- string
}

func handleMessages(messageChan <- chan string, addChan <- chan Client, removeChan <- chan Client) {
	channels := make(map[Client] chan <- string)

	for {
		select {
		case message := <- messageChan:
			log.Print("New message: ", message)
			for _, channel := range channels {
				go func (c chan <- string) {
					c <- message
				}(channel)
			}
		case client := <- addChan:
			log.Print("Client connected: ", client)
			channels[client] = client.channel
		case client := <- removeChan:
			log.Print("Client disconnected: ", client)
			delete(channels, client)
		}
	}
}

func handleStream(messageChan chan <- string, addChan chan <- Client, removeChan chan <- Client, writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(200)

	channel := make(chan string)
	client  := Client{channel}

	addChan <- client

	for {
		message := <- channel
		if _, error := writer.Write([]byte(message + "\r\n")); error != nil {
			log.Print("Write: ", error)
			break
		}
		writer.(http.Flusher).Flush()
	}

	removeChan <- client
}

func handleMessage(messageChan chan <- string, writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	message := request.FormValue("message")
	messageChan <- message
	writer.WriteHeader(200)
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
	http.HandleFunc("/messages", func (writer http.ResponseWriter, request *http.Request) {
		handleMessage(messagesChan, writer, request)
	})

	log.Print("Starting server on :8080")

	http.ListenAndServe(":8080", nil)
}
