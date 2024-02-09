package c2

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type HttpListener struct {
	Ip         string
	Port       string
	errorLog   *log.Logger
	requestLog *log.Logger
}

func (h *HttpListener) getTasksHandler(w http.ResponseWriter, r *http.Request) {
	h.requestLog.Printf("Time: %v, Method: %s, URL: %s, RemoteAddr: %s\n", time.Now(), r.Method, r.URL, r.RemoteAddr)
	task := map[string]interface{}{
		"cmd":  "echo",
		"args": []string{"latest", "config"},
	}
	jsonData, err := json.Marshal(task)
	if err != nil {
		h.errorLog.Println("Error marshalling JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *HttpListener) handler(w http.ResponseWriter, r *http.Request) {
	if id := r.Header.Get("Cookie"); id != "" {
		switch id {
		case "gt":
			h.getTasksHandler(w, r)
		default:
			return
		}
	}
}

func (h *HttpListener) StartListener() {
	// open log file for listener
	logFile, err := os.OpenFile("c2/data/listener.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %s", err.Error())
	}
	defer logFile.Close()

	h.errorLog = log.New(logFile, "ERROR: ", log.LstdFlags)
	h.requestLog = log.New(logFile, "REQUEST: ", log.LstdFlags)

	// init server
	server := &http.Server{
		Addr:     fmt.Sprintf("%s:%s", h.Ip, h.Port),
		ErrorLog: h.errorLog,
		Handler:  http.HandlerFunc(h.handler),
	}

	// start listening
	err = server.ListenAndServe()
	if err != nil {
		h.errorLog.Fatalf("Server failed to start: %v", err)
	}
}
