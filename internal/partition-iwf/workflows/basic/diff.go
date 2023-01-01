package basic

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/iworkflowio/iwf-golang-sdk/gen/iwfidl"
	"github.com/iworkflowio/iwf-golang-sdk/iwf"
)

// Diff State is when partitions are not matching ideal state
//	goes to Steady once confirmed all matching plan is done ..

// basic skeleton below; replace Diff with your own ..

const (
	DiffID = "Diff"
)

type DiffState struct{}

func (b DiffState) GetStateId() string {
	return DiffID
}

func (b DiffState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	return nil, nil
}

func (b DiffState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	//var i int
	return nil, nil
}

func (b DiffState) GetStateOptions() *iwfidl.WorkflowStateOptions {
	fmt.Println("Inside GetStateOptions ==>  ")
	v, err := iwfidl.NewWorkflowStateOptionsWithDefaults().MarshalJSON()
	if err != nil {
		panic(err)
	}
	spew.Dump(v)
	//iwfidl.NewNullableWorkflowStateOptions(nil)
	return nil
}
