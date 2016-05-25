/**
 * @file session.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2015
 * @brief queries for session table
 */

package db

import (
	"database/sql"
	"time"
)

// Session row
type Session struct {
	ID        int
	TeamID    int
	Session   string
	Timestamp time.Time
}

func createSessionTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "session" (
		id		SERIAL PRIMARY KEY,
		team_id		INTEGER NOT NULL,
		session		TEXT NOT NULL,
		timestamp	TIMESTAMP with time zone DEFAULT now()
	)`)

	return
}

// AddSession add session and fill id
func AddSession(db *sql.DB, s *Session) (err error) {

	stmt, err := db.Prepare("INSERT INTO session (team_id, session) " +
		"VALUES ($1, $2) RETURNING id")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(s.TeamID, s.Session).Scan(&s.ID)
	if err != nil {
		return
	}

	return
}

// GetSessionTeam get team id for session
func GetSessionTeam(db *sql.DB, session string) (teamID int, err error) {

	stmt, err := db.Prepare("SELECT team_id FROM session WHERE session=$1")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(session).Scan(&teamID)
	if err != nil {
		return
	}

	return
}
