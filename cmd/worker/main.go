package main

import (
	"app/internal/breakglass"
	"app/internal/teleport"
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var b breakglass.Bastion

func main() {

	c, err := client.Dial(client.Options{
		Namespace: "default",
	})
	if err != nil {
		panic(err)
	}
	w := worker.New(c, teleport.TaskQueue, worker.Options{
		EnableSessionWorker: true,
	})

	// Start Workflow concurrently ..
	go startWorkflow(c)

	// How to avoid this; useless ..
	berr, b := breakglass.NewBastion("postgres://foo:password@127.0.0.1:5432/myterraform")
	if berr != nil {
		panic(berr)
	}

	w.RegisterWorkflow(teleport.BreakGlassWorkflow)
	w.RegisterActivity(&b) // If register the whole struct; all the  method + data are passed
	// Below is how would be done traditionally
	//w.RegisterActivity(b.AddToRole)
	//w.RegisterActivity(b.RemoveFromRole)

	rerr := w.Run(worker.InterruptCh())
	if rerr != nil {
		panic(rerr)
	}
}

func startWorkflow(c client.Client) {
	wfr, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:        teleport.TestWFID,
		TaskQueue: teleport.TaskQueue,
	},
		teleport.BreakGlassWorkflow,
	)
	if err != nil {
		fmt.Println("Unexpected ERR: ", err)
		return
	}
	fmt.Println("WID:", wfr.GetID(), " RID:", wfr.GetRunID())
	return
}
