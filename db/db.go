/**
 * @file db.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date October, 2015
 * @brief work with database
 *
 * Contain functions for work with database.
 */

package db

import (
	"database/sql"
	_ "github.com/lib/pq" // import postgresql db engine
)

// All table names
var tables = [...]string{"category", "flag"}

// Create tables
func createSchema(db *sql.DB) (err error) {

	_, err = db.Exec("CREATE SCHEMA IF NOT EXISTS public")
	if err != nil {
		return
	}

	err = createCategoryTable(db)
	if err != nil {
		return
	}

	err = createFlagTable(db)
	if err != nil {
		return
	}

	return
}

// OpenDatabase need defer db.Close() after open
func OpenDatabase(path string) (db *sql.DB, err error) {

	db, err = sql.Open("postgres", path)
	if err != nil {
		return
	}

	err = createSchema(db)
	if err != nil {
		return
	}

	return
}

// Clean all values in table
func cleanTable(db *sql.DB, table string) (err error) {
	_, err = db.Exec("DELETE FROM " + table)
	return
}

// Restart id sequence in table
func restartSequence(db *sql.DB, table string) (err error) {
	_, err = db.Exec("ALTER SEQUENCE " + table + "_id_seq RESTART WITH 1;")
	return
}

// CleanDatabase clean all values and restart sequences in database without
// drop tables
func CleanDatabase(db *sql.DB) (err error) {

	for _, table := range tables {

		err = cleanTable(db, table)
		if err != nil {
			return
		}

		err = restartSequence(db, table)
		if err != nil {
			return
		}
	}

	return
}

// Drop public schema in current database
func dropSchema(db *sql.DB) (err error) {

	_, err = db.Exec("DROP SCHEMA public CASCADE")
	if err != nil {
		return
	}
	return
}

// InitDatabase recreate all database tables
func InitDatabase(path string) (db *sql.DB, err error) {

	db, err = OpenDatabase(path)
	if err != nil {
		return
	}

	err = dropSchema(db)
	if err != nil {
		// No schema not good, but ok
	}

	err = createSchema(db)
	if err != nil {
		return
	}

	return
}
