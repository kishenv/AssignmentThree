package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type UsageMetrics struct {
	CpuLoad        string // An exported string field for reporting the CPU load.
	MemoryUsagePct string //  An exported string filed for reporting the used/Total capacity of memory.
}

var (
	hbCounter                          int  // A variable to track the number of heartbeats missed to stop monitoring
	metricsReceived, heartBeatReceived bool // A variable of type bool to switch when server is initialized before client
	consecutiveHb                      bool // A variable of type bool to check for subsequent heartbeat.
)

func main() {
	go heartbeatCheck()
	//go metricsCheck(logger)
	log.Print("Client Monitoring service")
	log.Print("INFO: Starting health monitoring agent")
	mux := http.NewServeMux()
	mux.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) { heartbeatHandler(w, r) })
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) { usageMetricsHandler(w, r) })
	log.Print("INFO: Launched HTTP receiver for client events for metrics and heartbeats")
	http.ListenAndServe(":8080", mux)
}

// Handler when the endpoint /heartbeat is hit. Flaps the boolean variable to decide if the heartbeat is
// received well within the 3 second deadline.
func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	heartBeatReceived = true
	hbCounter = 0
	log.Print("INFO: Received Heartbeat from client")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"Message":"HeartbeatReceived"}`))
	consecutiveHb = true
}

// heartbeatCheck - Function to check for both metrics and heartbeats and log data.
// Since both Heartbeats and metrics are sent by the same service, the validation is done
// for both client liveliness and metrics.
func heartbeatCheck() {
	for {
		if !heartBeatReceived || !metricsReceived {
			log.Print("WARN: No heartbeats/Metrics received received yet.")
			time.Sleep(10 * time.Second)
		} else if consecutiveHb {
			time.Sleep(1 * time.Second)
			hbCounter++
			log.Print("INFO: Set TTL for ", hbCounter*3, " seconds until next Heartbeat")
			time.Sleep(3 * time.Second)
			if hbCounter == 3 {
				log.Print("WARN: Client unavailable, giving up on monitoring until next heartbeat")
				heartBeatReceived = false
				consecutiveHb = false
			}
		}
	}
}

// usageMetricsHandler - A handler for the endpoint `/metrics`` which accepts the data
// sent by the Client through POST operation.
func usageMetricsHandler(w http.ResponseWriter, r *http.Request) {
	metricsReceived = true
	var mData UsageMetrics
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print("ERROR: Cannot decode payload body", err)
		http.Error(w, "Malformed data received.", http.StatusBadRequest)
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	err = json.Unmarshal(b, &mData)
	if err != nil {
		log.Print("ERROR: Error during unmarshaling of data", err)
		http.Error(w, "Malformed data received.", http.StatusBadRequest)
	}
	log.Printf("INFO : CPU usage: %s Memory Usage: %s", mData.CpuLoad, mData.MemoryUsagePct)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"Message":"MetricsReceived"}`))
}
