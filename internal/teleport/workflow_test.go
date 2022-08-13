package teleport

import (
	"go.temporal.io/sdk/testsuite"
	"testing"
	"time"
)

func TestBastion_AddToRole(t *testing.T) {
	type fields struct {
		Identity string
		URL      string
		Token    string
	}
	type args struct {
		roleName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bastion{
				Identity: tt.fields.Identity,
				URL:      tt.fields.URL,
				Token:    tt.fields.Token,
			}
			if err := b.AddToRole(tt.args.roleName); (err != nil) != tt.wantErr {
				t.Errorf("AddToRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBastion_RemoveFromRole(t *testing.T) {
	type fields struct {
		Identity string
		URL      string
		Token    string
	}
	type args struct {
		roleName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bastion{
				Identity: tt.fields.Identity,
				URL:      tt.fields.URL,
				Token:    tt.fields.Token,
			}
			if err := b.RemoveFromRole(tt.args.roleName); (err != nil) != tt.wantErr {
				t.Errorf("RemoveFromRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBreakGlassWorkflow(t *testing.T) {
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
		{"can dump", args{callbacks: []signal{
			{"breakglass", BreakGlassSignal{
				Action: BG_OPS_DUMP,
			}, time.Millisecond},
		}}, true},
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
	}
	ts := testsuite.WorkflowTestSuite{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := ts.NewTestWorkflowEnvironment()
			// All tests should be done by 1s
			env.SetTestTimeout(time.Second)
			env.RegisterWorkflow(BreakGlassWorkflow)
			// Setuo the signals ..
			for _, signal := range tt.args.callbacks {
				// As per normal footgun :sweat:
				// Make a copy; otherwise the BGSignal content gets overwritten ,..
				signal := signal
				env.RegisterDelayedCallback(func() {
					env.SignalWorkflow(signal.name, signal.content)
				}, signal.when)
			}
			// Original
			//env.RegisterDelayedCallback(func() {
			//	env.SignalWorkflow("breakglass", BreakGlassSignal{
			//		Action: BG_OPS_DUMP,
			//	})
			//}, time.Millisecond)
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
