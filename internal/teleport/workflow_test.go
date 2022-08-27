package teleport

import (
	"go.temporal.io/sdk/testsuite"
	"testing"
	"time"
)

func TestBreakGlassWorkflow(t *testing.T) {
	t.Parallel()

	type signal struct {
		name    string
		content BreakGlassSignal
		when    time.Duration
	}
	type args struct {
		callbacks []signal
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test dump state", args{callbacks: []signal{
			{"breakglass", BreakGlassSignal{
				Action: BG_OPS_DUMP,
			}, time.Millisecond},
		}}, false},
		{"happy path", args{callbacks: []signal{
			{"breakglass", BreakGlassSignal{
				Action: BG_REQUEST_ACCESS,
			}, time.Millisecond},
			{"breakglass", BreakGlassSignal{
				Action: BG_REQUEST_ACCESS,
			}, 10 * time.Millisecond},
			{"breakglass", BreakGlassSignal{
				Action: BG_REQUEST_APPROVED,
			}, 50 * time.Millisecond},
			{"breakglass", BreakGlassSignal{
				Action: BG_OPS_DUMP,
			}, 100 * time.Millisecond},
		}}, false},
		{"happy path rejected", args{callbacks: []signal{
			{"breakglass", BreakGlassSignal{
				Action: BG_REQUEST_REJECTED,
			}, time.Millisecond},
			{"breakglass", BreakGlassSignal{
				Action: BG_REQUEST_ACCESS,
			}, 10 * time.Millisecond},
			{"breakglass", BreakGlassSignal{
				Action: BG_REQUEST_REJECTED,
			}, 50 * time.Millisecond},
			{"breakglass", BreakGlassSignal{
				Action: BG_OPS_DUMP,
			}, 100 * time.Millisecond},
		}}, false},
	}
	// Only need to set once ..
	ts := testsuite.WorkflowTestSuite{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// env needs to be reset each time so can run in parallel
			env := ts.NewTestWorkflowEnvironment()
			// All tests should be done by 1s
			env.SetTestTimeout(time.Second)
			env.RegisterWorkflow(BreakGlassWorkflow)
			// Setuo the signal simulation ..
			for _, signal := range tt.args.callbacks {
				// As per normal footgun :sweat:
				// Make a copy; otherwise the BGSignal content gets overwritten ,..
				signal := signal
				env.RegisterDelayedCallback(func() {
					env.SignalWorkflow(signal.name, signal.content)
				}, signal.when)
			}
			env.ExecuteWorkflow(BreakGlassWorkflow)
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
