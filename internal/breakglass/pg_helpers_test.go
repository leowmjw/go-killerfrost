package breakglass

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jackc/pgx/v5"
	"testing"
)

func Test_readGrantRole(t *testing.T) {
	t.Parallel()

	type args struct {
		conn     *pgx.Conn
		userName string
		roleName string
	}
	connString := "postgres://foo:password@127.0.0.1:5432/myterraform"
	conn, err := pgx.Connect(context.Background(), connString)
	//config, err := pgx.ParseConfig(connString)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"happy #1", args{
			conn:     conn,
			userName: "foo",
			roleName: "s2admin",
		}, false},
		{"sad #1", args{
			conn:     conn,
			userName: "foounknown",
			roleName: "s2admin",
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := readGrantRole(tt.args.conn, tt.args.userName, tt.args.roleName); (err != nil) != tt.wantErr {
				t.Errorf("readGrantRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_revokeRoleMembership(t *testing.T) {
	t.Parallel()
	type args struct {
		setup    func(*pgx.Conn)
		conn     *pgx.Conn
		userName string
		roleName string
	}
	connString := "postgres://foo:password@127.0.0.1:5432/myterraform"
	conn, err := pgx.Connect(context.Background(), connString)
	//config, err := pgx.ParseConfig(connString)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"happy #1", args{
			userName: "s2admin",
			roleName: "s2admin",
		}, false, false},
		{"happy #2", args{
			conn:     conn,
			userName: "s2admin",
			roleName: "foo",
		}, false, false},
		{"happy #3", args{
			setup: func(conn *pgx.Conn) {
				ct, xerr := conn.Exec(context.Background(),
					fmt.Sprintf("GRANT %s TO %s", "s2read", "backend"),
				)
				if xerr != nil {
					spew.Dump(xerr)
					return
				}
				fmt.Println("EXEC_RESULT: ", ct.String(), " NO: ", ct.RowsAffected())
			},
			conn:     conn,
			userName: "backend",
			roleName: "s2read",
		}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup to add membership ..
			if tt.args.setup != nil {
				// If has func; run it with conn!!
				tt.args.setup(tt.args.conn)
			}
			got, err := revokeRoleMembership(tt.args.conn, tt.args.userName, tt.args.roleName)
			if (err != nil) != tt.wantErr {
				t.Errorf("revokeRoleMembership() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("revokeRoleMembership() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_grantRoleMembership(t *testing.T) {
	t.Parallel()
	type args struct {
		conn     *pgx.Conn
		userName string
		roleName string
	}
	connString := "postgres://foo:password@127.0.0.1:5432/myterraform"
	conn, err := pgx.Connect(context.Background(), connString)
	//config, err := pgx.ParseConfig(connString)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"happy #0 - ignore role = user", args{
			conn:     conn,
			userName: "foo",
			roleName: "foo",
		}, false, false},
		{"happy #1 - nothing to do", args{
			conn:     conn,
			userName: "foo",
			roleName: "s2admin",
		}, false, false},
		{"happy #2", args{
			conn:     conn,
			userName: "backend",
			roleName: "s2read",
		}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// With TXN no need to worry about mutations ..
			tx, terr := conn.Begin(context.Background())
			if terr != nil {
				t.Error(terr)
				return
			}
			defer tx.Rollback(context.Background())
			got, err := grantRoleMembership(tt.args.conn, tt.args.userName, tt.args.roleName)
			if (err != nil) != tt.wantErr {
				t.Errorf("grantRoleMembership() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("grantRoleMembership() got = %v, want %v", got, tt.want)
			}
		})
	}
}
