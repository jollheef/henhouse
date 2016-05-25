/**
 * @file score_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2015
 * @brief test work with score table
 */

package db

import (
	"errors"
	"testing"
)

func TestCreateScoreTable(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = createScoreTable(db)
	if err != nil {
		panic(err)
	}
}

func TestAddScore(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	score := Score{ID: 255, TeamID: 1, Score: 10}

	err = AddScore(db, &score)
	if err != nil {
		panic(err)
	}

	if score.ID != 1 {
		panic(errors.New("Score id not correct"))
	}
}

// Test add score with closed database
func TestFailAddScore(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	db.Close()

	err = AddScore(db, &Score{})
	if err == nil {
		panic(err)
	}
}

func TestGetLastScore(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	teamID := 1
	resultScore := 30

	err = AddScore(db, &Score{ID: 255, TeamID: teamID, Score: 10})
	if err != nil {
		panic(err)
	}

	err = AddScore(db, &Score{ID: 255, TeamID: teamID, Score: resultScore})
	if err != nil {
		panic(err)
	}

	score, err := GetLastScore(db, teamID)
	if err != nil {
		panic(err)
	}

	if score.Score != resultScore {
		panic(errors.New("Score value not correct"))
	}
}

// Test get last score with closed database
func TestFailGetLastScore(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	db.Close()

	_, err = GetLastScore(db, 1)
	if err == nil {
		panic(err)
	}
}
