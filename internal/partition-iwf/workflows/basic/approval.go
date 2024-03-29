package basic

import (
	"fmt"
	"github.com/iworkflowio/iwf-golang-sdk/gen/iwfidl"
	"github.com/iworkflowio/iwf-golang-sdk/iwf"
	"time"
)

const (
	ApprovalStateID = "ApprovalState"
)

// Approval state happens when a gap is detected
//	for simple case it may auto-advance to Diff state

type ApprovalState struct{}

func (b ApprovalState) GetStateId() string {
	return ApprovalStateID
}

func (b ApprovalState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	fmt.Println("APPROVAL_START ===========")
	//var i int
	//err := input.Get(&i)
	//if err != nil {
	//	return nil, err
	//}
	// Setup local DB connection ...

	//var db PostgresDB
	//// Use it to store the schema + tables instead??
	//gerr := persistence.GetStateLocal("DB", &db)
	//if gerr != nil {
	//	panic(gerr)
	//}
	//if db.Conn == nil {
	//	// Setup the PGD connectivity the first time ..
	//	fmt.Println("ONLY ONCE for Approval!!")
	//	connString := "postgres://s2admin:password@127.0.0.1:5432/myterraform"
	//	connConfig, perr := pgx.ParseConfig(connString)
	//	if perr != nil {
	//		fmt.Println("PERR:", perr)
	//		panic(perr)
	//	}
	//	// Avoid attaching first ..
	//	//db.ConnConfig = connConfig
	//	spew.Dump(connConfig)
	//	// All checks out ..
	//	c, cerr := pgx.ConnectConfig(context.Background(), connConfig)
	//	if cerr != nil {
	//		fmt.Println("CERR:", cerr)
	//		panic(cerr)
	//	}
	//	//db.Conn = c
	//	spew.Dump(c)
	//	// Store it for use later ,..
	//	//serr := persistence.SetStateLocal("DB", db)
	//	//if serr != nil {
	//	//	panic(serr)
	//	//}
	//}
	// Either get signal or let 10 secs exoire ..
	return iwf.AnyCommandCompletedRequest(
		iwf.NewSignalCommand("signal", SignalName),
		iwf.NewInterStateChannelCommand("intSignal", SignalName),
		iwf.NewTimerCommand("check", time.Now().Add(30*time.Second)),
	), nil
}

func (b ApprovalState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	fmt.Println("APPROVAL_DECIDE")
	//var i int
	//err := input.Get(&i)
	//if err != nil {
	//	return nil, err
	//}
	//var db PostgresDB
	//gerr := persistence.GetStateLocal("DB", &db)
	//if gerr != nil {
	//	panic(gerr)
	//}
	//// CHeck can see; mayb eeven run wuery??
	//spew.Dump(db.Conn)

	// Start Signals
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
	// DEBUG
	//fmt.Println("APPROVE_INT_SIGNAL .. ==========>")
	//spew.Dump(intSig)
	if intSig.Status == iwfidl.RECEIVED {
		var ptSignal PTSignal
		gerr := intSig.Value.Get(&ptSignal)
		if gerr != nil {
			return nil, gerr
		}
		switch ptSignal.Action {
		case PT_OPS_DUMP:
			// Save state ..
			fmt.Println("GOT DUMP in ApprovalState")
			return iwf.GracefulCompleteWorkflow("WEIRD"), nil
		default:
			fmt.Println("UNKNOWNN ...")
		}
	}
	// Normsl route ..
	// Approve --> Diff
	// Reject --> SteadyState ..
	return iwf.SingleNextState("Diff", nil), nil
	// Below is safe; callobarate
	//return iwf.GracefulCompleteWorkflow("NORMAL"), nil
}

func (b ApprovalState) GetStateOptions() *iwfidl.WorkflowStateOptions {
	return nil
}
