package basic

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/iworkflowio/iwf-golang-sdk/gen/iwfidl"
	"github.com/iworkflowio/iwf-golang-sdk/iwf"
	"github.com/jackc/pgx/v5"
	"time"
)

type PostgresDB struct {
	ConnConfig *pgx.ConnConfig
	Conn       *pgx.Conn
}

// Steady state is when the partitions are matching approved plan;
//	until next cycle of evaluation

type SteadyState struct {
	db string
}

func (b SteadyState) GetStateId() string {
	return "SteadyState"
}

func (b SteadyState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	fmt.Println("STEADY_START ===========")
	if b.db == "" {
		fmt.Println("ONe time DB in Steady ..")
		b.db = "SETUP.."
	}

	//var i int
	//err := input.Get(&i)
	//if err != nil {
	//	return nil, err
	//}
	//var tt any
	//gerr := persistence.GetDataObject("TrackedTables", &tt)
	//if gerr != nil {
	//	panic(gerr)
	//}
	//// First time; setup as per init .. probably ..
	//if v, ok := tt.(TrackedTables); ok {
	//	fmt.Println("Should NOT be here!!")
	//	spew.Dump(v)
	//} else {
	//	fmt.Println("SETTING up TrackedTables")
	//	ttt := TrackedTables{
	//		ctx.GetWorkflowId(),
	//		[]Partition{},
	//		[]Partition{},
	//		[]Partition{},
	//	}
	//	persistence.SetDataObject("TrackedTables", ttt)
	//}
	// local is to pass between STart + Decide of same state .. DB connections?
	//persistence.GetStateLocal()
	var tt TrackedTables
	gerr := persistence.GetDataObject("TrackedTables", &tt)
	if gerr != nil {
		panic(gerr)
	}
	// Will be zero-valued ..
	if tt.ID == "" && tt.Planned == nil {
		tt.ID = ctx.GetWorkflowId()
		tt.Planned = []Partition{{Statement: "SELECT * FROM bob"}}
	}
	serr := persistence.SetDataObject("TrackedTables", tt)
	if serr != nil {
		panic(serr)
	}
	// Either get signal or let 10 secs exoire ..
	return iwf.AnyCommandCompletedRequest(
		iwf.NewSignalCommand("signal", SignalName),
		iwf.NewInterStateChannelCommand("intSignal", SignalName),
		iwf.NewTimerCommand("check", time.Now().Add(60*time.Second)),
	), nil
}

func (b SteadyState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	fmt.Println("STEADY_DECIDE ===========")
	if b.db == "" {
		fmt.Println("ERROR: SHOULD be named!!")
	} else {
		fmt.Println("VAI STRUCT: ", b.db)
	}
	//var i int
	//err := input.Get(&i)
	//if err != nil {
	//	return nil, err
	//}
	extSig := commandResults.GetSignalCommandResultByChannel(SignalName)
	// DEBUG
	//fmt.Println("APPROVE_EXT_SIGNAL .. ==========>")
	//spew.Dump(extSig)
	if extSig.Status == iwfidl.RECEIVED {
		var ptSignal PTSignal
		gerr := extSig.SignalValue.Get(&ptSignal)
		if gerr != nil {
			return nil, gerr
		}
		// DEBUG
		//fmt.Println("EXT_SIGNAL_VAL:")
		//spew.Dump(ptSignal)
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
	// Get data ..
	//var tt any
	//gerr := persistence.GetDataObject("TrackedTables", &tt)
	//if gerr != nil {
	//	panic(gerr)
	//}
	//// First time; setup as per init .. probably ..
	//if v, ok := tt.(TrackedTables); ok {
	//	fmt.Println("STEADY DECIDE")
	//	fmt.Println("========================")
	//	spew.Dump(v)
	//} else {
	//	fmt.Println("Should NOT be here!! DECIDE")
	//	spew.Dump(tt)
	//}
	var tt TrackedTables
	gerr := persistence.GetDataObject("TrackedTables", &tt)
	if gerr != nil {
		panic(gerr)
	}
	// Will be zero-valued ..
	if tt.ID == "" {
		fmt.Println("Should NOTz be here!! DECIDE")
	} else {
		fmt.Println("INIT:", tt.ID)
		spew.Dump(tt.Planned)
	}
	// See if the timer was fired ...
	return iwf.SingleNextState("ApprovalState", nil), nil
}

func (b SteadyState) GetStateOptions() *iwfidl.WorkflowStateOptions {
	// Unsure when this is used ..
	fmt.Println("Inside SteadyState GetStateOptions ==>  ")
	v, err := iwfidl.NewWorkflowStateOptionsWithDefaults().MarshalJSON()
	if err != nil {
		panic(err)
	}
	spew.Dump(v)
	//iwfidl.NewNullableWorkflowStateOptions(nil)

	return nil
}
