package main

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Hosein110011/go-master/pkg/models"
	"github.com/Hosein110011/go-master/pkg/routes"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	erre := godotenv.Load(".env") // Load .env file
	if erre != nil {
		log.Fatal("Error loading .env file", erre)
	}

	http.HandleFunc("/ws/ping/", wsHandler)
	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)
	http.Handle("/", r)
	portStr := os.Getenv("PORT")
	port, errr := strconv.Atoi(portStr)
	if errr != nil {
		log.Fatalf("Invalid port number: %v\n", errr)
	}

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
				host.TotalLostPacket = append(host.TotalLostPacket, data["packet_loss_count"].(float64))
				host.TotalTime = append(host.TotalTime, data["rtt_avg"].(float64))
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
				host.TotalLostPacket = append(host.TotalLostPacket, data["packet_loss_count"].(float64))
				host.TotalTime = append(host.TotalTime, data["rtt_avg"].(float64))
				// host.CloudID = cloud.ID
				models.UpdateHost(host)
				// hosts := models.GetHostsByCloud(cloud)
				// fmt.Println(hosts)
			}

		} else {
			fmt.Println("CURL")
			Url := &models.Url{}
			Curl := &models.Curl{}
			Url = models.GetUrlByUrl(data["url"].(string))
			if Url.ID == 0 {
				Url.Name = data["name"].(string)
				Url.Date = time.Now()
				Url.Url = data["url"].(string)
				Url.CreateUrl()
			}
			Curl.UrlID = Url.ID
			Curl.Status = int64(data["status_code"].(float64))
			Curl.Time = time.Now()
			Curl.CreateCurl()
			Url.Curls = append(Url.Curls, *Curl)
			models.UpdateUrl(Url)
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
