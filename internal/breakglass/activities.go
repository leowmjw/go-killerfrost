package breakglass

import (
	"context"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Bastion struct {
	Identity string
	URL      string
	Token    string
	PGConn   *pgx.Conn
}

func NewBastion(connString string) (error, Bastion) {
	//pgconn.ParseConfig()
	cc, err := pgx.ParseConfig(connString)
	if err != nil {
		return err, Bastion{}
	}
	conn, cerr := pgx.ConnectConfig(context.Background(), cc)
	if cerr != nil {
		return cerr, Bastion{}
	}
	//pgx.ConnectConfig()
	return nil, Bastion{
		Identity: "",
		URL:      "",
		Token:    "",
		PGConn:   conn,
	}
}
func (b Bastion) AddToRole(userName, roleName string) error {
	// GRANT s2write TO backend ;
	// Make the connection short term ..
	cc, err := pgx.ParseConfig(b.URL)
	if err != nil {
		return err
	}
	conn, cerr := pgx.ConnectConfig(context.Background(), cc)
	if cerr != nil {
		return cerr
	}
	defer conn.Close(context.Background())
	// With TXN no need to worry about inconsistency
	tx, terr := conn.Begin(context.Background())
	if terr != nil {
		return terr
	}
	// Do actual granting ...
	granted, gerr := grantRoleMembership(conn, userName, roleName)
	if gerr != nil {
		// Bail the TXN! when unexpected
		rerr := tx.Rollback(context.Background())
		if rerr != nil {
			fmt.Println("[WARN] Rollback ERR: ", rerr.Error())
		}
		return gerr
	}
	// Just flag unexpected behavior .
	if !granted {
		fmt.Println("[WARN] Unexpected behavior but no harm, continue on ..")
	}
	// All good.. commit!
	cmterr := tx.Commit(context.Background())
	if cmterr != nil {
		fmt.Println("[WARN] Commit ERR: ", cerr.Error())
	}
	return nil
}

func (b Bastion) RemoveFromRole(userName, roleName string) error {
	// REVOKE s2write FROM backend ;
	// Make the connection short term ..
	cc, err := pgx.ParseConfig(b.URL)
	if err != nil {
		return err
	}
	conn, cerr := pgx.ConnectConfig(context.Background(), cc)
	if cerr != nil {
		return cerr
	}
	defer conn.Close(context.Background())
	revoked, rerr := revokeRoleMembership(conn, userName, roleName)
	if rerr != nil {
		return rerr
	}
	// Just flag unexpected behavior .
	if !revoked {
		fmt.Println("[WARN] Unexpected behavior but no harm, continue on ..")
	}
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
	_, err := b.PGConn.Exec(context.Background(), fmt.Sprintf("SET ROLE %s", roleName))
	if err != nil {
		// If want to be more specific?
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			fmt.Println(err.Error())
			fmt.Println("DEETS:", pgerr.Code)
		} else {
			spew.Dump(err)
		}
		return false, err
	}
	//fmt.Println("ROWS: ", ct.RowsAffected())
	// Log out the current user .. diff from session!
	cs, xerr := b.PGConn.Query(context.Background(), "SELECT current_user, session_user")
	if xerr != nil {
		spew.Dump(xerr)
		return false, xerr
	}
	//fmt.Println("USER: ", cs.String())
	if cs.Next() {
		for _, col := range cs.RawValues() {
			spew.Dump(string(col))
		}
	}
	if cs.Next() {
		fmt.Println("ERR! More rows ..")
	}
	// If exec after a defer; conn will be busy ..
	//defer cs.Close()
	cs.Close()
	// If role exist and have permission by session ..

	// If role exist but NO permission by session

	// If does not exist
	return false, nil
}
