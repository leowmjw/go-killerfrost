# SETUP

For Demo; setup dummy to local Postgres; recommend Postgres.app

Considered the following as option to properly setup DB:

- https://github.com/rallyware/terraform-postgresql-bootstrap
- https://github.com/jparnaudeau/terraform-postgresql-database-admin <-- looks good!

Implements auto the best-practices as per documented in https://aws.amazon.com/blogs/database/managing-postgresql-users-and-roles/

Refer to the commands for use later to enable + revoke permissions for later demo:

- Read from existing table
- Write to existing table
- Create new table in the DB/scehema
- Temp app user with rotated credentials to allow safe switch of slot
- More sophisticated removal of inactive users; suspend, lock --> https://aws.amazon.com/blogs/database/managing-inactive-amazon-aurora-postgresql-users/
