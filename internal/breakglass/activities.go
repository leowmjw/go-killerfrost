package breakglass

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jackc/pgx/v5"
)

type Bastion struct {
	Identity string
	URL      string
	Token    string
	PGConn   *pgx.Conn
}

func (b Bastion) AddToRole(roleName string) error {
	// GRANT s2write TO backend ;
	return nil
}

func (b Bastion) RemoveFromRole(roleName string) error {
	// REVOKE s2write FROM backend ;

	return nil
}

func (b Bastion) GetRole() (string, error) {
	var roleName string

	// Query Teleport role ...
	if roleName == "" {
		return "", fmt.Errorf("INVLID ROLE ..!! %v", b)
	}
	return roleName, nil
}

func (b Bastion) RoleExists(roleName string) (bool, error) {
	fmt.Println("In RoleExists ..")
	fmt.Println("CONF: ", b.PGConn.Config().ConnString())
	// SET ROLE s2write
	ct, err := b.PGConn.Exec(context.Background(), "SET ROLE s2admin")
	if err != nil {
		spew.Dump(err)
		return false, err
	}
	fmt.Println("ROWS: ", ct.RowsAffected())
	// Log out the current session ..
	// If role exist and have permission by session ..

	// If role exist but NO permission by session

	// If does not exist
	return false, nil
}
