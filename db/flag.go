/**
 * @file flag.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2015
 * @brief queries for flag table
 */

package db

import (
	"database/sql"
	"time"
)

// Flag row
type Flag struct {
	ID        int
	TeamID    int
	TaskID    int
	Flag      string
	Solved    bool
	Timestamp time.Time
}

func createFlagTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "flag" (
		id		SERIAL PRIMARY KEY,
		team_id		INTEGER NOT NULL,
		task_id		INTEGER NOT NULL,
		flag		TEXT NOT NULL,
		solved		BOOLEAN NOT NULL,
		timestamp	TIMESTAMP with time zone DEFAULT now()
	)`)

	return
}

// AddFlag add flag to db and fill id
func AddFlag(db *sql.DB, flag *Flag) (err error) {

	stmt, err := db.Prepare("INSERT INTO flag " +
		"(team_id, task_id, flag, solved) " +
		"VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(flag.TeamID, flag.TaskID, flag.Flag,
		flag.Solved).Scan(&flag.ID)
	if err != nil {
		return
	}

	return
}

// GetFlags get all flags in flags table
func GetFlags(db *sql.DB) (flags []Flag, err error) {

	rows, err := db.Query("SELECT id, team_id, task_id, flag, solved, " +
		"timestamp FROM flag")
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var f Flag

		err = rows.Scan(&f.ID, &f.TeamID, &f.TaskID, &f.Flag,
			&f.Solved, &f.Timestamp)
		if err != nil {
			return
		}

		flags = append(flags, f)
	}

	return
}

// GetSolvedCount return amount of solved flags
func GetSolvedCount(db *sql.DB, taskID int) (count int, err error) {

	stmt, err := db.Prepare("SELECT count(*) FROM flag " +
		"WHERE task_id=$1 AND solved=TRUE")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(taskID).Scan(&count)
	if err != nil {
		return
	}

	return
}

// IsSolved return true if task solved by team
func IsSolved(db *sql.DB, teamID, taskID int) (solved bool, err error) {
	stmt, err := db.Prepare("SELECT EXISTS(SELECT id FROM flag " +
		"WHERE team_id=$1 AND task_id=$2 AND solved=TRUE)")

	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(teamID, taskID).Scan(&solved)
	if err != nil {
		return
	}

	return
}

// GetSolvedBy get all team ids who solved task
func GetSolvedBy(db *sql.DB, taskID int) (teamIDs []int, err error) {

	stmt, err := db.Prepare("SELECT team_id FROM flag " +
		"WHERE task_id=$1 AND solved=TRUE")
	if err != nil {
		return
	}

	defer stmt.Close()

	rows, err := stmt.Query(taskID)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var teamID int

		err = rows.Scan(&teamID)
		if err != nil {
			return
		}

		teamIDs = append(teamIDs, teamID)
	}

	return
}
