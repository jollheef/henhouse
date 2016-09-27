/**
 * @file task_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
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

	task := Task{ID: 255}

	err = AddTask(db, &task)
	if err != nil {
		panic(err)
	}

	if task.ID != 1 {
		panic(errors.New("Task id not correct"))
	}
}

// Test work with task at closed database
func TestOnClosedDatabase(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	db.Close()

	err = AddTask(db, &Task{})
	if err == nil {
		panic(err)
	}

	err = UpdateTask(db, &Task{})
	if err == nil {
		panic(err)
	}

	_, err = GetTask(db, 0)
	if err == nil {
		panic(err)
	}

	err = SetOpened(db, 0, false)
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

func TestUpdateTask(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	taskName := "__rand_task_100"

	task := Task{ID: 255, Name: taskName}

	err = AddTask(db, &task)
	if err != nil {
		panic(err)
	}

	newTaskName := "100__rand_task"

	err = UpdateTask(db, &Task{ID: task.ID, Name: newTaskName})
	if err != nil {
		panic(err)
	}

	t, err := GetTask(db, task.ID)
	if err != nil {
		panic(err)
	}

	if t.Name != newTaskName {
		panic("invalid task name")
	}
}

func TestGetTask(*testing.T) {

	db, err := InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	taskName := "__rand_task_1"

	task := Task{ID: 255, Name: taskName}

	err = AddTask(db, &task)
	if err != nil {
		panic(err)
	}

	t, err := GetTask(db, task.ID)
	if err != nil {
		panic(err)
	}

	if t.Name != taskName {
		panic("invalid task name")
	}
}
