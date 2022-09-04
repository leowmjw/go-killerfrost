package breakglass

import (
	"context"
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
