package partition

// Postgres 14 can detach partition without interruption
// When inactive; detached from main partition
// After retention time; mark for removal with humam last time review ..

// Supported Partition
// Test: Day -> Year; Hours -> Month; Min -> Day
// Production: Year; Month; Day
// Start with fixed; next time can be flexible; e.g. workday, weekend, holidays ..

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
