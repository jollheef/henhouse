/**
 * @file team.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date November, 2015
 * @brief queries for team table
 */

package db

import (
	"database/sql"
)

// Team row
type Team struct {
	ID    int
	Name  string
	Email string
	Desc  string
	Login string
	Pass  string
	Salt  string
}

func createTeamTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "team" (
		id		SERIAL PRIMARY KEY,
		name		TEXT NOT NULL,
		email		TEXT NOT NULL,
		description	TEXT NOT NULL,
		login		TEXT NOT NULL,
		password	TEXT NOT NULL,
		salt		TEXT NOT NULL
	)`)

	return
}

// AddTeam add team and fill id
func AddTeam(db *sql.DB, t *Team) (err error) {

	stmt, err := db.Prepare("INSERT INTO team (name, email, " +
		"description, login, password, salt) " +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(t.Name, t.Email, t.Desc, t.Login, t.Pass,
		t.Salt).Scan(&t.ID)
	if err != nil {
		return
	}

	return
}

// GetTeams get all teams
func GetTeams(db *sql.DB) (teams []Team, err error) {

	rows, err := db.Query("SELECT id, name, email, description, login, " +
		"password, salt FROM team")
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var t Team

		err = rows.Scan(&t.ID, &t.Name, &t.Email, &t.Desc, &t.Login,
			&t.Pass, &t.Salt)
		if err != nil {
			return
		}

		teams = append(teams, t)
	}

	return
}
