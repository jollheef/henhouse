/**
 * @file score.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date November, 2015
 * @brief queries for score table
 */

package db

import (
	"database/sql"
	"time"
)

// Score row
type Score struct {
	ID        int
	TeamID    int
	Score     int
	Timestamp time.Time
}

func createScoreTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "score" (
		id		SERIAL PRIMARY KEY,
		team_id		INTEGER NOT NULL,
		score		INTEGER NOT NULL,
		timestamp	TIMESTAMP with time zone DEFAULT now()
	)`)

	return
}

// AddScore add last result for team and fill id
func AddScore(db *sql.DB, score *Score) (err error) {

	stmt, err := db.Prepare("INSERT INTO score (team_id, score) " +
		"VALUES ($1, $2) RETURNING id")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(score.TeamID, score.Score).Scan(&score.ID)
	if err != nil {
		return
	}

	return
}

// GetLastScore get last result for team id
func GetLastScore(db *sql.DB, teamID int) (s Score, err error) {

	stmt, err := db.Prepare("SELECT id, score, timestamp FROM score " +
		" WHERE team_id=$1 AND id = " +
		"(SELECT MAX(id) FROM score WHERE team_id=$1)")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(teamID).Scan(&s.ID, &s.Score, &s.Timestamp)
	if err != nil {
		return
	}

	s.TeamID = teamID

	return
}
