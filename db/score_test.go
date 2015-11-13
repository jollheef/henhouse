/**
 * @file score_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
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
