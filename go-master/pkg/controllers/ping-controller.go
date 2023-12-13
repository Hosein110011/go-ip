package controllers

import (
	"encoding/json"
	"fmt"

	// "github.com/gorilla/mux"
	"net/http"
	// "strconv"
	"github.com/Hosein110011/go-master/pkg/schema"
	// "github.com/Hosein110011/go-master/pkg/models"
)

func GetDataCenter(w http.ResponseWriter, r *http.Request) {
	// clouds := models.GetAllClouds()
	// for _, cloud := range clouds {
	// 	cloud.Hosts = models.GetHostsByCloud(&cloud)
	// }
	// res, _ := json.Marshal(clouds)
	result, err := schema.GenerateApiResponse()
	if err != nil {
		fmt.Println("Error generating:", err)
	}
	res, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetUrls(w http.ResponseWriter, r *http.Request) {
	result, _ := schema.GenerateCurlApiResponse()
	res, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}


