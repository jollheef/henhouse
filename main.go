/**
 * @file main.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date October, 2015
 * @brief task-based ctf daemon
 *
 * Entry point for task-based ctf daemon
 */

package main

import (
	"database/sql"
	"github.com/jollheef/henhouse/config"
	"github.com/jollheef/henhouse/db"
	"github.com/jollheef/henhouse/game"
	"github.com/jollheef/henhouse/scoreboard"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
)

var (
	configPath = kingpin.Arg("config",
		"Path to configuration file.").Required().String()

	dbReinit = kingpin.Flag("reinit", "Reinit database.").Bool()
)

func reinitDatabase(database *sql.DB, cfg config.Config) (err error) {
	log.Println("Reinit database")

	for _, team := range cfg.Teams {
		log.Println("Add team", team.Name)
		err = db.AddTeam(database, &db.Team{
			Name:  team.Name,
			Desc:  team.Description,
			Login: team.Login,
			Pass:  team.Pass,
		})
		if err != nil {
			return
		}
	}

	entries, err := ioutil.ReadDir(cfg.TaskDir)
	if err != nil {
		return
	}

	var categories []db.Category

	for _, entry := range entries {

		if entry.IsDir() {
			continue
		}

		var content []byte
		content, err = ioutil.ReadFile(cfg.TaskDir + "/" +
			entry.Name())
		if err != nil {
			return
		}

		var task config.Task
		task, err = config.ParseXMLTask(content)
		if err != nil {
			return
		}

		var finded bool
		var taskCategory db.Category
		for _, cat := range categories {
			if cat.Name == task.Category {
				finded = true
				taskCategory = cat
				break
			}
		}

		if !finded {
			taskCategory.Name = task.Category

			err = db.AddCategory(database, &taskCategory)
			if err != nil {
				return
			}

			categories = append(categories, taskCategory)

			log.Println("Add category", taskCategory.Name)
		}

		err = db.AddTask(database, &db.Task{
			Name:          task.Name,
			Desc:          task.Description,
			CategoryID:    taskCategory.ID,
			Level:         task.Level,
			Flag:          task.Flag,
			Price:         500,   // TODO support non-shared task
			Shared:        true,  // TODO support non-shared task
			MaxSharePrice: 500,   // TODO support value from xml
			MinSharePrice: 100,   // TODO support value from xml
			Opened:        false, // by default task is closed
		})

		log.Println("Add task", task.Name)

		if err != nil {
			return
		}
	}

	return
}

func main() {

	kingpin.Parse()

	cfg, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatalln("Cannot open config:", err)
	}

	log.Println("Use db connection", cfg.Database.Connection)

	var database *sql.DB

	if *dbReinit {

		database, err = db.InitDatabase(cfg.Database.Connection)
		if err != nil {
			log.Fatalln("Error:", err)
		}

		defer database.Close()

		err = reinitDatabase(database, cfg)
		if err != nil {
			log.Fatalln("Error:", err)
		}

	} else {

		database, err = db.OpenDatabase(cfg.Database.Connection)
		if err != nil {
			log.Fatalln("Error:", err)
		}

		defer database.Close()
	}

	log.Println("Set max db connections to", cfg.Database.MaxConnections)
	database.SetMaxOpenConns(cfg.Database.MaxConnections)

	log.Println("Start game at", cfg.Game.Start.Time)
	log.Println("End game at", cfg.Game.End.Time)
	game, err := game.NewGame(database, cfg.Game.Start.Time,
		cfg.Game.End.Time)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	log.Println("Set task open timeout to", cfg.Task.OpenTimeout.Duration)
	game.OpenTimeout = cfg.Task.OpenTimeout.Duration

	if cfg.Task.AutoOpen {
		log.Println("Auto open tasks after", cfg.Task.AutoOpenTimeout.Duration)
	} else {
		log.Println("Auto open tasks disabled")
	}

	game.AutoOpen = cfg.Task.AutoOpen
	game.AutoOpenTimeout = cfg.Task.AutoOpenTimeout.Duration

	go game.Run()

	log.Println("Use html files from", cfg.Scoreboard.WwwPath)
	log.Println("Listen at", cfg.Scoreboard.Addr)
	err = scoreboard.Scoreboard(database, &game, cfg.Scoreboard.WwwPath,
		cfg.Scoreboard.Addr)
	if err != nil {
		log.Fatalln("Error:", err)
	}
}
