//
// This file starts the program,  then use the snapshot algorithm to monitor the
// program, and record the snapshot.
//

package main

import (
	"fmt"
	"snapshot"
	"time"
)

func main() {
	p, err := snapshot.CreateProcess("process1.json")
	if err != nil {
		snapshot.Dprintf("create process p error:", err)
		return
	}
	q, err := snapshot.CreateProcess("process2.json")
	if err != nil {
		snapshot.Dprintf("create process p error:", err)
		return
	}

	go p.Recv()
	go q.Recv()

	time.Sleep(time.Second)
	// P sends m first.
	p.Send(snapshot.M)
	fmt.Println("Start rotate.")

	// waiting...
	for !p.Check() || !q.Check() {
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Snapshot completed.")
	// print global state
	fmt.Println("global state:")
	fmt.Printf("Process P's state: %d\n", p.LocalState.StateOfpc)
	fmt.Printf("Process P's extra msgs: %v\n", p.LocalState.ExtraMsg)
	fmt.Printf("Process Q's state: %d\n", q.LocalState.StateOfpc)
	fmt.Printf("Process Q's extra msgs: %v\n", q.LocalState.ExtraMsg)
}
