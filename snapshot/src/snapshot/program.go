//
// This file implements the program that process P and process Q rotate
// the message m.
//

package snapshot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"sync"
)

const (
	buffersize    = 256
	Debug         = false
	M             = "m"
	Marker        = "marker"
	stateOfRotate = iota
	stateOfFinishLocalSS
	stateOfStop
)

// Config stores the info of process
type Config struct {
	Name       string `json:"name"`
	LocalAddr  string `json:"localAddress"`
	RemoteAddr string `json:"remoteAddress"`
	Num        int    `json:"num"`
}

// Process is the struct of a process in the program.
type Process struct {
	name       string
	localAddr  string
	remoteAddr string
	numOfm     int
	l          net.Listener
	mu         sync.RWMutex
	state      int
	LocalState *LocalSS
}

type LocalSS struct {
	StateOfpc int
	ExtraMsg  []string
}

// createProcess creates a node in the program.
func CreateProcess(cfgFile string) (*Process, error) {
	cfg, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		Dprintf("read cfgFile error:", err)
		return nil, err
	}
	var config Config
	err = json.Unmarshal(cfg, &config)
	if err != nil {
		Dprintf("unmarshal cfgFile error:", err)
		return nil, err
	}

	pc := &Process{
		name:       config.Name,
		localAddr:  config.LocalAddr,
		remoteAddr: config.RemoteAddr,
		numOfm:     config.Num,
	}
	l, err := net.Listen("tcp", config.LocalAddr)
	if err != nil {
		Dprintf("listen error:", err)
		return nil, err
	}
	pc.l = l
	// pc.Finish = make(chan bool)
	pc.state = stateOfRotate
	pc.LocalState = &LocalSS{
		StateOfpc: pc.numOfm,
		ExtraMsg:  []string{},
	}

	return pc, nil
}

// Send sends a message to the remote process
func (pc *Process) Send(msg string) {
	conn, err := net.Dial("tcp", pc.remoteAddr)
	// defer conn.Close()
	if err != nil {
		Dprintf("dial error:", err)
		return
	}
	if msg == M {
		_, err = conn.Write([]byte("m"))
		Dprintf("%s sends m", pc.name)
	} else if msg == Marker {
		_, err = conn.Write([]byte("marker"))
		Dprintf("%s sends marker", pc.name)
	} else {
		Dprintf("msg is unvalid.")
	}
	if err != nil {
		Dprintf("write error:", err)
		return
	}
}

// Recv receives a message from the remote process
func (pc *Process) Recv() {
	Dprintf("%s is listening...", pc.name)
	for {
		conn, err := pc.l.Accept()
		if err != nil {
			Dprintf("accept error:", err)
			break
		}
		pc.handleRecv(conn)
	}
}

func (pc *Process) handleRecv(conn net.Conn) {
	buffer := make([]byte, buffersize)
	n, err := conn.Read(buffer)
	if err != nil {
		Dprintf("read error:", err)
		return
	}
	// conn.Close()
	msg := string(buffer[:n])
	if msg == M {
		// receives m
		pc.mu.Lock()
		if pc.state == stateOfFinishLocalSS {
			// extra m, just record, not update
			Dprintf("%s receives extra m", pc.name)
			pc.LocalState.ExtraMsg = append(pc.LocalState.ExtraMsg, msg)
		} else {
			Dprintf("%s receives m", pc.name)
			pc.numOfm++
			Dprintf("state of %s is %d", pc.name, pc.numOfm)
		}

		pc.Send(M)

		if pc.name == "P" && pc.numOfm == 101 {
			// process P start snapshot
			// record local state
			pc.LocalState.StateOfpc = pc.numOfm
			// send marker
			pc.Send(Marker)
			// start to record extra msgs
			pc.state = stateOfFinishLocalSS
		}
		pc.mu.Unlock()
	} else if msg == Marker {
		// receives marker
		Dprintf("%s receives marker", pc.name)
		pc.mu.Lock()
		if pc.state == stateOfRotate {
			// hasn't recorded yet
			pc.LocalState.StateOfpc = pc.numOfm
			pc.Send(Marker)
		}
		pc.state = stateOfStop
		pc.mu.Unlock()
	}
}

func (pc *Process) Check() bool {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	if pc.state == stateOfStop {
		return true
	} else {
		return false
	}
}

/*
func (pc *Process) recordLocalState() {
	pc.localState.stateOfpc = pc.numOfm
}
*/

func Dprintf(format string, v ...interface{}) {
	if Debug {
		log.Printf(format, v...)
	}
}
