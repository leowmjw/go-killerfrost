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
		IsTest: true,
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
			pts.Status = PT_INITIAL
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

// applyPlan only means an approved signal given; or manual call ..
func applyPlan(ctx workflow.Context, tt *TrackedTable) error {
	// TODO: Sanity check on the DB Conn??
	// Check if the parent table exists; if not create it ..
	ctxact := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10,
	})
	workflow.ExecuteActivity(ctxact, PostgresDB.ArchiveDateRange)

	for _, pt_slice := range tt.Ranges {
		if pt_slice.Status == WAITING {
			workflow.GoNamed(ctx, tt.Name, func(ctx workflow.Context) {

			})
		} else {
			fmt.Println("SLICE: ", pt_slice.Name)
			fmt.Println("STATUS: ", pt_slice.Status)
		}
	}
	return nil
}

func checkerLoop(ctx workflow.Context, tt *TrackedTable) {
	// Find out the last run? Calculate currenttime against it?
	// Int test: boundary == Minute; in prod: boundary == Day
	serr := workflow.Sleep(ctx, time.Minute)
	if serr != nil {
		spew.Dump(serr)
		fmt.Println("CANCELED: ", serr)
		return
	}
	// TODO: See if there is a gap; only when in steady state
	// Always maintain x number of slices in the future

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
	ctxact := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10,
	})
	workflow.ExecuteActivity(ctxact, PostgresDB.ArchiveDateRange)
}

// getMajorTimeBoundary - for Test - Day (31d); for Prod - Year (2xyz)
//func (tt *TrackedTable) getMajorTimeBoundary(currentTime time.Time) int {
//	if tt.IsTest {
//		return currentTime.Day()
//	}
//	// Handle rollover
//	return currentTime.Year()
//}

// getMinorTimeBoundary - for Test - Hour (24h); for Prod - Month (12m)
//func (tt *TrackedTable) getMinorTimeBoundary(currentTime time.Time) int {
//	if tt.IsTest {
//		return currentTime.Hour()
//	}
//	return int(currentTime.Month())
//}

func (tt *TrackedTable) getSliceProjection(ct time.Time, numProjection int) {
	currentMinorIndex := func() int {
		if tt.IsTest {
			return ct.Hour()
		}
		return int(ct.Month())
	}()
	currentMajorIndex := func() int {
		if tt.IsTest {
			return ct.Day()
		}
		return ct.Year()
	}()
	// Iterate till see currentMajor + currentMinor; split at this seam
	// if finish with nothing use the full numProject; if found -1
	var startIndex, seamIndex int
	//var leftOver []DateRange
	if len(tt.Ranges) > 0 {
		labelName := fmt.Sprintf("y%04dm%02d", currentMajorIndex, currentMinorIndex)
		// then iterate through to find the seamIndex
		for si, v := range tt.Ranges {
			if seamIndex > 0 {
				// Found the currentSeam
				// Can check rest of items now ..
				numProjection--
				continue // next loop ..
			}
			// Default is to look for the seam ..
			if v.Name == labelName {
				seamIndex = si
				startIndex = 1
				break
			}
		}
		// then do slice; else it is full range
		fmt.Println("SEAMIDX: ", seamIndex)
		// extract slice sarting wirth seamIndex till len(slice)
		//leftOver = tt.Ranges[seamIndex:len(tt.Ranges)]
		//// DEBUG ...
		//spew.Dump(leftOver)
	}
	// DEBUG
	//fmt.Println("MAJOR: ", currentMajorIndex, " MINOR: ", currentMinorIndex)
	var nr []DateRange
	var labelName string
	// Project numProjection into future ..
	for i := startIndex; i < numProjection; i++ {
		// Rules for the boundary will be different if it is Test ..
		// If those index does not exist; add!
		// If those are WAITING .. do nothing ..
		var finalMajorIndex, finalMinorIndex int
		if tt.IsTest {
			finalMajorIndex = currentMajorIndex
			finalMinorIndex = (currentMinorIndex + i) % 24
			if finalMinorIndex < currentMinorIndex {
				finalMajorIndex = (finalMajorIndex + 1) % 31
			}
		} else {
			finalMajorIndex = currentMajorIndex
			finalMinorIndex = (currentMinorIndex + i) % 12 //Month
			if finalMinorIndex == 0 {
				finalMinorIndex++ // Month starts with 1
			}
			if finalMinorIndex < currentMinorIndex {
				finalMajorIndex++ // Year unlikely to run out for a while ..
			}
		}
		labelName = fmt.Sprintf("y%04dm%02d", finalMajorIndex, finalMinorIndex)
		// Figure iut between test + prod; append if needed
		// Decide if already WAIT/ATTACHED?
		nr = append(nr, DateRange{
			Name:    labelName,
			MinDate: "2022",
			Status:  WAITING,
		})
	}
	// if temp > 0; then can append ..
	if len(nr) > 0 {
		// DEBUG
		//		fmt.Println("APPEND these ...")
		//spew.Dump(nr)
		tt.Ranges = append(tt.Ranges, nr...)
	} else {
		fmt.Println("NOTHING to DO ....")
	}
	return
}

func determineBufferNeeded(ctx workflow.Context, policy Policy, tt *TrackedTable) error {
	// Where are we now; inclusive
	// slice it out of tt.Ranges; what remains?
	var currentTime time.Time
	encodedCurrent := workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
		return workflow.Now(ctx)
	})
	err := encodedCurrent.Get(&currentTime)
	if err != nil {
		fmt.Println("ERR: ", err)
		return err
	}
	// measurement_y2006m02
	// YYYY-MM/DD e.g. 2022-02/30
	// DD-HH/mm e.g. 31-24/55
	boundaryKey := fmt.Sprintf("y%4dm%2d", currentTime.Day(), currentTime.Hour()) // 01-24 // MONTH
	fmt.Println("KEY: ", boundaryKey)
	// Policy will determine unit and how far ahead needed
	// for now we hard code
	bufferPolicy := 3
	slicesToCreate := bufferPolicy
	for _, pt_slice := range tt.Ranges {
		if pt_slice.Status == WAITING {
			slicesToCreate--
		} else {
			fmt.Println("SLICE: ", pt_slice.Name)
			fmt.Println("STATUS: ", pt_slice.Status)
			// Nothng to do here ,.
		}
	}

	if slicesToCreate > 0 {
		// Appedn the remainig from now ..
	}
	// Mark WAITING until Buffer filled
	return nil
}

// determineNextCheckPoint finds next micro boundary ; if passed do now
func determineNextCheckPoint(currentTime time.Time) time.Time {
	// if currentTime > Now; fire it off immediately ..
	// Dummy vakue .. if passed; do now
	return time.Now()
}
