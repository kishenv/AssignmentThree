// Client side script that sends heartbeats to the master node on a regular basis, and the OS stats.
// Use log rotation

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

var group sync.WaitGroup
var gigaB float64
var httpClient *http.Client

type UsageMetrics struct {
	CpuLoad        string
	MemoryUsagePct string
}

func main() {
	httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
	gigaB = math.Pow(10, 9)
	log.Print("INFO: Initializing client code to send Heartbeats and CPU/ Memory usage profile to server")
	group.Add(1)
	go sendHeartBeatstoServer()
	go sendCPUandMemoryProfile()
	group.Wait()
}

// sendHeartBeatstoServer() - Posts data to the metrics server which tracks the CPU/Memory metrics
// and the heartbeat of the client every 3s.
func sendHeartBeatstoServer() {
	for {
		log.Print("INFO: Sending Heartbeat to server")
		healthMsg := []byte(`{"Health":"Ok"}`)
		req, err := http.NewRequest(http.MethodPost, "http://metrics-server.default.svc.cluster.local:8080/heartbeat", bytes.NewReader(healthMsg))
		if err != nil {
			log.Print("ERROR: Error while creating a request", err)
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Print("ERROR: Error while creating a request", err)
			continue
		}
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("INFO: Heartbeat API response. statusCode: %d ResposeBody: %s", resp.StatusCode, string(respBytes))
		resp.Body.Close()
		time.Sleep(time.Second * 3)
	}
}

// sendCPUandMemoryProfile() - Fetches the CPU and Memory metrics and POSTs them to the metrics-server
// every 20s
func sendCPUandMemoryProfile() {
	for {
		memory, err := memory.Get()
		if err != nil {
			log.Printf("ERROR: Failed to get the Memory related information", err)
			return
		}
		memoryUsage := fmt.Sprintf("%.3f", (float64(memory.Used)/gigaB)/(float64(memory.Total)/gigaB)*100)
		before, err := cpu.Get()
		if err != nil {
			log.Printf("ERROR: Failed to get the CPU related information", err)
			return
		}
		time.Sleep(time.Duration(50) * time.Millisecond)
		after, err := cpu.Get()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}
		total := float64(after.Total - before.Total)
		cpuUsage := fmt.Sprintf("%.3f", float64(after.User-before.User)+float64(after.System-before.System)/total*100)
		log.Printf("INFO: CPU usage: %s. Memory Usage: %s%%", cpuUsage, memoryUsage)
		data, _ := json.Marshal(UsageMetrics{CpuLoad: cpuUsage, MemoryUsagePct: memoryUsage})
		log.Print("INFO: Sending CPU and mem profile to server")
		req, err := http.NewRequest(http.MethodPost, "http://metrics-server.default.svc.cluster.local:8080/metrics", bytes.NewReader(data))

		if err != nil {
			log.Print("ERROR: Error while creating a request", err)
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Print("ERROR: Error while creating a request", err)
			continue
		}
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("INFO: Metrics API response statusCode: %d ResposeBody: %s", resp.StatusCode, string(respBytes))
		resp.Body.Close()
		time.Sleep(20 * time.Second)
	}
}
