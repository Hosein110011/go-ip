package schema

import (
	"fmt"
	"github.com/Hosein110011/go-master/pkg/models"
)


type ApiResponse struct {
    Message    string                   `json:"message"`
    StatusCode int                      `json:"statusCode"`
    IsSuccess  bool                     `json:"isSuccess"`
    Result     []DataCenterApiResponse  `json:"result"`
}

type ApiResponse2 struct {
	Message    string                   `json:"message"`
    StatusCode int                      `json:"statusCode"`
    IsSuccess  bool                     `json:"isSuccess"`
    Result     []UrlApiResponse  `json:"result"`
}

type DataCenterApiResponse struct {
    DataCenterName string            `json:"data_center_name"`
    Types          map[string]map[string]HostMetrics `json:"types"`
}

type HostMetrics struct {
    LostPacketLastMin  string `json:"lost_packet_last_min"`
    LostPacketLastHour string `json:"lost_packet_last_hour"`
    LostPacketLastDay  string `json:"lost_packet_last_day"`
    RTTAvgLastHalf     string `json:"rtt_avg_last_half"`
}

type UrlApiResponse struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	Date    string `json:"date"`
	HourlyStatus  []CurlsApiResponse `json:"hourly_status"`
}

type CurlsApiResponse struct {
	Hour      string `json:"hour"`
	Details   []CurlApiResponse `json:"details"`
}

type CurlApiResponse struct {
	Time      string `json:"time"`
	Status    string `json:"status"`
} 

func GenerateApiResponse() (*ApiResponse, error) {
	var clouds []models.Cloud
	clouds = models.GetAllClouds()
	var result []DataCenterApiResponse
	for _, cloud := range clouds {
        types := make(map[string]map[string]HostMetrics)
		hosts := models.GetHostsByCloud(&cloud)
        for _, host := range hosts {
            if _, exists := types[host.IpType]; !exists {
                types[host.IpType] = make(map[string]HostMetrics)
            }

            // Assuming TotalLostPacket and TotalTime are slices with at least 3 and 2 elements respectively
            lenLostPackets := len(host.TotalLostPacket)
			var lastMin []float64
			if lenLostPackets >= 60 {
				lastMin = host.TotalLostPacket[lenLostPackets-60:]
			} else {
				lastMin = host.TotalLostPacket
			}
			var lastHour []float64
			if lenLostPackets >= 3600 {
				lastHour = host.TotalLostPacket[lenLostPackets-3600:]
			} else {
				lastHour = host.TotalLostPacket
			}
			var lastDay []float64
			if lenLostPackets >= 86400 {
				lastDay = host.TotalLostPacket[lenLostPackets-86400:]
			} else {
				lastDay = host.TotalLostPacket
			}
			var lastHalf []float64
			if len(host.TotalTime) >= 1800 {
				lastHalf = host.TotalTime[len(host.TotalTime)-1800:]
			} else {
				lastHalf = host.TotalTime
			}

			types[host.IpType][host.Ip] = HostMetrics{
                LostPacketLastMin:  fmt.Sprintf("%.2f", (sumFloats(lastMin)/60)*100),
                LostPacketLastHour: fmt.Sprintf("%.2f", (sumFloats(lastHour)/3600)*100),
                LostPacketLastDay:  fmt.Sprintf("%.2f", (sumFloats(lastDay)/86400)*100),
                RTTAvgLastHalf:     fmt.Sprintf("%.2f", (sumFloats(lastHalf))/1800),
            }
        }

        result = append(result, DataCenterApiResponse{
            DataCenterName: cloud.DatacenterName,
            Types:          types,
        })
    }

    return &ApiResponse{
        Message:    "hosts fetched successfully!",
        StatusCode: 200,
        IsSuccess:  true,
        Result:     result,
    }, nil
}


func sumFloats(slice []float64) float64 {
	var total float64
	for _, value := range slice {
		total += value
	}
	return total
}

func GenerateCurlApiResponse() (*ApiResponse2, error) {
	var urls []models.Url
	urls = models.GetAllUrls()
	UrlResponse = []UrlApiResponse{}
	for _, url := range urls {
		var curls []models.Curl
		curls = models.GetCurlsByUrl(&url)
		for i := 0; i <= 23; i++ {
			var CurlSerializer []CurlApiResponse
			var UrlDetails []CurlsApiResponse
			UrlDetails.Hour = i

		}
	} 

}