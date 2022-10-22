package partition

import "time"

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
