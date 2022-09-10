package breakglass

import (
	"github.com/jackc/pgx/v5"
	"testing"
)

func TestRoleExistence(t *testing.T) {
	// Setup testcontainers?
	connString := "postgres://foo:password@127.0.0.1:5432/myterraform"
	err, b := NewBastion(connString)
	if err != nil {
		t.Fatal(err)
	}
	var exist bool
	var xerr error

	exist, xerr = b.RoleExists("app_admin_role")
	// Expect err
	if xerr == nil {
		t.Error("EXPECTED error!")
	}
	// If the role does not exist; even if no permission
	// ERROR:  role "s2read_2" does not exist
	exist, xerr = b.RoleExists("s2read_2")
	// Expect err
	if xerr == nil {
		t.Error("EXPECTED error!")
	}
	if exist {
		t.Error("Expected role s2read_2 NOT exists!")
	}

	exist, xerr = b.RoleExists("s2admin")
	if xerr != nil {
		t.Fail()
	}
	if !exist {
		t.Error("Expected role s2admin exists!")
	}
}

func TestBastion_RemoveFromRole(t *testing.T) {
	t.Parallel()
	type fields struct {
		Identity     string
		URL          string
		Token        string
		PGConnConfig *pgx.ConnConfig
	}
	type args struct {
		userName string
		roleName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"happy #1", fields{
			URL: "postgres://foo:password@127.0.0.1:5432/myterraform",
		}, args{
			userName: "backend",
			roleName: "s2read",
		}, false},
		{"happy #2", fields{
			URL: "postgres://foo:password@127.0.0.1:5432/myterraform",
		}, args{
			userName: "backend",
			roleName: "backend",
		}, false},
		{"happy #3", fields{
			URL: "postgres://foo:password@127.0.0.1:5432/myterraform",
		}, args{
			userName: "s2read",
			roleName: "backend",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nerr, b := NewBastion(tt.fields.URL)
			if nerr != nil {
				t.Fatal(nerr)
			}
			if err := b.RemoveFromRole(tt.args.userName, tt.args.roleName); (err != nil) != tt.wantErr {
				t.Errorf("RemoveFromRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBastion_AddToRole(t *testing.T) {
	t.Parallel()
	type fields struct {
		Identity     string
		URL          string
		Token        string
		PGConnConfig *pgx.ConnConfig
	}
	type args struct {
		userName string
		roleName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"happy #0", fields{
			URL: "postgres://foo:password@127.0.0.1:5432/myterraform",
		}, args{
			userName: "s2read",
			roleName: "s2read",
		}, false},
		{"happy #1", fields{
			URL: "postgres://foo:password@127.0.0.1:5432/myterraform",
		}, args{
			userName: "backend",
			roleName: "s2read",
		}, false},
		{"happy #2", fields{
			URL: "postgres://foo:password@127.0.0.1:5432/myterraform",
		}, args{
			userName: "s2read",
			roleName: "backend",
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nerr, b := NewBastion(tt.fields.URL)
			if nerr != nil {
				t.Fatal(nerr)
			}
			// Full integration test, is it better to have func attached?
			// pre, post? no rollback .. commit all the way!
			if err := b.AddToRole(tt.args.userName, tt.args.roleName); (err != nil) != tt.wantErr {
				t.Errorf("AddToRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
