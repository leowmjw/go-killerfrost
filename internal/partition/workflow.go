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
	PT_REQUEST_ACCESS
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
	PT_PENDING
	PT_APPROVED
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

	// If empty; and not loading from a previous state
	// Setup X number of SmallestUnit
	// Set current SmallestUnit as Archival
	trackedTable := TrackedTable{
		Schema: "s2",
		Name:   "foobar",
	}

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
			return nil
		case PT_OPS_DUMP:
			spew.Dump(p)
			fmt.Println("Dumping state .... finish flow ..")
			spew.Dump(pts)
			return nil
		case PT_REQUEST_ACCESS:
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

// determineNextCheckPoint finds next micro boundary ; if passed do now
func determineNextCheckPoint(currentTime time.Time) (time.Time, error) {
	// if currentTime > Now; fire it off immediately ..
	// Dummy vakue .. if passed; do now
	return time.Now(), nil
}
