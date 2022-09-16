package breakglass

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go.temporal.io/sdk/workflow"
	"time"
)

type WorkAction int

const (
	BG_NOOPS WorkAction = iota
	BG_OPS_RESTORE
	BG_OPS_DUMP
	BG_REQUEST_ACCESS
	BG_REQUEST_APPROVED
	BG_REQUEST_REJECTED
)

const (
	TaskQueue  = "breakglass.queue"
	SignalName = "breakglass"
	TestWFID   = "backend"
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

// Workflow WorkflowID will be the OIDC identity
func Workflow(ctx workflow.Context) error {
	var bgs State
	bgs.ID = workflow.GetInfo(ctx).WorkflowExecution.ID
	// DEBUG
	//spew.Dump(bgs)

	// Restore from state from previous Dump to continue
	// TODO: A BreakGlassReq coming in; considered as restore?
	connString := "postgres://foo:password@127.0.0.1:5432/myterraform" // Get from env?
	nerr, b := NewBastion(connString)
	if nerr != nil {
		return nerr
	}
	// DEBUG
	//fmt.Println("First time init ...")
	//spew.Dump(b.PGConnConfig)

	var bgsig Signal
	recv := workflow.GetSignalChannel(ctx, SignalName)
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

		case BG_NOOPS:
			fmt.Println("NoOPS .. nothing to see here ..")
			return nil
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
			// DEBUG
			//spew.Dump(b.PGConnConfig)

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
			fmt.Println("APPROVED ..")
			// DEBUG
			//spew.Dump(b.PGConnConfig)

			// Extract userName + roleName from signaml..
			fmt.Println("RAW_DATA received .. ")
			spew.Dump(bgsig.Body)
			alcerr := ApprovedLifeCycle(ctx, b, "backend", "s2read")
			if alcerr != nil {
				return alcerr
			}
		default:
			fmt.Println("BAD SIG", bgsig.Action, " ignoring ..")
			continue

		}

	}

	return nil
}

// ApprovedLifeCycle just factor out what happens once approval given; cleaner ..
func ApprovedLifeCycle(ctx workflow.Context, b Bastion, userName, roleName string) error {
	fmt.Println("Inside ApprovedLifeCycle .. ")
	// DEBUG
	//spew.Dump(b.PGConnConfig)

	fmt.Println("TIME_START: ", workflow.Now(ctx))
	ao := workflow.ActivityOptions{
		TaskQueue:           TaskQueue,
		StartToCloseTimeout: time.Second,
	}
	ctx2, cancel := workflow.WithCancel(ctx)
	// TODO: Handle cancel??
	defer cancel()
	ctx2 = workflow.WithActivityOptions(ctx2, ao)

	// With an Approve; have a closing timer
	addf := workflow.ExecuteActivity(ctx2, b.AddToRole, userName, roleName)
	adderr := addf.Get(ctx, nil)
	if adderr != nil {
		fmt.Println("ERR: AddToRole")
		spew.Dump(adderr)
		return adderr
	}
	fmt.Println("TIME_AFTER_ASSIGNED: ", workflow.Now(ctx))
	// Sleep now ..
	fmt.Println("Now call timer ....")
	// Call the waiting handler ..
	serr := workflow.Sleep(ctx, time.Second*15)
	if serr != nil {
		// is this fatal?
		return serr
	}
	f := workflow.ExecuteActivity(ctx2, b.RemoveFromRole, userName, roleName)
	err := f.Get(ctx, nil)
	if err != nil {
		spew.Dump(err)
		return err
	}
	fmt.Println("TIME_AFTER_REVOKED: ", workflow.Now(ctx))

	return nil
}
