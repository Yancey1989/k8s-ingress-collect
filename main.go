package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/topicai/candy"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func run(host string) {
	nginxHost := getenv("NGINX_HOST", "192.168.61.73")
	nginxPort := getenv("NGINX_PORT", "46193")
	influxDBHost := getenv("INFLUX_DB_HOST", "k8sinfluxapi.k8s.baifendian.com")
	influxDBPort := getenv("INFLUX_DB_PORT", "80")
	influxDBName := getenv("INFLUX_DB_NAME", "test")
	nginxStatusURL := fmt.Sprintf("http://%s:%s/nginx_status/format/json",
		nginxHost, nginxPort)
	influxDBURL := fmt.Sprintf("http://%s:%s/write?db=%s", influxDBHost, influxDBPort, influxDBName)
	valList := []string{}
	timestamp := time.Now().Unix() * 1000000000
	resp, err := http.Get(nginxStatusURL)
	candy.Must(err)
	body, err := ioutil.ReadAll(resp.Body)
	candy.Must(err)
	var v map[string]interface{}
	json.Unmarshal(body, &v)

	// Formatting request stat
	for key, value := range v["connections"].(map[string]interface{}) {
		tmp := fmt.Sprintf("nginx_value,type=nginx_connections,type_instance=%s,host=%s value=%d %d", key, host, int(value.(float64)), timestamp)
		valList = append(valList, tmp)
	}

	// Formatting server zone
	for key, value := range v["serverZones"].(map[string]interface{}) {
		requestCounter := int(value.(map[string]interface{})["requestCounter"].(float64))
		tmp := fmt.Sprintf("nginx_value,type=nginx_zone,type_instance=requestCounter,zone=%s,host=%s value=%d %d", key, host, requestCounter, timestamp)
		valList = append(valList, tmp)
	}

	points := strings.Join(valList, "\n")
	fmt.Println(points)

	// Write date to influxdb
	resp, err = http.Post(influxDBURL, "application/x-www-form-urlencoded", strings.NewReader(points))
	candy.Must(err)
	body, err = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	return
}

// GetLocalIP returns the non loopback local IP of the host
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
func main() {
	interval, _ := strconv.Atoi(getenv("INTERVAL", "30"))
	localIP := getLocalIP()
	//c := make(chan int)
	//go run(c, localIP)
	for {
		go run(localIP)
		//c <- 1
		time.Sleep(time.Second * time.Duration(interval))
	}
}
