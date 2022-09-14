package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type speedTestRecords struct {
	Timestamp string `json:"timestamp"`
	UserInfo  struct {
		IP  string `json:"IP"`
		Lat string `json:"Lat"`
		Lon string `json:"Lon"`
		Isp string `json:"Isp"`
	} `json:"user_info"`
	Servers []struct {
		URL      string  `json:"url"`
		Lat      string  `json:"lat"`
		Lon      string  `json:"lon"`
		Name     string  `json:"name"`
		Country  string  `json:"country"`
		Sponsor  string  `json:"sponsor"`
		ID       string  `json:"id"`
		URL2     string  `json:"url_2"`
		Host     string  `json:"host"`
		Distance float64 `json:"distance"`
		Latency  int     `json:"latency"`
		DlSpeed  float64 `json:"dl_speed"`
		UlSpeed  float64 `json:"ul_speed"`
	} `json:"servers"`
}

func getSpeedTestResult() (string, float64, float64, int, string, string, string, string, string) {

	flag.Parse()
	app, arg1 := strings.Join(flag.Args()[:1], " "), strings.Join(flag.Args()[1:], " ")
	cmd := exec.Command(app, arg1)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	results := speedTestRecords{}
	json.Unmarshal([]byte(stdout), &results)

	upload := math.RoundToEven(results.Servers[0].UlSpeed)
	download := math.RoundToEven(results.Servers[0].DlSpeed)
	name := results.Servers[0].Name
	country := results.Servers[0].Country
	sponser := results.Servers[0].Sponsor
	latency := ((results.Servers[0].Latency) / 1000000)
	ipAddr := results.UserInfo.IP
	isp := results.UserInfo.Isp
	timestamp := results.Timestamp

	return timestamp, download, upload, latency, ipAddr, isp, sponser, name, country
}

func dbConnect(timestamp string, download float64, upload float64, latency int, ipAddr string, isp string, peerServer string) {

	pswd := os.Getenv("MYSQL_PASSWORD")
	user := os.Getenv("MYSQL_USER")
	database := os.Getenv("MYSQL_DATABASE")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dsn := user + ":" + pswd + "@tcp(" + host + ":" + port + ")/" + database

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Error while validating sql.Open arguments")
		panic(err.Error())
	}
	defer db.Close()

	conn := db.Ping()
	if conn != nil {
		fmt.Println("Error while connecting to database")
		panic(conn.Error())
	}

	insert, err := db.Prepare("INSERT INTO speedtest.speedtest (TimeStamp, DownloadSpeed, UploadSpeed, Latency, PublicIp, ISP, Peers) VALUES (? ,? ,? ,? ,? ,? ,?)")
	if err != nil {
		panic(err.Error())
	}
	insert.Exec(timestamp, download, upload, latency, ipAddr, isp, peerServer)
	defer insert.Close()
	fmt.Println("1 record inserted")
}

func main() {
	timestamp, download, upload, latency, ipAddr, isp, sponser, name, country := getSpeedTestResult()
	peerServer := sponser + " " + name + " " + " " + country
	dbConnect(timestamp, download, upload, latency, ipAddr, isp, peerServer)
}
