package partition

import (
	"fmt"
	"go.temporal.io/sdk/workflow"
	"time"
)

type PartitionNotFound error

// Postgres 14 can detach partition without interruption
// When inactive; detached from main partition
// After retention time; mark for removal with humam last time review ..

type PolUnit int

const (
	POL_UNKNOWN PolUnit = iota
	POL_DAY             = 60 * 20
	POL_MONTH           = POL_DAY * 30
	POL_YEAR            = POL_MONTH * 12
)

type Policy struct {
	UnitInMin int
}

type Status int

const (
	UNKNOWN Status = iota
	WAITING
	ATTACHED
	DETACHING
	ARCHIVED
)

// DateRange is partition segment; to set constraints + range
type DateRange struct {
	Name    string
	MinDate string
	MaxDate string
	Status
}

// TrackedTable represents the lifecycle tracked
type TrackedTable struct {
	Schema string
	Name   string
	Ranges []DateRange
	IsTest bool
}

// Supported Partition
// Test: Day -> Year; Hours -> Month; Min -> Day
// Production: Year; Month; Day
// Start with fixed; next time can be flexible; e.g. workday, weekend, holidays ..
func determineMissingRanges(currentTime time.Time) error {
	// With the currentTime
	//	is it within or more than the SmallestUnit Boundary
	//	If Y:
	//		Create X
	//	If N: Do Nothing
	return nil
}

// Safe: Planner; Execute after Review
// For each minor unit == partition; calculate + plan the missing items

// Plans can be triggered manually or using the new scheduling ..

// Supervisor Workflow tracks the lifetime + lifecycle + audit actions
// have the standard load/unload .. be as backwards compatibel as possible ..
// WorkflowID: <schema>.<table>

// Docs: https://www.postgresql.org/docs/current/ddl-partitioning.html
// Partition Types:
//	- Range Partitioning
//	- List Partitioning
// 	- Hash Partitioning

// Activities
//	- Initialize table if not exist already ..
//	- parititon

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
	var leftOver []DateRange
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
		leftOver = tt.Ranges[seamIndex:len(tt.Ranges)]
		// DEBUG ...
		//spew.Dump(leftOver)
	}
	// DEBUG
	//fmt.Println("MAJOR: ", currentMajorIndex, " MINOR: ", currentMinorIndex)
	var nr []DateRange
	var labelName string
	lenRanges := len(leftOver)
	fmt.Println("Total items: ", lenRanges)
	// Project numProjection into future ..
	for i := startIndex; i < numProjection; i++ {
		// If those are ATTACHED or WAITING .. do nothing ..
		if i < lenRanges {
			if (tt.Ranges[i].Status == ATTACHED) || (tt.Ranges[i].Status == WAITING) {
				continue
			}
		}
		// Rules for the boundary will be different if it is Test ..
		// If those index does not exist; add!
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

// applyPlan only means an approved signal given; or manual call ..
func applyPlan(ctx workflow.Context, tt *TrackedTable) error {
	// TODO: Sanity check on the DB Conn??
	// Check if the parent table exists; if not create it ..
	ctxact := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10,
	})
	workflow.ExecuteActivity(ctxact, PostgresDB.ArchiveDateRange, DateRange{
		Name: "bobo",
	})

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
