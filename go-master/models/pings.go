package models

import (
	// "gorm.io/driver/postgres"
  	"gorm.io/gorm"
	"github.com/Hosein1100011/go-api/pkg/config"
)

var db *gorm.DB

type Cloud struct {
	gorm.Model
	CloudID          uint   `gorm:"primaryKey"`
	DatacenterName   string `gorm:"uniqueIndex;not null"`
	Hosts      		 []Host `json:"hosts"`
}

type Host struct {
	gorm.Model
	CloudID uint `gorm:"uniqueIndex;not null"`
	IpType  string `json:"ip_type"`
	Ip      string `gorm:"uniqueIndex;not null"`
	TotalLostPacket []int `json:"total_lost_packet"`
	TotalTime  []float64 `json:"total_time"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func init() {
	config.Connect()
    db = config.GetDB()
	db.AutoMigrate(&Cloud{})
	db.AutoMigrate(&Host{})
}

func (c *Cloud) CreateCloud() *Cloud {
	db.Create(&c)
	return c
}

func (h *Host) CreateHost() *Host {
	db.Create(&h)
    return h
}

func GetAllClouds() []Cloud {
	var Clouds []Cloud
	db.Find(&Clouds)
	return Clouds
}

