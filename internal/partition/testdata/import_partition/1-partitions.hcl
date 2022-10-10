// As per in: https://www.postgresql.org/docs/current/ddl-partitioning.html#DDL-PARTITIONING-DECLARATIVE
// CREATE TABLE measurement_y2008m02
//  (LIKE measurement INCLUDING DEFAULTS INCLUDING CONSTRAINTS)

// ALTER TABLE measurement_y2008m02 ADD CONSTRAINT y2008m02
//   CHECK ( logdate >= DATE '2008-02-01' AND logdate < DATE '2008-03-01' );

// ALTER TABLE measurement ATTACH PARTITION measurement_y2008m02
//    FOR VALUES FROM ('2008-02-01') TO ('2008-03-01' );

// Apply as per below:
// $ atlas schema apply -s s2 -u "postgres://foo:password@localhost:5432/myterraform?sslmode=disable" -f ./import_partition
// Schema is synced, no changes to be made
// Table below created from the parent level
table "measurement_y2008m02" {
  schema = schema.s2
  column "city_id" {
    null = false
    type = integer
  }
  column "logdate" {
    null = false
    type = date
  }
  column "peaktemp" {
    null = true
    type = integer
  }
  column "unitsales" {
    null = true
    type = integer
  }
  column "description" {
    null = true
    type = text
  }
  check "y2008m02" {
    expr = "((logdate >= '2008-02-01'::date) AND (logdate < '2008-03-01'::date))"
  }
}
