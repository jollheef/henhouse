/**
 * @file session_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date November, 2015
 * @brief test work with session table
 */

package db

import (
	"errors"
	"testing"
)

func TestCreateSessionTable(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = createSessionTable(db)
	if err != nil {
		panic(err)
	}
}

func TestAddSession(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	session := Session{ID: 255, TeamID: 10, Session: "test"}

	err = AddSession(db, &session)
	if err != nil {
		panic(err)
	}

	if session.ID != 1 {
		panic(errors.New("Session id not correct"))
	}
}

func TestGetSessionTeam(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	sessionText := "test"
	teamID := 10

	session := Session{ID: 255, TeamID: teamID, Session: sessionText}

	err = AddSession(db, &session)
	if err != nil {
		panic(err)
	}

	id, err := GetSessionTeam(db, sessionText)
	if err != nil {
		panic(err)
	}

	if id != teamID {
		panic(errors.New("Session text not correct"))
	}
}
