/**
 * @file task_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date November, 2015
 * @brief test work with task table
 */

package db

import (
	"errors"
	"fmt"
	"testing"
)

func TestCreateTaskTable(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = createTaskTable(db)
	if err != nil {
		panic(err)
	}
}

func TestAddTask(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	task := Task{255, "n", "d", 10, 100, 50, true, "f", 1, 10, true}

	err = AddTask(db, &task)
	if err != nil {
		panic(err)
	}

	if task.ID != 1 {
		panic(errors.New("Task id not correct"))
	}
}

// Test add task with closed database
func TestFailAddTask(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	db.Close()

	err = AddTask(db, &Task{})
	if err == nil {
		panic(err)
	}
}

func TestGetTasks(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	ntasks := 150

	for i := 0; i < ntasks; i++ {

		task := Task{ID: 255, Name: fmt.Sprintf("%d", i)}

		err = AddTask(db, &task)
		if err != nil {
			panic(err)
		}
	}

	tasks, err := GetTasks(db)
	if err != nil {
		panic(err)
	}

	if len(tasks) != ntasks {
		panic(errors.New("Mismatch get tasks length"))
	}

	for i := 0; i < ntasks; i++ {

		if tasks[i].Name != fmt.Sprintf("%d", i) && tasks[i].ID != i {
			panic(errors.New("Get invalid task"))
		}
	}
}

// Test get tasks with closed database
func TestFailGetTasks(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	db.Close()

	_, err = GetTasks(db)
	if err == nil {
		panic(err)
	}
}

func TestSetOpened(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	task := Task{Opened: false}

	err = AddTask(db, &task)
	if err != nil {
		panic(err)
	}

	tasks, err := GetTasks(db)
	if err != nil {
		panic(err)
	}
	if tasks[0].Opened != false {
		panic("closed task added as opened")
	}

	err = SetOpened(db, task.ID, true)
	if err != nil {
		panic(err)
	}

	tasks, err = GetTasks(db)
	if err != nil {
		panic(err)
	}
	if tasks[0].Opened != true {
		panic("opened task closed")
	}

	err = SetOpened(db, task.ID, false)
	if err != nil {
		panic(err)
	}

	tasks, err = GetTasks(db)
	if err != nil {
		panic(err)
	}
	if tasks[0].Opened != false {
		panic("closed task opened")
	}
}
