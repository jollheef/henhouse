/**
 * @file task.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date November, 2015
 * @brief queries for task table
 */

package db

import (
	"database/sql"
	"time"
)

// Task row
type Task struct {
	ID            int
	Name          string
	Desc          string
	CategoryID    int
	Level         int
	Price         int
	Shared        bool
	Flag          string
	MaxSharePrice int
	MinSharePrice int
	Opened        bool
	Author        string
	OpenedTime    time.Time
}

func createTaskTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "task" (
		id		SERIAL PRIMARY KEY,
		name		TEXT NOT NULL,
		description	TEXT NOT NULL,
		category_id	INTEGER NOT NULL,
		level		INTEGER NOT NULL,
		price		INTEGER NOT NULL,
		shared		BOOLEAN NOT NULL,
		flag		TEXT NOT NULL,
		max_share_price	INTEGER NOT NULL,
		min_share_price	INTEGER NOT NULL,
		opened		BOOLEAN NOT NULL,
		author		TEXT NOT NULL,
		opened_time	TIMESTAMP with time zone
	)`)

	return
}

// AddTask add task and fill id
func AddTask(db *sql.DB, t *Task) (err error) {

	stmt, err := db.Prepare("INSERT INTO task (name, description, " +
		"category_id, level, price, shared, flag, max_share_price, " +
		"min_share_price, opened, author, opened_time) " +
		"VALUES ($1, $2, $3, $4, " +
		"$5, $6, $7, $8, $9, $10, $11, $12) RETURNING id")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(t.Name, t.Desc, t.CategoryID, t.Level, t.Price,
		t.Shared, t.Flag, t.MaxSharePrice, t.MinSharePrice,
		t.Opened, t.Author, t.OpenedTime).Scan(&t.ID)
	if err != nil {
		return
	}

	return
}

// GetTasks get all tasks in tasks table
func GetTasks(db *sql.DB) (tasks []Task, err error) {

	rows, err := db.Query("SELECT id, name, description, category_id, " +
		"level, price, shared, flag, max_share_price, " +
		"min_share_price, opened, author, opened_time FROM task")
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var t Task

		err = rows.Scan(&t.ID, &t.Name, &t.Desc, &t.CategoryID,
			&t.Level, &t.Price, &t.Shared, &t.Flag,
			&t.MaxSharePrice, &t.MinSharePrice, &t.Opened,
			&t.Author, &t.OpenedTime)
		if err != nil {
			return
		}

		tasks = append(tasks, t)
	}

	return
}

// SetOpened open or close task
func SetOpened(db *sql.DB, taskID int, opened bool) (err error) {

	stmt, err := db.Prepare("UPDATE task SET opened=$1, opened_time=$2 " +
		"WHERE id=$3")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(opened, time.Now(), taskID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateTask update task
func UpdateTask(db *sql.DB, t *Task) (err error) {

	stmt, err := db.Prepare("UPDATE task SET name=$1, description=$2, " +
		"category_id=$3, level=$4, price=$5, shared=$6, flag=$7, " +
		"max_share_price=$8, min_share_price=$9, opened=$10, " +
		"author=$11, opened_time=$12 WHERE id=$13")
	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(t.Name, t.Desc, t.CategoryID, t.Level, t.Price,
		t.Shared, t.Flag, t.MaxSharePrice, t.MinSharePrice, t.Opened,
		t.Author, t.OpenedTime, t.ID)
	if err != nil {
		return
	}

	return
}

// GetTask get task by id
func GetTask(db *sql.DB, taskID int) (t Task, err error) {

	stmt, err := db.Prepare("SELECT id, name, description, category_id, " +
		"level, price, shared, flag, max_share_price, " +
		"min_share_price, opened, author, opened_time " +
		"FROM task WHERE id=$1")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(taskID).Scan(&t.ID, &t.Name, &t.Desc,
		&t.CategoryID, &t.Level, &t.Price, &t.Shared, &t.Flag,
		&t.MaxSharePrice, &t.MinSharePrice, &t.Opened,
		&t.Author, &t.OpenedTime)
	if err != nil {
		return
	}

	return
}
