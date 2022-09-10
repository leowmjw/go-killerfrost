package main

import (
	"app/internal/teleport"
	"context"
	"encoding/json"
	"fmt"
	"go.temporal.io/sdk/client"
	"os"
	"time"
)

func main() {
	fmt.Println("Breakglass CLI!!")
	os.Exit(run())
}

func run() int {
	c, err := client.Dial(client.Options{
		Namespace: "default",
	})
	if err != nil {
		panic(err)
	}
	// Reusable workflowID; should be use controllable; within validated list?
	wfID := "background"
	// Signal to kick off the process for access ..
	serr := c.SignalWorkflow(context.Background(), wfID, "",
		teleport.SignalName,
		teleport.BreakGlassSignal{
			Action: teleport.BG_REQUEST_ACCESS,
			Body:   json.RawMessage(`{"content": "dood"}`),
		})
	if serr != nil {
		fmt.Println("ERR_SIG: ", serr)
	}
	// SImulate approve; reject ..
	time.Sleep(time.Second * 10)
	serr = c.SignalWorkflow(context.Background(), wfID, "",
		teleport.SignalName,
		teleport.BreakGlassSignal{
			Action: teleport.BG_REQUEST_APPROVED,
			Body: json.RawMessage(`{
	"userName": "backend",
	"roleName": "s2read"
}`),
		})
	if serr != nil {
		fmt.Println("ERR_SIG: ", serr)
	}

	time.Sleep(time.Second * 15)
	// Simulate pending requests queue
	serr = c.SignalWorkflow(context.Background(), wfID, "",
		teleport.SignalName,
		teleport.BreakGlassSignal{
			Action: teleport.BG_OPS_DUMP,
			Body:   nil,
		})
	if serr != nil {
		fmt.Println("ERR_SIG: ", serr)
	}
	return 0
}
