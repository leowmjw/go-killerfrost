package basic

import (
	"encoding/json"
	"github.com/iworkflowio/iwf-golang-sdk/iwf"
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

type PartitionWorkflow struct{}

func (b PartitionWorkflow) GetStates() []iwf.StateDef {
	return []iwf.StateDef{
		//iwf.NewStartingState(&stead{}),
		//iwf.NewNonStartingState(&basicWorkflowState2{}),
	}
}

func (b PartitionWorkflow) GetPersistenceSchema() []iwf.PersistenceFieldDef {
	return nil
}

func (b PartitionWorkflow) GetCommunicationSchema() []iwf.CommunicationMethodDef {
	return []iwf.CommunicationMethodDef{
		iwf.NewSignalChannelDef(SignalName),
	}
}

func (b PartitionWorkflow) GetWorkflowType() string {
	return ""
}
