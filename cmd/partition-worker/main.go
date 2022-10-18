package main

import (
	"app/internal/partition"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var pgd partition.PostgresDB

func main() {

	c, err := client.Dial(client.Options{
		Namespace: "default",
	})
	if err != nil {
		panic(err)
	}
	w := worker.New(c, partition.TaskQueue, worker.Options{
		EnableSessionWorker: true,
	})

	// Start Workflow concurrently .. after registration ..
	go startWorkflow(c)

	// Setup the PGD connectivity
	connString := "postgres://s2admin:password@127.0.0.1:5432/myterraform"
	connConfig, perr := pgx.ParseConfig(connString)
	if perr != nil {
		panic(perr)
	}
	pgd := partition.PostgresDB{ConnConfig: connConfig}
	//w.RegisterWorkflow(teleport.BreakGlassWorkflow)
	w.RegisterWorkflow(partition.LifeCycleWorkflow)
	w.RegisterActivity(&pgd) // If register the whole struct; all the  method + data are passed
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
		ID:        partition.TestWFID,
		TaskQueue: partition.TaskQueue,
	},
		partition.LifeCycleWorkflow,
	)
	if err != nil {
		fmt.Println("Unexpected ERR: ", err)
		return
	}
	fmt.Println("WID:", wfr.GetID(), " RID:", wfr.GetRunID())
	return
}
