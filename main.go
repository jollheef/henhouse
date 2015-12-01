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
	"github.com/jollheef/henhouse/config"
	"github.com/jollheef/henhouse/db"
	"github.com/jollheef/henhouse/game"
	"github.com/jollheef/henhouse/scoreboard"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
)

var (
	configPath = kingpin.Arg("config",
		"Path to configuration file.").Required().String()

	dbReinit = kingpin.Flag("reinit", "Reinit database.").Bool()
)

func main() {

	kingpin.Parse()

	cfg, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatalln("Cannot open config:", err)
	}

	log.Println("Use db connection", cfg.Database.Connection)

	if *dbReinit {
		log.Println("Reinit database")

		database, err := db.InitDatabase(cfg.Database.Connection)
		if err != nil {
			log.Fatalln("Error:", err)
		}

		defer database.Close()

		// TODO add teams
		// TODO add categories from xml
		// TODO add tasks from xml
	} else {
		database, err := db.OpenDatabase(cfg.Database.Connection)
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

	log.Print("Set task open timeout to", cfg.Task.OpenTimeout.Duration)
	game.OpenTimeout = cfg.Task.OpenTimeout.Duration

	if cfg.Task.AutoOpen {
		log.Print("Auto open tasks after", cfg.Task.AutoOpenTimeout.Duration)
	} else {
		log.Print("Auto open tasks disabled")
	}

	game.AutoOpen = cfg.Task.AutoOpen
	game.AutoOpenTimeout = cfg.Task.AutoOpenTimeout.Duration

	go game.Run()

	log.Println("Use html files from", cfg.Scoreboard.WwwPath)
	log.Println("Listen at", cfg.Scoreboard.Addr)
	err = scoreboard.Scoreboard(&game, cfg.Scoreboard.WwwPath,
		cfg.Scoreboard.Addr)
	if err != nil {
		log.Fatalln("Error:", err)
	}
}
