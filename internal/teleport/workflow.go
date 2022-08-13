package teleport

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go.temporal.io/sdk/workflow"
)

type BreakGlassWorkAction int

const (
	BG_UNKNOWN BreakGlassWorkAction = iota
	BG_OPS_RESTORE
	BG_OPS_DUMP
	BG_REQUEST_ACCESS
	BG_REQUEST_APPROVED
	BG_REQUEST_REJECTED
)

type BreakGlassSignal struct {
	Action BreakGlassWorkAction
	Body   json.RawMessage
}

type BreakGlassStatus int

const (
	BGS_INITIAL BreakGlassStatus = iota
	BGS_PENDING
	BGS_APPROVED
)

type BreakGlassState struct {
	ID     string // OIDC Username == Github username/identity?
	Status BreakGlassStatus
}

// BreakGlassWorkflow WorkflowID will be the OIDC identity
func BreakGlassWorkflow(ctx workflow.Context) error {
	var bgs BreakGlassState
	bgs.ID = workflow.GetInfo(ctx).WorkflowExecution.ID
	// DEBUG
	//spew.Dump(bgs)
	var bgsig BreakGlassSignal
	recv := workflow.GetSignalChannel(ctx, "breakglass")

	// Below triggers workflow error - ScheduleToClose
	//for recv.Receive(ctx, &bgsig) {}

	for {
		more := recv.Receive(ctx, &bgsig)
		if !more {
			return fmt.Errorf("Unexpected channel closed! State: %v", bgs)
		}
		spew.Dump(bgsig)
		// Restore from state from previous Dump to continue
		// Receive Operational Signal to Dump state and close
		if bgsig.Action == BG_OPS_DUMP {
			fmt.Println("Dumping state .... finish flow ..")
			spew.Dump(bgs)
			return nil
		}

		// First cut StateMachine:
		// 	BGS_INITIAL --> SIG<BG_REQUEST_ACCESS> --> BGS_PENDING
		//		BGS_PENDING --> SIG<BG_REQUEST_APPROVED> --> BGS_APPROVED
		//		BGS_PENDING --> SIG<BG_REQUEST_REJECTED> --> BGS_INITIAL
		//			BGS_APPROVED --> TimerDone --> BGS_INITIAL
		// Gets signal for RequestAccess

		switch bgsig.Action {
		case BG_REQUEST_ACCESS:
			if bgs.Status != BGS_INITIAL {
				fmt.Println("BAD SIG BG_REQUEST_ACCESS for STATUS: ", bgs.Status, " ignoring ..")
				continue
			}
			bgs.Status = BGS_PENDING
			fmt.Println("Set status to PENDING")
			continue
		case BG_REQUEST_REJECTED:
			if bgs.Status != BGS_PENDING {
				fmt.Println("BAD SIG BG_REQUEST_REJECTED for STATUS: ", bgs.Status, " ignoring ..")
				continue
			}
			fmt.Println("Rejected .. rset the state")
			bgs.Status = BGS_INITIAL
		case BG_REQUEST_APPROVED:
			if bgs.Status != BGS_PENDING {
				fmt.Println("BAD SIG BG_REQUEST_APPROVED for STATUS: ", bgs.Status, " ignoring ..")
				continue
			}
			fmt.Println("Now call timer ....")
		default:
			fmt.Println("BAD SIG", bgsig.Action, " ignoring ..")
			continue

		}
		// Depending on the current state machine; some signals may be ignored!!

		// Gets signal to Approve; including the request timing
		// validate the requested timing ..
		// Gets signal to Reject

		// With an Approve; have a closing timer
		// Finish; continue as new

	}

	return nil
}
