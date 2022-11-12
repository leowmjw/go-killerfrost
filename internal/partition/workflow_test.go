package partition

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"go.temporal.io/sdk/testsuite"
	"testing"
	"time"
)

func TestMockedLifeCycleWorkflow(t *testing.T) {
	t.Parallel()
	// For this, we mock the Activity output .. faster ..
}

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
			// Setup Activity + DBs ..
			// although needed only for full integration; can be mocked out ,,
			// Setup the PGD connectivity
			connString := "postgres://s2admin:password@127.0.0.1:5432/myterraform"
			connConfig, perr := pgx.ParseConfig(connString)
			pgx.ConnectConfig(context.Background(), connConfig)
			if perr != nil {
				panic(perr)
			}
			pgd := PostgresDB{ConnConfig: connConfig}
			env.RegisterActivity(&pgd)
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

func TestTrackedTable_getSliceProjection(t *testing.T) {
	type fields struct {
		Schema string
		Name   string
		Ranges []DateRange
		IsTest bool
	}
	type args struct {
		currentTime    time.Time
		numProjection  int
		expectedRanges []DateRange
	}
	// Common test time .. Leap year Feb 29 23:55 ..
	leapYearFeb := time.Date(2020, time.February, 29, 23, 55, 0, 0, time.UTC)
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"happy isTest #1", fields{
			Schema: "s1",
			Name:   "measurement",
			IsTest: true,
		}, args{
			currentTime:   leapYearFeb,
			numProjection: 3,
			expectedRanges: []DateRange{
				{Name: "y0029m23", MinDate: "2022", Status: WAITING},
				{Name: "y0030m00", MinDate: "2022", Status: WAITING},
				{Name: "y0030m01", MinDate: "2022", Status: WAITING},
			},
		}},
		{"happy existing isTest #2", fields{
			Schema: "s1",
			Name:   "measurement",
			Ranges: []DateRange{
				{Name: "y0029m23", MinDate: "2022", Status: ATTACHED},
				{Name: "y0030m00", MinDate: "2022", Status: WAITING},
			},
			IsTest: true,
		}, args{
			currentTime:   leapYearFeb,
			numProjection: 3,
			expectedRanges: []DateRange{
				{Name: "y0029m23", MinDate: "2022", Status: ATTACHED},
				{Name: "y0030m00", MinDate: "2022", Status: WAITING},
				{Name: "y0030m01", MinDate: "2022", Status: WAITING},
			},
		}},
		{"happy existing isTest #3", fields{
			Schema: "s1",
			Name:   "measurement",
			Ranges: []DateRange{
				{Name: "y0029m22", MinDate: "2022", Status: ATTACHED},
				{Name: "y0029m23", MinDate: "2022", Status: WAITING},
				{Name: "y0030m00", MinDate: "2022", Status: WAITING},
			},
			IsTest: true,
		}, args{
			currentTime:   leapYearFeb,
			numProjection: 3,
			expectedRanges: []DateRange{
				{Name: "y0029m22", MinDate: "2022", Status: ATTACHED},
				{Name: "y0029m23", MinDate: "2022", Status: WAITING},
				{Name: "y0030m00", MinDate: "2022", Status: WAITING},
				{Name: "y0030m01", MinDate: "2022", Status: WAITING},
			},
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tt := &TrackedTable{
				Schema: tc.fields.Schema,
				Name:   tc.fields.Name,
				Ranges: tc.fields.Ranges,
				IsTest: tc.fields.IsTest,
			}
			tt.getSliceProjection(tc.args.currentTime, tc.args.numProjection)
			// Check if the projection is correct ..
			if diff := cmp.Diff(tt.Ranges, tc.args.expectedRanges); diff != "" {
				t.Fatalf("want,got: %s", diff)
			}
		})
	}
}
