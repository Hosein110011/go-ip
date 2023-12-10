package main

import (
    // "encoding/json"
    "fmt"
	"encoding/json"
    // "io/ioutil"
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"github.com/Hosein110011/go-master/pkg/models"
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
			// fmt.Printf("Key: %s, Value: %v\n", key, value)
		if data["duty"] == "ping" {
			cloud := &models.Cloud{}
			host := &models.Host{}
			// db := &models.GetDB()
			cloud = models.getCloudByDcName(data["data_center"].(string))
			if cloud == (&models.Cloud{}) {
				cloud.DatacenterName = data["data_center"].(string)
				host.Ip = data["destination"].(string)
				host.IpType = data["type"].(string)
				host.TotalLostPacket = append(host.TotalLostPacket ,data["packet_loss_count"].(int))
				host.TotalTime = append(host.TotalTime ,data["rtt_avg"].(float64))
				cloud = cloud.CreateCloud()
				host = host.CreateHost()
				cloud.Hosts = append(cloud.Hosts, *host)
			}
		} else {
			fmt.Println("CURL")
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