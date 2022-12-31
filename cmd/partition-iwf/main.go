package main

import (
	partition_iwf "app/internal/partition-iwf"
	basic "app/internal/partition-iwf/workflows/basic"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iworkflowio/iwf-golang-sdk/iwf"
	"log"
	"net/http"
	"sync"
)

func main() {
	fmt.Println("Partition Implementation with iWF for FSM ..")
	Run()
}

func Run() {
	fmt.Println("start running samples")
	closeFn := startWorkflowWorker()
	// TODO improve the waiting with process signal
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	closeFn()
}

func startWorkflowWorker() (closeFunc func()) {
	router := gin.Default()
	router.POST(iwf.WorkflowStateStartApi, partition_iwf.ApiV1WorkflowStateStart)
	router.POST(iwf.WorkflowStateDecideApi, partition_iwf.ApiV1WorkflowStateDecide)
	//
	//
	router.GET(
		"/partition/start",
		partition_iwf.StartWorkflow(&basic.PartitionWorkflow{}, "SteadyState", 1),
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
