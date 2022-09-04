package breakglass

import (
	"context"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jackc/pgx/v5"
)

type RoleNotFound error

// Follow like: https://github.com/cyrilgdn/terraform-provider-postgresql
const (
	// This returns the role membership for role, grant_role
	getGrantRoleQuery = `
SELECT
  pg_get_userbyid(member) as role,
  pg_get_userbyid(roleid) as grant_role
FROM
  pg_auth_members
WHERE
  pg_get_userbyid(member) = $1 AND
  pg_get_userbyid(roleid) = $2;
`
	// This returns the role membership for role, grant_role
	//	getRoleExistQuery = `
	//SELECT 1
	//FROM
	//  pg_auth_members
	//WHERE
	//  pg_get_userbyid(member) = $1 AND
	//  pg_get_userbyid(roleid) = $2;
	//`
)

// TODO: If need to check if is a member/user or ROLE
// Follow SQL convention; ROLE == NOLOGIN while user/member == LOGIN
// ROLE --> SELECT oid, rolname, rolcanlogin  FROM pg_roles WHERE rolcanlogin = 'f';
// USER/MEMBER --> SELECT oid, rolname, rolcanlogin  FROM pg_roles WHERE rolcanlogin = 't';

// func readGrantRole(db QueryAble, d *schema.ResourceData) error {
// readGrantRole to confirm GRANT is as expected
func readGrantRole(conn *pgx.Conn, userName, roleName string) error {
	var checkRoleName, checkGrantRoleName string
	values := []any{
		&checkRoleName,
		&checkGrantRoleName,
	}
	row := conn.QueryRow(context.Background(), getGrantRoleQuery, userName, roleName)
	err := row.Scan(values...)
	if err != nil {
		// If no row; WARN ..
		if errors.Is(err, pgx.ErrNoRows) {
			fmt.Println("[WARN] PostgreSQL grant role ", roleName, " to ", userName, " NOT found!")
			return errors.New("NOT FOUND!")
			//return nil
		}
		// Unexpected ..
		return fmt.Errorf("Error reading grant role: %w", err)
	}
	// DEBUG ..
	spew.Dump(values)
	return nil
}

// revokeRoleMembership revokes the role *role* from the user *member*.
// It returns false if the revoke is not needed because the user is not a member of this role.
// func revokeRoleMembership(db QueryAble, role, member string) (bool, error) {
// revokeRoleMembership
func revokeRoleMembership(conn *pgx.Conn, userName, roleName string) (bool, error) {
	// Cannot revoke user from itself!
	if userName == roleName {
		return false, nil
	}
	// Check if already a member; if not mark it and exit
	var checkRoleName, checkGrantRoleName string
	values := []any{
		&checkRoleName,
		&checkGrantRoleName,
	}
	row := conn.QueryRow(context.Background(), getGrantRoleQuery, userName, roleName)
	err := row.Scan(values...)
	if err != nil {
		// If no row; WARN ..
		if errors.Is(err, pgx.ErrNoRows) {
			fmt.Println("[WARN] PostgreSQL grant role ", roleName, " to ", userName, " NOT found!")
			return false, nil
		}
		// Unexpected ..
		return false, fmt.Errorf("Error reading grant role: %w", err)
	}
	// now actually revoke!
	ct, xerr := conn.Exec(context.Background(), fmt.Sprintf("REVOKE %s FROM %s", roleName, userName))
	// Below does not work .. :sweat: ..
	//ct, xerr := conn.Exec(context.Background(), "REVOKE '$1' FROM '$2'", roleName, userName)
	if xerr != nil {
		spew.Dump(xerr)
		return false, xerr
	}
	fmt.Println("EXEC_RESULT: ", ct.String(), " NO: ", ct.RowsAffected())

	return true, nil
}

//func createGrantRoleQuery(d *schema.ResourceData) string {
//	grantRole, _ := d.Get("grant_role").(string)
//	role, _ := d.Get("role").(string)
//
//	query := fmt.Sprintf(
//		"GRANT %s TO %s",
//		pq.QuoteIdentifier(grantRole),
//		pq.QuoteIdentifier(role),
//	)
//	if wao, _ := d.Get("with_admin_option").(bool); wao {
//		query = query + " WITH ADMIN OPTION"
//	}
//
//	return query
//}
//
//func createRevokeRoleQuery(d *schema.ResourceData) string {
//	grantRole, _ := d.Get("grant_role").(string)
//	role, _ := d.Get("role").(string)
//
//	return fmt.Sprintf(
//		"REVOKE %s FROM %s",
//		pq.QuoteIdentifier(grantRole),
//		pq.QuoteIdentifier(role),
//	)
//}

//func grantRole(txn *sql.Tx, d *schema.ResourceData) error {
//	query := createGrantRoleQuery(d)
//	if _, err := txn.Exec(query); err != nil {
//		return fmt.Errorf("could not execute grant query: %w", err)
//	}
//	return nil
//}
//
//func revokeRole(txn *sql.Tx, d *schema.ResourceData) error {
//	query := createRevokeRoleQuery(d)
//	if _, err := txn.Exec(query); err != nil {
//		return fmt.Errorf("could not execute revoke query: %w", err)
//	}
//	return nil
//}
