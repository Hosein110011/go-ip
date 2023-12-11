package models

import (
	// "gorm.io/driver/postgres"
	"time"
  	"gorm.io/gorm"
	"github.com/Hosein110011/go-master/pkg/config"
	pq "github.com/lib/pq"
)

var db *gorm.DB

type Cloud struct {
	gorm.Model
	ID          	 uint   `gorm:"primaryKey;autoIncrement"`
	DatacenterName   string `gorm:"column:datacentername;uniqueIndex;not null"`
	Hosts      		 []Host `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:CloudID"`
}

type Host struct {
	gorm.Model
	// ID      uint `gorm:"primaryKey;autoIncrement"`
	CloudID uint   `gorm:"index;not null"`
	IpType  string `gorm:"column:iptype;"`
	Ip      string `gorm:"column:ip;uniqueIndex;not null"`
	TotalLostPacket pq.Float64Array `gorm:"column:totallostpacket;type:float[]"`
	TotalTime  pq.Float64Array `gorm:"column:totaltime;type:float[]"`
}

type Url struct {
	gorm.Model 
	ID      uint `gorm:"primaryKey;autoIncrement"`
	Name    string `gorm:"column:name;not null"`
	Url     string `gorm:"column:url;uniqueIndex;not null"`
	Date    time.Time `gorm:"column:date;not null"`
	Curls   []Curl    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UrlID"`
}

type Curl struct {
	gorm.Model
	UrlID      uint `gorm:"column:urlid;index;not null"`
	Status     int64 `gorm:"column:status;not null"`
	Time       time.Time `gorm:"column:time;not null"`
}


func init() {
	config.Connect()
    db = config.GetDB()
	db.AutoMigrate(&Cloud{})
	db.AutoMigrate(&Host{})
	db.AutoMigrate(&Url{})
	db.AutoMigrate(&Curl{})
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

func GetCloudByDcName(dcName string) *Cloud {
	var cloud Cloud
    db.Where("DatacenterName=?", dcName).Find(&cloud)
	return &cloud
}


func GetHostsByCloudId(cloud *Cloud) ([]Host, error){
	if err := db.Where("DatacenterName = ?", cloud.DatacenterName).Preload("Hosts").First(&cloud).Error; err != nil {
		// fmt.Println("Error:", err)
		return nil, err
	}	
	return cloud.Hosts, nil 
}

func GetHostByIp(ip string) *Host {
	var host Host
    db.Where("Ip=?", ip).Find(&host)
	return &host
}

func UpdateHost(host *Host) {
	db.Save(&host)
}

func UpdateCloud(cloud *Cloud) {
	db.Save(&cloud)
}

func GetHostsByCloud(c *Cloud) []Host {
	var cloud Cloud
	db.Preload("Hosts").First(&cloud, c.ID)
	return cloud.Hosts
}

func (u *Url) CreateUrl() *Url {
	db.Create(&u)
    return u
}

func (c *Curl) CreateCurl() *Curl {
	db.Create(&c)
    return c
}

func GetUrlByUrl(u string) *Url {
	var url Url
    db.Where("Url=?", u).Find(&url)
	return &url
}

func UpdateUrl(url *Url) {
	db.Save(&url)
}

func GetCurlsByUrl(url *Url) ([]Curl, error){
	// First, find the Url record
    if err := db.Where("url = ?", url.Url).First(url).Error; err != nil {
        return nil, err
    }

    // Then, load associated Curls sorted by Time
    var curls []Curl
    if err := db.Model(&Curl{}).Where("UrlID = ?", url.ID).Order("time DESC").Limit(288).Find(&curls).Error; err != nil {
        return nil, err
    }

    return curls, nil
}

func GetAllUrls() []Url {
	var Urls []Url
	db.Find(&Urls)
	return Urls
}

