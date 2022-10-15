package partition

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go.temporal.io/sdk/workflow"
	"time"
)

type PTAction int

const (
	PT_NOOPS PTAction = iota
	PT_OPS_RESTORE
	PT_OPS_DUMP
	PT_REQUEST_PLAN
	PT_REQUEST_APPROVED
	PT_REQUEST_REJECTED
)

const (
	TaskQueue  = "partition.queue"
	SignalName = "partition"
	TestWFID   = "measurement" // Sample Demo Table used to show PG Partition
)

type PTSignal struct {
	Action PTAction
	Body   json.RawMessage
}

// PTStatus not used by .. remove for now?
type PTStatus int

const (
	PT_INITIAL PTStatus = iota
	PT_PENDING          // Plan Pending Approval
	PT_LOCKED           // Plan approved; locked
	PT_HASPLAN
	PT_APPLIED
	PT_ARCHIVED
)

type PTState struct {
	ID     string // OIDC Username == Github username/identity?
	Status PTStatus
}

// Complement a coplete solution like:
//	ByteBase - https://www.bytebase.com/docs/introduction/what-is-bytebase

// Another attarctive option; supports partition:
//  AtlasGo - https://atlasgo.io/atlas-schema/sql-resources#partitions
// Other options; https://github.com/amacneil/dbmate

// LifeCycleWorkflow keeps track of all the partition slices as it goes through life
func LifeCycleWorkflow(ctx workflow.Context) error {
	// TODO: Passed in via LifeCycleReq?
	// If empty; and not loading from a previous state
	// Setup X number of SmallestUnit
	// Set current SmallestUnit as Archival
	trackedTable := TrackedTable{
		Schema: "s2",
		Name:   TestWFID,
	}
	partitionedTable := PTState{
		ID:     TestWFID,
		Status: PT_INITIAL,
	}
	// TODO: Need to check that the target table is partition and acceptable .
	// QUERY: For UI + actions ..
	qerr := workflow.SetQueryHandler(ctx, SignalName, func() (string, error) {
		ptStatus := partitionedTable.Status
		switch ptStatus {
		case PT_HASPLAN:
			return "plan", nil
		case PT_LOCKED:
			return "locked", nil
		default:
			fmt.Println("UNEXPECTED: ", ptStatus)
		}

		return "unknown", nil
	})
	if qerr != nil {
		fmt.Println("ERR: ", qerr)
		spew.Dump(partitionedTable)
		return qerr
	}
	// Info from LifeCycleReq?
	type partition_slices struct {
		Name string
		Status
	}
	var p []partition_slices

	// Add a new porcesing ... until spare of x WAITING
	p = append(p, partition_slices{
		Name:   "bobo",
		Status: WAITING,
	})

	// If x < len(WAITING); continue on ..

	// Slice is ATTACHED; updte vi sideeffect?
	// Slice is DETACHED
	// Slice is ARCHIVED

	//heap.Push(p, partition_slices{
	//	Name:   "foofoo",
	//	Status: INUSE,
	//})
	//heap.Pop(nil)

	if !(len(trackedTable.Ranges) > 0) {
		// Edge case; first time just set things up ..
	}

	// Every microslice of time;
	// Supported Partition
	// Test: Day -> Year; Hours -> Month; Min -> Day
	// Production: Year; Month; Day
	// Start with fixed; next time can be flexible; e.g. workday, weekend, holidays ..

	//ok, err := workflow.AwaitWithTimeout(ctx, time.Minute, func() bool {
	//	return len(trackedTable.Ranges) == 5
	//})
	//if err != nil {
	//	panic(err)
	//}
	//if ok {
	//	fmt.Println("Issat OK??")
	//}

	// Spawn the timer to check back on the last SmallestUnit (..race..)
	// Look for signal to do checks too .. on-demand ..
	// Look for sugnal to bring back sanity in case of missing; with a plan ..

	workflow.Go(ctx, func(ctx workflow.Context) {
		for {
			checkerLoop(ctx, &trackedTable)
		}
	})
	// Loop through the actives; mark those detach + those pending

	// Loop through those marked for detach; fire off activities in concurrent
	// Loop through those marked pending; fire off activities in concurrent ..

	var pts PTState
	var ptsig PTSignal
	recv := workflow.GetSignalChannel(ctx, SignalName)
	for {
		more := recv.Receive(ctx, &ptsig)
		if !more {
			return fmt.Errorf("Unexpected channel closed! PTState: %v", pts)
		}
		// DEBUG; content of Signal Received
		//spew.Dump(ptsig)

		// First cut StateMachine:
		// 	PAS_INITIAL --> SIG<PA_REQUEST_ACCESS> --> PAS_PENDING
		//		PAS_PENDING --> SIG<PA_REQUEST_APPROVED> --> PAS_APPROVED
		//		PAS_PENDING --> SIG<PA_REQUEST_REJECTED> --> PAS_INITIAL
		//			PAS_APPROVED --> TimerDone --> PAS_INITIAL
		// Gets signal for RequestAccess

		switch ptsig.Action {
		// Depending on the current state machine; some signals may be ignored!!

		// Gets signal to Approve; including the request timing
		// validate the requested timing ..
		// Gets signal to Reject

		case PT_NOOPS:
			fmt.Println("NoOPS .. nothing to see here ..")
			continue
		case PT_OPS_DUMP:
			spew.Dump(p)
			fmt.Println("Dumping state .... finish flow ..")
			spew.Dump(pts)
			return nil
		case PT_REQUEST_PLAN:
			//if pts.Status != PT_INITIAL {
			//	fmt.Println("BAD SIG PA_REQUEST_ACCESS for STATUS: ", bgs.Status, " ignoring ..")
			//	continue
			//}
			//pts.Status = PT_PENDING
			fmt.Println("Set status to PENDING")
			// DEBUG
			//spew.Dump(b.PGConnConfig)

			continue
		case PT_REQUEST_REJECTED:
			//if pts.Status != PT_PENDING {
			//	fmt.Println("BAD SIG PA_REQUEST_REJECTED for STATUS: ", bgs.Status, " ignoring ..")
			//	continue
			//}
			fmt.Println("Rejected .. rset the state")
			//pts.Status = PT_INITIAL
		case PT_REQUEST_APPROVED:
			//if pts.Status != PT_PENDING {
			//	fmt.Println("BAD SIG PA_REQUEST_APPROVED for STATUS: ", bgs.Status, " ignoring ..")
			//	continue
			//}
			fmt.Println("APPROVED ..")
			// DEBUG
			// Extract userName + roleName from signaml..
			fmt.Println("RAW_DATA received .. ")
			spew.Dump(ptsig.Body)
			//alcerr := ApprovedLifeCycle(ctx, b, "backend", "s2read")
			//if alcerr != nil {
			//	return alcerr
			//}
		default:
			fmt.Println("BAD SIG", ptsig.Action, " ignoring ..")
			continue

		}

	}

	return nil
}

func checkerLoop(ctx workflow.Context, tt *TrackedTable) {
	// Find out the last run? Calculate currenttime against it?
	// Int test: boundary == Minute; in prod: boundary == Day
	serr := workflow.Sleep(ctx, time.Minute)
	if serr != nil {
		fmt.Println("CANCELED: ", serr)
		return
	}
	// TODO: See if there is a gap; only when in steady state

	// NOTE: This just mark things up; it will then appear
	// in State and will need human intervention? as query?
	timeNow := workflow.Now(ctx)
	// DEBUG
	//fmt.Println("TIME: ", timeNow)
	//spew.Dump(tt)
	// If want to change like below; need to pass in ref ..
	tt.Name = "boo" + timeNow.String()
	// Below not needed but if want something more sophisticated..
	//var nextRunTime time.Time
	//encodedNextRun := workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
	//	return determineNextCheckPoint(workflow.Now(ctx))
	//
	//})
	//encodedNextRun.Get(&nextRunTime)
}

// determineNextCheckPoint finds next micro boundary ; if passed do now
func determineNextCheckPoint(currentTime time.Time) time.Time {
	// if currentTime > Now; fire it off immediately ..
	// Dummy vakue .. if passed; do now
	return time.Now()
}
