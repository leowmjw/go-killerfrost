package partition

// Functions utilizing Atlas to do the Plan, Diff + Execute ..
// Try out in CLI first; before try to reverse engineer it ..

// Before any action; should use the schema app role ..
/*
myterraform=> RESET ROLE;
RESET
myterraform=# SELECT SESSION_USER, CURRENT_USER;
 session_user | current_user
--------------+--------------
 leow         | leow
(1 row)

myterraform=# SET ROLE 's2admin';
SET
myterraform=> SELECT SESSION_USER, CURRENT_USER;
 session_user | current_user
--------------+--------------
 leow         | s2admin
(1 row)
*/

// Create virtual table as partiiton
/*
myterraform=> DROP TABLE measurement;
DROP TABLE
myterraform=> CREATE TABLE measurement (
myterraform(>     city_id         int not null,
myterraform(>     logdate         date not null,
myterraform(>     peaktemp        int,
myterraform(>     unitsales       int
myterraform(> ) PARTITION BY RANGE (logdate);
CREATE TABLE
myterraform=> \dt
                 List of relations
 Schema |    Name     |       Type        |  Owner
--------+-------------+-------------------+---------
 s2     | measurement | partitioned table | s2admin
(1 row)
myterraform=> \dp+
                                         Access privileges
 Schema |    Name     |       Type        |    Access privileges    | Column privileges | Policies
--------+-------------+-------------------+-------------------------+-------------------+----------
 s2     | measurement | partitioned table | s2admin=arwdDxt/s2admin+|                   |
        |             |                   | s2read=rxt/s2admin     +|                   |
        |             |                   | s2write=arwdxt/s2admin  |                   |
(1 row)

myterraform=> \dP+
                             List of partitioned relations
 Schema |    Name     |  Owner  |       Type        | Table | Total size | Description
--------+-------------+---------+-------------------+-------+------------+-------------
 s2     | measurement | s2admin | partitioned table |       | 0 bytes    |
(1 row)
*/

// Add more partitions + indexes ..
/*
myterraform=> CREATE TABLE measurement_y2006m02 PARTITION OF measurement
myterraform->     FOR VALUES FROM ('2006-02-01') TO ('2006-03-01');
CREATE TABLE
myterraform=>
myterraform=> CREATE TABLE measurement_y2006m03 PARTITION OF measurement
myterraform->     FOR VALUES FROM ('2006-03-01') TO ('2006-04-01');
CREATE TABLE
myterraform=> DROP TABLE measurement_y2006m02 ;
DROP TABLE
myterraform=> CREATE TABLE measurement_y2006m02 PARTITION OF measurement                                                           FOR VALUES FROM ('2006-02-01') TO ('2006-03-01')                                                                               PARTITION BY RANGE (peaktemp);
CREATE TABLE
 CREATE INDEX ON measurement (logdate);
CREATE INDEX
                                                                    List of relations
 Schema |               Name               |       Type        |  Owner  |        Table         | Persistence | Access method |    Size    | Description
--------+----------------------------------+-------------------+---------+----------------------+-------------+---------------+------------+-------------
 s2     | measurement_logdate_idx          | partitioned index | s2admin | measurement          | permanent   | btree         | 0 bytes    |
 s2     | measurement_y2006m02_logdate_idx | partitioned index | s2admin | measurement_y2006m02 | permanent   | btree         | 0 bytes    |
 s2     | measurement_y2006m03_logdate_idx | index             | s2admin | measurement_y2006m03 | permanent   | btree         | 8192 bytes |
(3 rows)

(END)
                                                                    List of relations
 Schema |               Name               |       Type        |  Owner  |        Table         | Persistence | Access method |    Size    | Description
--------+----------------------------------+-------------------+---------+----------------------+-------------+---------------+------------+-------------
 s2     | measurement_logdate_idx          | partitioned index | s2admin | measurement          | permanent   | btree         | 0 bytes    |
 s2     | measurement_y2006m02_logdate_idx | partitioned index | s2admin | measurement_y2006m02 | permanent   | btree         | 0 bytes    |
 s2     | measurement_y2006m03_logdate_idx | index             | s2admin | measurement_y2006m03 | permanent   | btree         | 8192 bytes |
(3 rows)
*/

// Detach for use by month access; by regulation up till 3 years; concurrently to not block
// Can reattach again ..
/*
myterraform=> ALTER TABLE measurement DETACH PARTITION measurement_y2006m03 CONCURRENTLY ;
ALTER TABLE
myterraform=> \dp+
                                             Access privileges
 Schema |         Name         |       Type        |    Access privileges    | Column privileges | Policies
--------+----------------------+-------------------+-------------------------+-------------------+----------
 s2     | measurement          | partitioned table | s2admin=arwdDxt/s2admin+|                   |
        |                      |                   | s2read=rxt/s2admin     +|                   |
        |                      |                   | s2write=arwdxt/s2admin  |                   |
 s2     | measurement_y2006m02 | partitioned table | s2admin=arwdDxt/s2admin+|                   |
        |                      |                   | s2read=rxt/s2admin     +|                   |
        |                      |                   | s2write=arwdxt/s2admin  |                   |
 s2     | measurement_y2006m03 | table             | s2admin=arwdDxt/s2admin+|                   |
        |                      |                   | s2read=rxt/s2admin     +|                   |
        |                      |                   | s2write=arwdxt/s2admin  |                   |
(3 rows)

myterraform=> \dP+
                                      List of partitioned relations
 Schema |          Name           |  Owner  |       Type        |    Table    | Total size | Description
--------+-------------------------+---------+-------------------+-------------+------------+-------------
 s2     | measurement             | s2admin | partitioned table |             | 0 bytes    |
 s2     | measurement_logdate_idx | s2admin | partitioned index | measurement | 0 bytes    |
(2 rows)

myterraform=> \di+
                                                                    List of relations
 Schema |               Name               |       Type        |  Owner  |        Table         | Persistence | Access method |    Size    | Description
--------+----------------------------------+-------------------+---------+----------------------+-------------+---------------+------------+-------------
 s2     | measurement_logdate_idx          | partitioned index | s2admin | measurement          | permanent   | btree         | 0 bytes    |
 s2     | measurement_y2006m02_logdate_idx | partitioned index | s2admin | measurement_y2006m02 | permanent   | btree         | 0 bytes    |
 s2     | measurement_y2006m03_logdate_idx | index             | s2admin | measurement_y2006m03 | permanent   | btree         | 8192 bytes |
(3 rows)

myterraform=> ALTER TABLE measurement ATTACH  PARTITION measurement_y2006m03 FOR VALUES FROM ('2006-03-01') TO ('2006-04-01');
ALTER TABLE
myterraform=> \di+
                                                                    List of relations
 Schema |               Name               |       Type        |  Owner  |        Table         | Persistence | Access method |    Size    | Description
--------+----------------------------------+-------------------+---------+----------------------+-------------+---------------+------------+-------------
 s2     | measurement_logdate_idx          | partitioned index | s2admin | measurement          | permanent   | btree         | 0 bytes    |
 s2     | measurement_y2006m02_logdate_idx | partitioned index | s2admin | measurement_y2006m02 | permanent   | btree         | 0 bytes    |
 s2     | measurement_y2006m03_logdate_idx | index             | s2admin | measurement_y2006m03 | permanent   | btree         | 8192 bytes |
(3 rows)

myterraform=> \dP+
                                      List of partitioned relations
 Schema |          Name           |  Owner  |       Type        |    Table    | Total size | Description
--------+-------------------------+---------+-------------------+-------------+------------+-------------
 s2     | measurement             | s2admin | partitioned table |             | 0 bytes    |
 s2     | measurement_logdate_idx | s2admin | partitioned index | measurement | 8192 bytes |
(2 rows)

myterraform=> \dp+
                                             Access privileges
 Schema |         Name         |       Type        |    Access privileges    | Column privileges | Policies
--------+----------------------+-------------------+-------------------------+-------------------+----------
 s2     | measurement          | partitioned table | s2admin=arwdDxt/s2admin+|                   |
        |                      |                   | s2read=rxt/s2admin     +|                   |
        |                      |                   | s2write=arwdxt/s2admin  |                   |
 s2     | measurement_y2006m02 | partitioned table | s2admin=arwdDxt/s2admin+|                   |
        |                      |                   | s2read=rxt/s2admin     +|                   |
        |                      |                   | s2write=arwdxt/s2admin  |                   |
 s2     | measurement_y2006m03 | table             | s2admin=arwdDxt/s2admin+|                   |
        |                      |                   | s2read=rxt/s2admin     +|                   |
        |                      |                   | s2write=arwdxt/s2admin  |                   |
(3 rows)

*/

// Last 6 months snpashot == last 5 months + current_month view?

// More Refs:
// - https://www.enterprisedb.com/postgres-tutorials/how-use-table-partitioning-scale-postgresql
// - https://evilmartians.com/chronicles/a-slice-of-life-table-partitioning-in-postgresql-databases
// - https://hevodata.com/learn/postgresql-partitions/
