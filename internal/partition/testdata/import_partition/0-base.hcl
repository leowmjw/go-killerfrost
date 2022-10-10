// CReated by running below ..
// atlas schema inspect -u "postgres://foo:password@localhost:5432/myterraform?sslmode=disable"
table "measurement" {
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
  index "measurement_logdate_idx" {
    columns = [column.logdate]
    type    = BTREE
  }
  partition {
    type    = RANGE
    columns = [column.logdate]
  }
}

schema "s2" {
}