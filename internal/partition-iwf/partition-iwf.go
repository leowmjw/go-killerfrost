package partition_iwf

import (
	"app/internal/partition-iwf/workflows"
	"app/internal/partition-iwf/workflows/basic"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iworkflowio/iwf-golang-sdk/gen/iwfidl"
	"github.com/iworkflowio/iwf-golang-sdk/iwf"
	"net/http"
)

var client = iwf.NewClient(workflows.GetRegistry(), nil)
var workerService = iwf.NewWorkerService(workflows.GetRegistry(), nil)

// StartWorkflow Calls the commands using script ..
func StartWorkflow(wf iwf.Workflow, startStateId string, input interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		wfId := "partition-mleow"
		runId, err := client.StartWorkflow(c.Request.Context(), wf, startStateId, wfId, 3600, input, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, fmt.Sprintf("workflowId: %v runId: %v", wfId, runId))
		return
	}
}

// SignalWorkflow symbolizes acceptance ..
func SignalWorkflow(wf iwf.Workflow, signal basic.PTSignal) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// DEBUG
		//spew.Dump(ctx)
		err := client.SignalWorkflow(ctx, wf, "partition-mleow", "", basic.SignalName, signal)
		if err != nil {
			fmt.Println("ERR:", err)
		}
	}
}

func ApiV1WorkflowStateStart(c *gin.Context) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := workerService.HandleWorkflowStateStart(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

func ApiV1WorkflowStateDecide(c *gin.Context) {
	var req iwfidl.WorkflowStateDecideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := workerService.HandleWorkflowStateDecide(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}
