package partition

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jackc/pgx/v5"
)

type PostgresDB struct {
	ConnConfig *pgx.ConnConfig
}

// In Workflow for table WFID
// 	for those WAITING
//	Go Concurrent
//		if non-fatal-err: if already running; is unique; will not do anything; just log it

func (pgd PostgresDB) CreateNewPartitionSlice(policy Policy) error {
	conn, cerr := pgx.ConnectConfig(context.Background(), pgd.ConnConfig)
	if cerr != nil {
		// logs stuff ..
		return cerr
	}
	// Idempotent; check if table exist? capture this error?
	// Copy from existing parent which is active ..
	copyParentTable := ""
	pct, xerr := conn.Exec(context.Background(), copyParentTable)
	if xerr != nil {
		// logs stuff
		return xerr
	}

	pct.RowsAffected() // how many?
	// is copy wht type of action??

	// Set the constraint
	/*
		SET ROLE to <SCHEMA>admin <-- default pattern
		CREATE TABLE measurement_y2008m02
		  (LIKE measurement INCLUDING DEFAULTS INCLUDING CONSTRAINTS)
		  TABLESPACE fasttablespace;

		ALTER TABLE measurement_y2008m02 ADD CONSTRAINT y2008m02
		   CHECK ( logdate >= DATE '2008-02-01' AND logdate < DATE '2008-03-01' );

	*/
	return nil
}

func (pgd PostgresDB) AttachNewDateRange(table TrackedTable) error {
	/*
			\copy measurement_y2008m02 from 'measurement_y2008m02'
		-- possibly some other data preparation work

		ALTER TABLE measurement ATTACH PARTITION measurement_y2008m02
		    FOR VALUES FROM ('2008-02-01') TO ('2008-03-01' );

	*/

	return nil
}

func (pgd PostgresDB) DetachDateRange(table TrackedTable) error {
	return nil
}

func (pgd PostgresDB) ArchiveDateRange(dateRange DateRange) error {
	fmt.Println("IN: ArchiveDateRange ===================>")
	fmt.Println("<====== ***************************** ==========>")
	spew.Dump(pgd.ConnConfig.ConnString())
	// Confirm the DateRange has been detached .. how?
	// https://stackoverflow.com/questions/71078738/how-to-test-if-a-postgres-partition-has-been-populated-or-not
	// select exists (select * from partition_name limit 1)
	// HOWTO check existing partitions with data?
	/*
		myterraform=# SELECT EXISTS  (select * from measurement_unknown limit 1);
		ERROR:  relation "measurement_unknown" does not exist
		LINE 1: SELECT EXISTS  (select * from measurement_unknown limit 1);
		                                      ^
		myterraform=# SELECT EXISTS  (select * from measurement_y2008m01 limit 1);
		ERROR:  relation "measurement_y2008m01" does not exist
		LINE 1: SELECT EXISTS  (select * from measurement_y2008m01 limit 1);
		                                      ^
		myterraform=# SELECT EXISTS  (select * from measurement_y2008m02 limit 1);
		 exists
		--------
		 t

	*/
	// Bundle into icebreg + S3?
	// Delete the table via DROP ..
	return nil
}

// Helper functions ..
func relationExists(tableName string) (bool, error) {

	// If reach here; it exists?
	return true, nil
}

// Create the dummy test table in the known schema?
// Should it be Setup and TearDown?
// if locked in battle; use the semantics of ubique ..
