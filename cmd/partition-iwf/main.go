package main

import (
	"app/internal/partition"
	partition_iwf "app/internal/partition-iwf"
	basic "app/internal/partition-iwf/workflows/basic"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iworkflowio/iwf-golang-sdk/iwf"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Partition Implementation with iWF for FSM ..")
	Run()
}

func Run() {
	fmt.Println("start running samples")
	closeFn := startWorkflowWorker()
	// We can run normal Temporal workers too ..
	var w worker.Worker
	go func() {
		c, err := client.Dial(client.Options{
			Namespace: "default",
		})
		if err != nil {
			panic(err)
		}
		w = worker.New(c, partition.TaskQueue, worker.Options{
			EnableSessionWorker: true,
		})
		// Will stop later ..
		rerr := w.Start()
		if rerr != nil {
			panic(rerr)
		}
	}()

	// Block till SIGTERM ..
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	s := <-c
	fmt.Println("GOT_SIG:", s.String())
	fmt.Println("UNBLOCKED ..")
	// Clean up iWF
	closeFn()
	// Clean up normal Worker
	w.Stop()
}

func startWorkflowWorker() (closeFunc func()) {
	router := gin.Default()
	router.POST(iwf.WorkflowStateStartApi, partition_iwf.ApiV1WorkflowStateStart)
	router.POST(iwf.WorkflowStateDecideApi, partition_iwf.ApiV1WorkflowStateDecide)
	//
	//
	router.GET("/partition/start",
		partition_iwf.StartWorkflow(&basic.PartitionWorkflow{}, "SteadyState", 1),
	)
	router.GET("/partition/signal",
		partition_iwf.SignalWorkflow(&basic.PartitionWorkflow{}, basic.PTSignal{
			Action: basic.PT_OPS_DUMP,
		}),
	)
	//router.GET("/basic/start", startWorklfow(&basic.BasicWorkflow{}, basic.BasicWorkflowState1Id, 1))
	//router.GET("/interstateChannel/start", startWorklfow(&interstate.InterStateWorkflow{}, interstate.InterStateWorkflowState0Id, nil))
	//router.GET("/persistence/start", startWorklfow(&persistence.PersistenceWorkflow{}, persistence.PersistenceWorkflowState1Id, input))
	//router.GET("/signal/start", startWorklfow(&signal.SignalWorkflow{}, signal.SignalWorkflowState1Id, nil))
	//router.GET("/timer/start", startWorklfow(&timer.TimerWorkflow{}, timer.TimerWorkflowState1Id, 5))

	wfServer := &http.Server{
		Addr:    ":" + iwf.DefaultWorkerPort,
		Handler: router,
	}
	go func() {
		if err := wfServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return func() { wfServer.Close() }
}
