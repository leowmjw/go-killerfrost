package breakglass

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go.temporal.io/sdk/workflow"
)

type WorkAction int

const (
	BG_UNKNOWN WorkAction = iota
	BG_OPS_RESTORE
	BG_OPS_DUMP
	BG_REQUEST_ACCESS
	BG_REQUEST_APPROVED
	BG_REQUEST_REJECTED
)

type Signal struct {
	Action WorkAction
	Body   json.RawMessage
}

type Status int

const (
	BGS_INITIAL Status = iota
	BGS_PENDING
	BGS_APPROVED
)

type State struct {
	ID     string // OIDC Username == Github username/identity?
	Status Status
}

// RequestsWorkflow WorkflowID by Team?
func RequestsWorkflow(ctx workflow.Context) error {
	// Starts BreakGlassWorkflow ..
	// Register action; in this case is Postgres; next time can be
	//	MongoDB, Snowflake, Redshift etc.
	var bgs State
	bgs.ID = workflow.GetInfo(ctx).WorkflowExecution.ID
	// DEBUG
	//spew.Dump(bgs)

	// Restore from state from previous Dump to continue
	// TODO: A BreakGlassReq coming in; considered as restore?

	var bgsig Signal
	recv := workflow.GetSignalChannel(ctx, "breakglass")
	for {
		more := recv.Receive(ctx, &bgsig)
		if !more {
			return fmt.Errorf("Unexpected channel closed! State: %v", bgs)
		}
		// DEBUG; content of Signal Received
		//spew.Dump(bgsig)

		// First cut StateMachine:
		// 	BGS_INITIAL --> SIG<BG_REQUEST_ACCESS> --> BGS_PENDING
		//		BGS_PENDING --> SIG<BG_REQUEST_APPROVED> --> BGS_APPROVED
		//		BGS_PENDING --> SIG<BG_REQUEST_REJECTED> --> BGS_INITIAL
		//			BGS_APPROVED --> TimerDone --> BGS_INITIAL
		// Gets signal for RequestAccess

		switch bgsig.Action {
		// Depending on the current state machine; some signals may be ignored!!

		// Gets signal to Approve; including the request timing
		// validate the requested timing ..
		// Gets signal to Reject

		case BG_OPS_DUMP:
			fmt.Println("Dumping state .... finish flow ..")
			spew.Dump(bgs)
			return nil

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
			// With an Approve; have a closing timer
			// Finish; continue as new
			fmt.Println("Now call timer ....")
			// Call the waiting handler ..
		default:
			fmt.Println("BAD SIG", bgsig.Action, " ignoring ..")
			continue

		}

	}

	return nil
}
