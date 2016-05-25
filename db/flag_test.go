/**
 * @file flag.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2015
 * @brief test work with flag table
 */

package db

import (
	"errors"
	"fmt"
	"testing"
)

func TestCreateFlagTable(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = createFlagTable(db)
	if err != nil {
		panic(err)
	}
}

func TestAddFlag(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	flag := Flag{ID: 255, TeamID: 1, TaskID: 1, Flag: "test", Solved: false}

	err = AddFlag(db, &flag)
	if err != nil {
		panic(err)
	}

	if flag.ID != 1 {
		panic(errors.New("Flag id not correct"))
	}
}

// Test add flag with closed database
func TestFailAddFlag(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	db.Close()

	err = AddFlag(db, &Flag{})
	if err == nil {
		panic(err)
	}
}

func TestGetFlags(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	nflags := 150

	for i := 0; i < nflags; i++ {

		flag := Flag{ID: 255, TeamID: 1, TaskID: 1,
			Flag: fmt.Sprintf("%d", i), Solved: false}

		err = AddFlag(db, &flag)
		if err != nil {
			panic(err)
		}
	}

	flags, err := GetFlags(db)
	if err != nil {
		panic(err)
	}

	if len(flags) != nflags {
		panic(errors.New("Mismatch get flags length"))
	}

	for i := 0; i < nflags; i++ {

		if flags[i].Flag != fmt.Sprintf("%d", i) && flags[i].ID != i {
			panic(errors.New("Get invalid flag"))
		}
	}
}

// Test get flags with closed database
func TestFailGetFlags(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	db.Close()

	_, err = GetFlags(db)
	if err == nil {
		panic(err)
	}
}

func TestGetSolvedCount(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	nflags := 150
	teamID := 1
	taskID := 1

	for i := 0; i < nflags; i++ {

		flag := Flag{ID: 255, TeamID: teamID, TaskID: taskID,
			Flag: fmt.Sprintf("%d", i)}

		if i%2 == 0 {
			flag.Solved = true
		} else {
			flag.Solved = false
		}

		err = AddFlag(db, &flag)
		if err != nil {
			panic(err)
		}
	}

	solvedFlags, err := GetSolvedCount(db, taskID)
	if err != nil {
		panic(err)
	}

	if solvedFlags != nflags/2 {
		panic(errors.New("Mismatch solved flags length"))
	}
}

// Test get solved count with closed database
func TestFailGetSolvedCount(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	db.Close()

	_, err = GetSolvedCount(db, 0)
	if err == nil {
		panic(err)
	}
}

func TestIsSolved(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	teamID := 10
	taskID := 15

	flag := Flag{TeamID: teamID, TaskID: taskID, Solved: true}

	err = AddFlag(db, &flag)
	if err != nil {
		panic(err)
	}

	solved, err := IsSolved(db, teamID, taskID)
	if !solved {
		panic("Solved task unsolved")
	}

	if err != nil {
		panic(err)
	}
}

// Test is task solved with closed database
func TestFailIsSolved(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	db.Close()

	_, err = IsSolved(db, 0, 0)
	if err == nil {
		panic(err)
	}
}

// Test is task solved with unsolved task
func TestFailIsSolved2(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	solved, err := IsSolved(db, 0, 0)
	if err != nil {
		panic(err)
	}

	if solved {
		panic("Unsolved task solved")
	}
}
