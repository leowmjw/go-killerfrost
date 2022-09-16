package main

import (
	"app/internal/breakglass"
	"context"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
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
	wfID := breakglass.TestWFID
	// Signal to kick off the process for access .. start if completed ..
	wfr, serr := c.SignalWithStartWorkflow(context.Background(), wfID,
		breakglass.SignalName,
		breakglass.Signal{
			Action: breakglass.BG_REQUEST_ACCESS,
			Body:   json.RawMessage(`{"content": "dood"}`),
		},
		client.StartWorkflowOptions{
			ID:        breakglass.TestWFID,
			TaskQueue: breakglass.TaskQueue,
		},
		breakglass.Workflow,
	)
	fmt.Println("WID:", wfr.GetID(), " RID:", wfr.GetRunID())
	//serr := c.SignalWorkflow(context.Background(), wfID, "",
	//	teleport.SignalName,
	//	teleport.BreakGlassSignal{
	//		Action: teleport.BG_REQUEST_ACCESS,
	//		Body:   json.RawMessage(`{"content": "dood"}`),
	//	})
	if serr != nil {
		fmt.Println("ERR: BG_REQUEST_ACCESS")
		spew.Dump(serr)
		fmt.Println("ERR_SIG: ", serr)
	}
	// SImulate approve; reject ..
	time.Sleep(time.Second * 1)
	serr = c.SignalWorkflow(context.Background(), wfID, "",
		breakglass.SignalName,
		breakglass.Signal{
			Action: breakglass.BG_REQUEST_APPROVED,
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
		breakglass.SignalName,
		breakglass.Signal{
			Action: breakglass.BG_OPS_DUMP,
			Body:   nil,
		})
	if serr != nil {
		fmt.Println("ERR_SIG: ", serr)
	}

	return 0
}
