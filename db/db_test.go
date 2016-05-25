/**
 * @file db_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date October, 2015
 * @brief test general work with database functions
 */

package db

import (
	"crypto/rand"
	"errors"
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

// Test create schema with closed database
func TestFailCreateSchema(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	db.Close()

	err = createSchema(db)
	if err == nil {
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

// Test init database with incorrect path
func TestFailInitDatabase(*testing.T) {

	rndStr, err := rndString(15) // 15 is random len
	if err != nil {
		panic(err)
	}

	_, err = InitDatabase(fmt.Sprintf("user=%s  sslmode=disable", rndStr))
	if err == nil {
		panic(err)
	}
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

// Test clean database with closed database
func TestFailCleanDatabase(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	db.Close()

	err = CleanDatabase(db)
	if err == nil {
		panic(err)
	}
}

// Test clean database with removed sequence
func TestFailCleanDatabase2(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	for _, table := range tables {
		_, err = db.Exec("DROP SEQUENCE " + table + "_id_seq CASCADE;")
		if err != nil {
			panic(err)
		}
	}

	err = CleanDatabase(db)
	if err == nil {
		panic(err)
	}
}

func TestTablesCount(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	ntables := 0
	db.QueryRow("SELECT count(*) FROM information_schema.tables " +
		" WHERE table_schema = 'public';").Scan(&ntables)

	if ntables != len(tables) {
		panic(errors.New("Invalid table list"))
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
