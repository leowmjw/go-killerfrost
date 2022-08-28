package breakglass

import (
	"context"
	"github.com/jackc/pgx/v5"
	"testing"

	"github.com/jackc/pgx/v5/pgxtest"
)

func TestRoleExistence(t *testing.T) {
	// Setup testcontainers?
	b := Bastion{
		Identity: "",
		URL:      "",
		Token:    "",
	}
	// Setup PG
	defaultConnTestRunner := pgxtest.DefaultConnTestRunner()
	connString := "postgres://foo:password@127.0.0.1:5432/myterraform"
	conn, err := pgx.Connect(context.Background(), connString)
	b.PGConn = conn
	if err != nil {
		t.Fatalf("ERR: %v", err)
	}
	//defaultConnTestRunner.CreateConfig = func(ctx context.Context, t testing.TB) *pgx.ConnConfig {
	//	// connection string: postgres://foo:password@127.0.0.1:5432/myterraform
	//	//config, err := pgx.ParseConfig(os.Getenv("TEST_DATABASE"))
	//	connString := "postgres://foo:password@127.0.0.1:5432/myterraform"
	//	pgx.Connect(context.Background(), connString)
	//	config, err := pgx.ParseConfig(connString)
	//	require.NoError(t, err)
	//
	//	config.OnNotice = func(_ *pgconn.PgConn, n *pgconn.Notice) {
	//		t.Logf("PostgreSQL %s: %s", n.Severity, n.Message)
	//	}
	//	return config
	//}
	//b.PGConn = defaultConnTestRunner.CreateConfig(context.Background(), t)

	defaultConnTestRunner.RunTest(context.Background(), t, func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		// Assume is a powerful admin of DB role minimum; could be superuser ..
		// If role already exist; and try to create and succeed
		// If role exist but no permission to assume ..
		// ERROR:  permission denied to set role "app_admin_role"
		// If the role does not exist; even if no permission
		// ERROR:  role "s2read_2" does not exist

		b.RoleExists("s2admin")

	})

	t.Fail()
}

func TestHappyPathPostgres(t *testing.T) {
	// Revoke any existing role
	// Add new role
	// Set and confirm it is, SET ROLE <role> ; SET ROLE NONE ; /.
	/* WHEN NOT set sessiosn
	myterraform=> SELECT session_user;
	 session_user
	--------------
	 backend
	(1 row)

	myterraform=> SELECT current_user;
	 current_user
	--------------
	 backend
	(1 row)

	myterraform=> SELECT current_setting('role');
	 current_setting
	-----------------
	 none
	(1 row)
	*/

	/* When have roles
	SET ROLE s2admin
	myterraform=> SELECT current_user;
	 current_user
	--------------
	 s2admin
	(1 row)

	myterraform=> SELECT session_user;
	 session_user
	--------------
	 foo
	(1 row)

	myterraform=> SELECT current_setting('role');
	 current_setting
	-----------------
	 s2admin
	(1 row)
	*/
}

func TestBastion_AddToRole(t *testing.T) {
	type fields struct {
		Identity string
		URL      string
		Token    string
	}
	type args struct {
		roleName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bastion{
				Identity: tt.fields.Identity,
				URL:      tt.fields.URL,
				Token:    tt.fields.Token,
			}
			if err := b.AddToRole(tt.args.roleName); (err != nil) != tt.wantErr {
				t.Errorf("AddToRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBastion_RemoveFromRole(t *testing.T) {
	type fields struct {
		Identity string
		URL      string
		Token    string
	}
	type args struct {
		roleName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bastion{
				Identity: tt.fields.Identity,
				URL:      tt.fields.URL,
				Token:    tt.fields.Token,
			}
			if err := b.RemoveFromRole(tt.args.roleName); (err != nil) != tt.wantErr {
				t.Errorf("RemoveFromRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBastion_GetRole(t *testing.T) {
	type fields struct {
		Identity string
		URL      string
		Token    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bastion{
				Identity: tt.fields.Identity,
				URL:      tt.fields.URL,
				Token:    tt.fields.Token,
			}
			got, err := b.GetRole()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRole() got = %v, want %v", got, tt.want)
			}
		})
	}
}
