# Setup

## Plan

``
$ atlas schema apply --dry-run -s s2 -u "postgres://foo:password@localhost:5432/myterraform?sslmode=disable" -f ./import_partition
-- Planned Changes:
-- Modify "measurement" table
ALTER TABLE "s2"."measurement" ADD COLUMN "description" text NULL
``
## Apply 

