package basic

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/iworkflowio/iwf-golang-sdk/gen/iwfidl"
	"github.com/iworkflowio/iwf-golang-sdk/iwf"
	"github.com/jackc/pgx/v5"
	"time"
)

// Diff State is when partitions are not matching ideal state
//	goes to Steady once confirmed all matching plan is done ..

// basic skeleton below; replace Diff with your own ..

const (
	DiffID = "Diff"
)

type DiffState struct {
	db PostgresDB
}

func (b DiffState) GetStateId() string {
	return DiffID
}

func (b DiffState) Start(ctx iwf.WorkflowContext, input iwf.Object, persistence iwf.Persistence, communication iwf.Communication) (*iwf.CommandRequest, error) {
	fmt.Println("DIFF_START ******")
	spew.Dump(b.db)
	if b.db.Conn == nil {
		fmt.Println("ONe time DB in Diff ..")
		// connString can be passed along probably; from setup?
		connString := "postgres://s2admin:password@127.0.0.1:5432/myterraform"
		connConfig, err := pgx.ParseConfig(connString)
		if err != nil {
			fmt.Println("ERR:", err)
		} else {
			// DEBUG
			//fmt.Println("CTX: yyyy")
			//spew.Dump(ctx)
			//fmt.Println("CONN: xxxxxxx")
			//spew.Dump(connConfig.ConnectTimeout)
			// All checks out ..
			c, cerr := pgx.ConnectConfig(context.Background(), connConfig)
			if cerr != nil {
				fmt.Println("CERR:", cerr)
				panic(cerr)
			}
			b.db = PostgresDB{
				ConnConfig: connConfig,
				Conn:       c,
			}
			//fmt.Println("AFTER SETUP ..")
			//spew.Dump(b.db.ConnConfig)
		}
	}

	fmt.Println("BEFORE EXIT START ..")
	spew.Dump(b.db.ConnConfig)
	// Attach dynamic timer ..
	//return iwf.EmptyCommandRequest(), nil

	return iwf.AnyCommandCompletedRequest(
		iwf.NewTimerCommand("check", time.Now().Add(10*time.Second)),
	), nil
}

func (b DiffState) Decide(ctx iwf.WorkflowContext, input iwf.Object, commandResults iwf.CommandResults, persistence iwf.Persistence, communication iwf.Communication) (*iwf.StateDecision, error) {
	fmt.Println("DIFF_DECIDE ******")
	spew.Dump(b.db)
	//var i int
	// Loop to itself ..
	var tt TrackedTables
	gerr := persistence.GetDataObject("TrackedTables", &tt)
	if gerr != nil {
		fmt.Println("ERR??")
		//panic(gerr)
	}
	// Will be zero-valued ..
	if tt.ID == "" {
		fmt.Println("Should NOTz be here!! DECIDE")
	} else {
		fmt.Println("FROM_STEADY:", tt.ID)
		spew.Dump(tt.Planned)
	}

	return iwf.SingleNextState(DiffID, nil), nil
}

func (b DiffState) GetStateOptions() *iwfidl.WorkflowStateOptions {
	fmt.Println("Inside GetStateOptions ==>  ")
	v, err := iwfidl.NewWorkflowStateOptionsWithDefaults().MarshalJSON()
	if err != nil {
		panic(err)
	}
	spew.Dump(v)
	//iwfidl.NewNullableWorkflowStateOptions(nil)
	return nil
}
