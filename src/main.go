package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var users []*websocket.Conn

// why go
type Message struct {
	Type    string `json:"type"`
	Id      string `json:"id"`
	Payload string `json:"payload"`
}

var addr = flag.String("ws", "localhost:8080", "http service")

var upgrader = websocket.Upgrader{}

func echo(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)

	if err != nil {
		log.Print("Upgrade Failed")
		return
	}

	defer conn.Close()

	for {
		mt, message, err := conn.ReadMessage()

		if err != nil {
			log.Print("Read Message Failed")
			break
		}

		var jsonData Message
		json.Unmarshal(message, &jsonData)
		log.Println("The messsage type is ", jsonData)

		switch jsonData.Type {
		case "join":
			users = append(users, conn)
			log.Print(users)
		case "message":
			for _, val := range users {
				if conn != val {
					log.Print("Id is ", jsonData.Id)
					err = val.WriteMessage(mt, message)
					if err != nil {
						log.Print("Writing Message Failed ", err)
						break
					}
				}
			}
		}

	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
