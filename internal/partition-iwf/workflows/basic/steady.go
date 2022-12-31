package basic

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/iworkflowio/iwf-golang-sdk/gen/iwfidl"
	"github.com/iworkflowio/iwf-golang-sdk/iwf"
	"time"
)

// Steady state is when the partitions are matching approved plan;
//	until next cycle of evaluation

type SteadyState struct{}

func (b SteadyState) GetStateId() string {
	return "SteadyState"
}

func (b SteadyState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	//var i int
	//err := input.Get(&i)
	//if err != nil {
	//	return nil, err
	//}
	// Either get signal or let 10 secs exoire ..
	return iwf.AnyCommandCompletedRequest(
		iwf.NewSignalCommand("signal", SignalName),
		iwf.NewInterStateChannelCommand("intSignal", SignalName),
		iwf.NewTimerCommand("check", time.Now().Add(60*time.Second)),
	), nil
}

func (b SteadyState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	//var i int
	//err := input.Get(&i)
	//if err != nil {
	//	return nil, err
	//}
	extSig := commandResults.GetSignalCommandResultByChannel(SignalName)
	fmt.Println("APPROVE_EXT_SIGNAL .. ==========>")
	spew.Dump(extSig)
	if extSig.Status == iwfidl.RECEIVED {
		var ptSignal PTSignal
		gerr := extSig.SignalValue.Get(&ptSignal)
		if gerr != nil {
			return nil, gerr
		}
		fmt.Println("EXT_SIGNAL_VAL:")
		spew.Dump(ptSignal)
	}

	intSig := commandResults.GetInterStateChannelCommandResultByChannel(SignalName)
	if intSig.Status == iwfidl.RECEIVED {

		var sigAction PTAction
		gerr := intSig.Value.Get(&sigAction)
		if gerr != nil {
			return nil, gerr
		}
		switch sigAction {
		case PT_OPS_DUMP:
			// Save state ..
			fmt.Println("GOT DUMP in SteadyState")
			return iwf.GracefulCompleteWorkflow(1), nil
		default:
			fmt.Println("UNKNOWNN ...")

		}
	}
	// Send signal to dump on stuff for the next state ..
	//communication.PublishInterstateChannel(SignalName, PTSignal{
	//	Action: PT_OPS_DUMP,
	//})
	// See if the timer was fired ...
	return iwf.SingleNextState("ApprovalState", nil), nil
}

func (b SteadyState) GetStateOptions() *iwfidl.WorkflowStateOptions {
	return nil
}
