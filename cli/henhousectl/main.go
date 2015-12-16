/**
 * @file main.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date December, 2015
 * @brief contest checking system CLI
 *
 * Entry point for contest checking system CLI
 */

package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/jollheef/henhouse/config"
	"github.com/jollheef/henhouse/db"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

var (
	configPath = kingpin.Flag("config", "Path to configuration file.").String()

	// Task
	task = kingpin.Command("task", "Work with tasks.")

	taskList = task.Command("list", "List tasks.")

	taskUpdate    = task.Command("update", "Update task.")
	taskUpdateID  = taskUpdate.Arg("id", "ID of task.").Required().Int()
	taskUpdateXML = taskUpdate.Arg("xml", "Path to xml.").Required().String()

	taskOpen   = task.Command("open", "Open task.")
	taskOpenID = taskOpen.Arg("id", "ID of task").Required().Int()

	taskClose   = task.Command("close", "Close task.")
	taskCloseID = taskClose.Arg("id", "ID of task").Required().Int()

	taskDump   = task.Command("dump", "Dump task to xml.")
	taskDumpID = taskDump.Arg("id", "ID of task").Required().Int()

	// Category
	category = kingpin.Command("category", "Work with categories.")

	categoryList = category.Command("list", "List categories.")

	categoryAdd  = category.Command("add", "Add category.")
	categoryName = categoryAdd.Arg("name", "Name.").Required().String()
)

func getCategoryByID(categoryID int, categories []db.Category) string {
	for _, cat := range categories {
		if cat.ID == categoryID {
			return cat.Name
		}
	}
	return "Unknown"
}

func getCategoryByName(name string, categories []db.Category) (id int, err error) {
	for _, cat := range categories {
		if cat.Name == name {
			return cat.ID, nil
		}
	}

	return 0, errors.New("Category " + name + " not found")
}

func taskRow(task db.Task, categories []db.Category) (row []string) {
	row = append(row, fmt.Sprintf("%d", task.ID))
	row = append(row, task.Name)
	row = append(row, getCategoryByID(task.CategoryID, categories))
	row = append(row, task.Flag)
	row = append(row, fmt.Sprintf("%v", task.Opened))
	return
}

type byID []db.Task

func (t byID) Len() int           { return len(t) }
func (t byID) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t byID) Less(i, j int) bool { return t[i].ID < t[j].ID }

func parseTask(path string, categories []db.Category) (t db.Task, err error) {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	task, err := config.ParseXMLTask(content)
	if err != nil {
		return
	}

	t.Name = task.Name
	t.Desc = task.Description
	t.CategoryID, err = getCategoryByName(task.Category, categories)
	if err != nil {
		return
	}

	t.Level = task.Level
	t.Flag = task.Flag
	t.Price = 500         // TODO support non-shared task
	t.Shared = true       // TODO support non-shared task
	t.MaxSharePrice = 500 // TODO support value from xml
	t.MinSharePrice = 100 // TODO support value from xml
	t.Opened = false      // by default task is closed
	t.Author = task.Author

	return
}

func main() {

	kingpin.Parse()

	var cfgPath string

	if *configPath != "" {
		cfgPath = *configPath
	} else {
		cfgPath = "/etc/henhouse/cli.toml"
	}

	cfg, err := config.ReadConfig(cfgPath)
	if err != nil {
		log.Fatalln("Cannot open config:", err)
	}

	database, err := db.OpenDatabase(cfg.Database.Connection)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	defer database.Close()

	database.SetMaxOpenConns(cfg.Database.MaxConnections)

	categories, err := db.GetCategories(database)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	switch kingpin.Parse() {
	case "task update":
		task, err := db.GetTask(database, *taskUpdateID)
		if err != nil {
			log.Fatalln("Error:", err)
		}

		id := task.ID

		task, err = parseTask(*taskUpdateXML, categories)
		if err != nil {
			log.Fatalln("Error:", err)
		}

		task.ID = id

		err = db.UpdateTask(database, &task)
		if err != nil {
			log.Fatalln("Error:", err)
		}

	case "task list":
		tasks, err := db.GetTasks(database)
		if err != nil {
			log.Fatalln("Error:", err)
		}

		sort.Sort(byID(tasks))

		table := tablewriter.NewWriter(os.Stdout)
		header := []string{"ID", "Name", "Category", "Flag", "Opened"}
		table.SetHeader(header)

		for _, task := range tasks {
			table.Append(taskRow(task, categories))
		}

		table.Render()

	case "task open":
		err = db.SetOpened(database, *taskOpenID, true)
		if err != nil {
			log.Fatalln("Error:", err)
		}

	case "task close":
		err = db.SetOpened(database, *taskCloseID, false)
		if err != nil {
			log.Fatalln("Error:", err)
		}

	case "task dump":
		task, err := db.GetTask(database, *taskDumpID)
		if err != nil {
			log.Fatalln("Error:", err)
		}

		xmlTask := config.Task{
			Name:        task.Name,
			Description: task.Desc,
			Category:    getCategoryByID(task.CategoryID, categories),
			Level:       task.Level,
			Flag:        task.Flag,
			Author:      task.Author,
		}

		output, err := xml.MarshalIndent(xmlTask, "", "	")
		if err != nil {
			log.Fatalln("Error:", err)
		}

		os.Stdout.Write(output)

	case "category add":
		err = db.AddCategory(database, &db.Category{Name: *categoryName})
		if err != nil {
			log.Fatalln("Error:", err)
		}

	case "category list":
		categories, err := db.GetCategories(database)
		if err != nil {
			log.Fatalln("Error:", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name"})

		for _, cat := range categories {
			row := []string{fmt.Sprintf("%d", cat.ID), cat.Name}
			table.Append(row)
		}

		table.Render()
	}
}
