package teleport

import (
	"app/internal/breakglass"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go.temporal.io/sdk/workflow"
	"time"
)

type ProxyAccessWorkAction int

const (
	PA_NOOPS ProxyAccessWorkAction = iota
	PA_OPS_RESTORE
	PA_OPS_DUMP
	PA_REQUEST_ACCESS
	PA_REQUEST_APPROVED
	PA_REQUEST_REJECTED
)

const (
	TaskQueue  = "proxyaccess.queue"
	SignalName = "proxyaccess"
	TestWFID   = "xxxyy"
)

type ProxyAccessSignal struct {
	Action ProxyAccessWorkAction
	Body   json.RawMessage
}

type ProxyAccessStatus int

const (
	PAS_INITIAL ProxyAccessStatus = iota
	PAS_PENDING
	PAS_APPROVED
)

type ProxyAccessState struct {
	ID     string // OIDC Username == Github username/identity?
	Status ProxyAccessStatus
}

// ProxyAccessWorkflow WorkflowID will be the OIDC identity
func ProxyAccessWorkflow(ctx workflow.Context) error {
	var bgs ProxyAccessState
	bgs.ID = workflow.GetInfo(ctx).WorkflowExecution.ID
	// DEBUG
	//spew.Dump(bgs)

	// Restore from state from previous Dump to continue
	// TODO: A BreakGlassReq coming in; considered as restore?
	connString := "postgres://foo:password@127.0.0.1:5432/myterraform" // Get from env?
	nerr, b := breakglass.NewBastion(connString)
	if nerr != nil {
		return nerr
	}
	// DEBUG
	//fmt.Println("First time init ...")
	//spew.Dump(b.PGConnConfig)

	var bgsig ProxyAccessSignal
	recv := workflow.GetSignalChannel(ctx, SignalName)
	for {
		more := recv.Receive(ctx, &bgsig)
		if !more {
			return fmt.Errorf("Unexpected channel closed! State: %v", bgs)
		}
		// DEBUG; content of Signal Received
		//spew.Dump(bgsig)

		// First cut StateMachine:
		// 	PAS_INITIAL --> SIG<PA_REQUEST_ACCESS> --> PAS_PENDING
		//		PAS_PENDING --> SIG<PA_REQUEST_APPROVED> --> PAS_APPROVED
		//		PAS_PENDING --> SIG<PA_REQUEST_REJECTED> --> PAS_INITIAL
		//			PAS_APPROVED --> TimerDone --> PAS_INITIAL
		// Gets signal for RequestAccess

		switch bgsig.Action {
		// Depending on the current state machine; some signals may be ignored!!

		// Gets signal to Approve; including the request timing
		// validate the requested timing ..
		// Gets signal to Reject

		case PA_NOOPS:
			fmt.Println("NoOPS .. nothing to see here ..")
			return nil
		case PA_OPS_DUMP:
			fmt.Println("Dumping state .... finish flow ..")
			spew.Dump(bgs)
			return nil
		case PA_REQUEST_ACCESS:
			if bgs.Status != PAS_INITIAL {
				fmt.Println("BAD SIG PA_REQUEST_ACCESS for STATUS: ", bgs.Status, " ignoring ..")
				continue
			}
			bgs.Status = PAS_PENDING
			fmt.Println("Set status to PENDING")
			// DEBUG
			//spew.Dump(b.PGConnConfig)

			continue
		case PA_REQUEST_REJECTED:
			if bgs.Status != PAS_PENDING {
				fmt.Println("BAD SIG PA_REQUEST_REJECTED for STATUS: ", bgs.Status, " ignoring ..")
				continue
			}
			fmt.Println("Rejected .. rset the state")
			bgs.Status = PAS_INITIAL
		case PA_REQUEST_APPROVED:
			if bgs.Status != PAS_PENDING {
				fmt.Println("BAD SIG PA_REQUEST_APPROVED for STATUS: ", bgs.Status, " ignoring ..")
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
func ApprovedLifeCycle(ctx workflow.Context, b breakglass.Bastion, userName, roleName string) error {
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
