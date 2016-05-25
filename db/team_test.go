/**
 * @file team_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2015
 * @brief test work with team table
 */

package db

import (
	"errors"
	"fmt"
	"testing"
)

func TestCreateTeamTable(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = createTeamTable(db)
	if err != nil {
		panic(err)
	}
}

func TestAddTeam(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	team := Team{255, "n", "e", "d", "l", false}

	err = AddTeam(db, &team)
	if err != nil {
		panic(err)
	}

	if team.ID != 1 {
		panic(errors.New("Team id not correct"))
	}
}

// Test add team with closed database
func TestFailAddTeam(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	db.Close()

	err = AddTeam(db, &Team{})
	if err == nil {
		panic(err)
	}
}

func TestGetTeams(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	nteams := 150

	for i := 0; i < nteams; i++ {

		team := Team{255, fmt.Sprintf("%d", i),
			"e", "d", "l", false}

		err = AddTeam(db, &team)
		if err != nil {
			panic(err)
		}
	}

	teams, err := GetTeams(db)
	if err != nil {
		panic(err)
	}

	if len(teams) != nteams {
		panic(errors.New("Mismatch get teams length"))
	}

	for i := 0; i < nteams; i++ {

		if teams[i].Name != fmt.Sprintf("%d", i) && teams[i].ID != i {
			panic(errors.New("Get invalid team"))
		}
	}
}

// Test get teams with closed database
func TestFailGetTeams(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	db.Close()

	_, err = GetTeams(db)
	if err == nil {
		panic(err)
	}
}
