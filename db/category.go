/**
 * @file category.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2015
 * @brief queries for category table
 */

package db

import "database/sql"

// Category row
type Category struct {
	ID   int
	Name string
}

func createCategoryTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "category" (
		id	SERIAL PRIMARY KEY,
		name	TEXT NOT NULL UNIQUE
	)`)

	return
}

// AddCategory add category to db and fill Id
func AddCategory(db *sql.DB, category *Category) (err error) {

	stmt, err := db.Prepare("INSERT INTO category (name) " +
		"VALUES ($1) RETURNING id")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(category.Name).Scan(&category.ID)
	if err != nil {
		return
	}

	return
}

// GetCategories get all categories in category table
func GetCategories(db *sql.DB) (categories []Category, err error) {

	rows, err := db.Query("SELECT id, name FROM category")
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var category Category

		err = rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return
		}

		categories = append(categories, category)
	}

	return
}
