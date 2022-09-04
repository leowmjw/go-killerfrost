package breakglass

import (
	"context"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jackc/pgx/v5"
)

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
)

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
			//return nil
		}
		// Unexpected ..
		return fmt.Errorf("Error reading grant role: %w", err)
	}
	// DEBUG ..
	spew.Dump(values)
	return nil
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
