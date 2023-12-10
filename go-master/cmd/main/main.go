package main

import (
    // "encoding/json"
    "fmt"
	"encoding/json"
	"github.com/gorilla/mux"
    "github.com/Hosein110011/go-master/pkg/routes"
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
	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)
	http.Handle("/", r)
	port := 9001 // Use a different port for the WebSocket server
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
			cloud = models.GetCloudByDcName(data["data_center"].(string))
			if cloud.ID == 0 {
				cloud.DatacenterName = data["data_center"].(string)
				cloud = cloud.CreateCloud()
				host = models.GetHostByIp(data["destination"].(string))
				if host.ID != 0 {
					return
				}             
				if host.ID == 0 {
					host.Ip = data["destination"].(string)
					host.IpType = data["type"].(string)
					host.CloudID = cloud.ID
					host.CreateHost()
					cloud.Hosts = append(cloud.Hosts, *host)
					models.UpdateCloud(cloud)
				}
				host.TotalLostPacket = append(host.TotalLostPacket ,data["packet_loss_count"].(float64))
				host.TotalTime = append(host.TotalTime ,data["rtt_avg"].(float64))
				// host.CloudID = cloud.ID
				models.UpdateHost(host)
			} else {
				host = models.GetHostByIp(data["destination"].(string))
				if host.ID == 0 {
					fmt.Println("0")
					host.Ip = data["destination"].(string)
					host.IpType = data["type"].(string)
					host.CloudID = cloud.ID
					host.CreateHost()
					cloud.Hosts = append(cloud.Hosts, *host)
					models.UpdateCloud(cloud)
				}
				host.TotalLostPacket = append(host.TotalLostPacket ,data["packet_loss_count"].(float64))
				host.TotalTime = append(host.TotalTime ,data["rtt_avg"].(float64))
				// host.CloudID = cloud.ID
				models.UpdateHost(host)
				hosts := models.GetHostsByCloud(cloud)
				fmt.Println(hosts)
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