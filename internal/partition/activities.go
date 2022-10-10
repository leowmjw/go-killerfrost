package partition

func CreateNewPartitionSlice(policy Policy) error {
	// Copy from existing one

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

func AttachNewDateRange(table TrackedTable) error {
	/*
			\copy measurement_y2008m02 from 'measurement_y2008m02'
		-- possibly some other data preparation work

		ALTER TABLE measurement ATTACH PARTITION measurement_y2008m02
		    FOR VALUES FROM ('2008-02-01') TO ('2008-03-01' );

	*/

	return nil
}

func DetachDateRange(table TrackedTable) error {
	return nil
}

func relationExists(tableName string) (bool, error) {

	// If reach here; it exists?
	return true, nil
}

func ArchiveDateRange(dateRange DateRange) error {
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
