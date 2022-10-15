package partition

import (
	"go.temporal.io/sdk/testsuite"
	"testing"
	"time"
)

func TestLifeCycleWorkflow(t *testing.T) {
	t.Parallel()

	type signal struct {
		name    string
		content PTSignal
		when    time.Duration
	}
	type args struct {
		callbacks []signal
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool // NOt used currently ..
	}{
		{"test dump state", args{callbacks: []signal{
			{"partition", PTSignal{
				Action: PT_OPS_DUMP,
			}, time.Millisecond},
		}}, true},
		{"test checkerloop", args{callbacks: []signal{
			{"partition", PTSignal{
				Action: PT_NOOPS,
			}, time.Millisecond},
			{"partition", PTSignal{
				Action: PT_OPS_DUMP,
			}, time.Minute * 5},
		}}, true},
	}
	// Only need to set once ..
	ts := testsuite.WorkflowTestSuite{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// env needs to be reset each time so can run in parallel
			env := ts.NewTestWorkflowEnvironment()
			// All tests should be done by 1s
			env.SetTestTimeout(time.Second)
			env.RegisterWorkflow(LifeCycleWorkflow)
			// Setuo the signal simulation ..
			for _, signal := range tt.args.callbacks {
				// As per normal footgun :sweat:
				// Make a copy; otherwise the BGSignal content gets overwritten ,..
				signal := signal
				env.RegisterDelayedCallback(func() {
					env.SignalWorkflow(signal.name, signal.content)
				}, signal.when)
			}
			env.ExecuteWorkflow(LifeCycleWorkflow)
			// Expect WF completed ..
			if !env.IsWorkflowCompleted() {
				t.Fail()
			}
			// Ensure no err!
			if xerr := env.GetWorkflowError(); xerr != nil {
				t.Fatal(xerr)
			}
		})
	}
}
