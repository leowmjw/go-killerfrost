# RESEARCH

## Postgres w/ Gzolang

- Tern - https://github.com/jackc/tern - migration to setup as app user assume roles
- pgx - https://github.com/jackc/pgx
- puddle - https://github.com/jackc/puddle - PG pooling 
- PG Mock - https://github.com/jackc/pgproto3/tree/master/example/pgfortune
- SQL Builder - https://github.com/jackc/pgsql
- PG Proxy / Sniffer - https://github.com/jackc/pgmock
- Fast Rollback tests - https://github.com/jackc/pgundolog
- util collections - https://github.com/jackc/pgxutil/blob/master/collect_test.go
- util connectivity - https://github.com/jackc/pgxutil/blob/master/pgxutil_test.go
- PGErrCode - https://github.com/jackc/pgerrcode
- ZeroLog Adapter - https://github.com/jackc/pgx-zerolog/blob/master/adapter_test.go
- 
## Teleport

Try out Teleport to Postgres + Snowflake

Look into its plugin system if SQL can be manipulated

Create a Teleport role with read-only access to various schema read users

## Temporal 

Have a Golang WebRPC endpoint representing request-response pattern for
moving a user in and out of a Teleport role

## Slack-wrench

Try out BoltJS with higher level abstraction to have elevated access to pre-defined schema / views;
powered by Temporalite

Dynamically add user into pre-defined Teleport role
