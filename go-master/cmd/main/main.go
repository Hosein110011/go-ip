package main

import (
    // "encoding/json"
    "fmt"
	"encoding/json"
    // "io/ioutil"
	"github.com/gorilla/websocket"
	"net/http"
	"log"
)


var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


func main() {
	http.HandleFunc("/ws/ping/", wsHandler)

	port := 8000 // Use a different port for the WebSocket server
	serverAddr := fmt.Sprintf(":%d", port)
	fmt.Printf("WebSocket Server listening on http://localhost%s\n", serverAddr)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Println("Error:", err)
	}
	select {}
}


func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("Received message: %s\n", p)
		var data map[string]interface{}
		err = json.Unmarshal(p, &data)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			continue
		}

		// Accessing keys and values in the map
		for key, value := range data {
			// fmt.Printf("Key: %s, Value: %v\n", key, value)
			if data["duty"] == "ping" {
				switch key {
				case "data":
					subData, ok := value.(map[string]interface{})
					if !ok {
						log.Println("Error decoding nested JSON:", err)
						continue
					}
					if subData["type"] == "ping" {
						fmt.Println("Received ping")
					} else {
						fmt.Println("CURL")
					}

					// Accessing keys and values in the nested map
					for subKey, subValue := range subData {
						fmt.Printf("Nested Key: %s, Nested Value: %v\n", subKey, subValue)
					}
				default:
					fmt.Printf("Key: %s, Value: %v\n", key, value)
				
				}
			} else {
				fmt.Println("CURL")
			}
		}
		if err != nil {
			log.Println(err)
            return
		}
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			log.Println(err)
			return
		}
	}
}