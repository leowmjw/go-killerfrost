package workflows

import (
	"app/internal/partition-iwf/workflows/basic"
	"github.com/iworkflowio/iwf-golang-sdk/iwf"
)

var registry = iwf.NewRegistry()

func init() {
	err := registry.AddWorkflows(
		&basic.PartitionWorkflow{},
		//&basic.BasicWorkflow{},
		//&interstate.InterStateWorkflow{},
		//&persistence.PersistenceWorkflow{},
		//&signal.SignalWorkflow{},
		//&timer.TimerWorkflow{},
	)
	if err != nil {
		panic(err)
	}
}

func GetRegistry() iwf.Registry {
	return registry
}
