/**
 * @file db_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date October, 2015
 * @brief test general work with database functions
 */

package db

import (
	"crypto/rand"
	"fmt"
	"testing"
)

const dbPath string = "user=postgres dbname=henhouse_test sslmode=disable"

func rndString(len int) (str string, err error) {

	randBuf := make([]byte, len)

	_, err = rand.Read(randBuf)
	if err != nil {
		return
	}

	str = string(randBuf)

	return
}

func TestOpenDatabase(*testing.T) {

	db, err := OpenDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	// Error in path may not show in open
	_, err = db.Exec("SELECT 1")
	if err != nil {
		panic(err)
	}

	db.Close()
}

func TestFailOpenDatabase(*testing.T) {

	rndStr, err := rndString(15) // 15 is random len
	if err != nil {
		panic(err)
	}

	db, err := OpenDatabase(fmt.Sprintf("user=%s  sslmode=disable", rndStr))
	if err == nil {
		panic(err)
	}

	// Error in path may not show in open
	_, err = db.Exec("SELECT 1")
	if err == nil {
		panic(err)
	}

	defer db.Close()
}

func TestCreateSchema(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = createSchema(db)
	if err != nil {
		panic(err)
	}
}

func TestInitDatabase(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	db.Close()
}

func TestCleanDatabase(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = CleanDatabase(db)
	if err != nil {
		panic(err)
	}
}

func TestFailCleanTable(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	rndStr, err := rndString(15) // 15 is random len
	if err != nil {
		panic(err)
	}

	// Invalid table name
	err = cleanTable(db, rndStr)
	if err == nil {
		panic(err)
	}
}

func TestFailRestartSequence(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	rndStr, err := rndString(15) // 15 is random len
	if err != nil {
		panic(err)
	}

	// Invalid table name
	err = restartSequence(db, rndStr)
	if err == nil {
		panic(err)
	}
}

func TestFailDropSchema(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	_, err = db.Exec("DROP SCHEMA public CASCADE")
	if err != nil {
		panic(err)
	}

	// Drop dropped schema
	err = dropSchema(db)
	if err == nil {
		// must be fail
		panic(err)
	}
}
