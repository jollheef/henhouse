/**
 * @file main.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date October, 2015
 * @brief task-based ctf daemon
 *
 * Entry point for task-based ctf daemon
 */

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/jollheef/henhouse/config"
	"github.com/jollheef/henhouse/db"
	"github.com/jollheef/henhouse/game"
	"github.com/jollheef/henhouse/scoreboard"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	configPath = kingpin.Arg("config",
		"Path to configuration file.").Required().String()

	dbReinit = kingpin.Flag("reinit", "Reinit database.").Bool()
)

var (
	// CommitID fill in ./build.sh
	CommitID string
	// BuildDate fill in ./build.sh
	BuildDate string
	// BuildTime fill in ./build.sh
	BuildTime string
)

func fillTranslateFallback(task * config.Task) {
	if task.NameEn == "" {
		task.NameEn = task.Name
	}
	if task.Name == "" {
		task.Name = task.NameEn
	}
	if task.DescriptionEn == "" {
		task.DescriptionEn = task.Description
	}
	if task.Description == "" {
		task.Description = task.DescriptionEn
	}
	return
}

func reinitDatabase(database *sql.DB, cfg config.Config) (err error) {
	log.Println("Reinit database")

	for _, team := range cfg.Teams {
		log.Println("Add team", team.Name)
		err = db.AddTeam(database, &db.Team{
			Name:  team.Name,
			Desc:  team.Description,
			Token: team.Token,
			Test:  team.Test,
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

		fillTranslateFallback(&task)

		err = db.AddTask(database, &db.Task{
			Name:          task.Name,
			Desc:          task.Description,
			NameEn:        task.NameEn,
			DescEn:        task.DescriptionEn,
			Tags:          task.Tags,
			CategoryID:    taskCategory.ID,
			Level:         task.Level,
			Flag:          task.Flag,
			Price:         500,   // TODO support non-shared task
			Shared:        true,  // TODO support non-shared task
			MaxSharePrice: 500,   // TODO support value from xml
			MinSharePrice: 100,   // TODO support value from xml
			Opened:        false, // by default task is closed
			Author:        task.Author,
			ForceClosed:   task.ForceClosed,
		})

		log.Println("Add task", task.Name)

		if err != nil {
			return
		}
	}

	return
}

func checkTaskPrices(cfg *config.Config)(err error){
	if cfg.TaskPrice.P200 == 0 || cfg.TaskPrice.P300 == 0 ||
		cfg.TaskPrice.P400 == 0 || cfg.TaskPrice.P500 == 0 {
		err = errors.New("Error: Task price not setted")
	}
	return
}

func initGame(database *sql.DB, cfg config.Config) (err error) {

	var teamBase float64

	if cfg.TaskPrice.UseNonLinear {
		teamBase, err = game.CalcTeamsBase(database)
		if err != nil {
			return
		}
		log.Println("Use teams amount based on session counter")
	} else if cfg.TaskPrice.UseTeamsBase {
		teamBase = float64(cfg.TaskPrice.TeamsBase)
		log.Println("Set teams base to", cfg.TaskPrice.TeamsBase)
	} else {
		teamBase = float64(len(cfg.Teams))
		log.Println("Use teams amount as teams base")
	}

	log.Println("Start game at", cfg.Game.Start.Time)
	log.Println("End game at", cfg.Game.End.Time)
	g, err := game.NewGame(database, cfg.Game.Start.Time,
		cfg.Game.End.Time, teamBase)
	if err != nil {
		return
	}

	if cfg.TaskPrice.UseNonLinear {
		go g.TeamsBaseUpdater(database,
			cfg.Scoreboard.RecalcTimeout.Duration)
	}

	err = checkTaskPrices(&cfg)
	if err != nil{
		return
	}

	fmt := "Set task price %d if solved less than %d%%\n"
	log.Printf(fmt, 200, cfg.TaskPrice.P200)
	log.Printf(fmt, 300, cfg.TaskPrice.P300)
	log.Printf(fmt, 400, cfg.TaskPrice.P400)
	log.Printf(fmt, 500, cfg.TaskPrice.P500)

	g.SetTaskPrice(cfg.TaskPrice.P500, cfg.TaskPrice.P400,
		cfg.TaskPrice.P300, cfg.TaskPrice.P200)

	log.Println("Set task open timeout to", cfg.Task.OpenTimeout.Duration)
	g.OpenTimeout = cfg.Task.OpenTimeout.Duration

	if cfg.Task.AutoOpen {
		log.Println("Auto open tasks after",
			cfg.Task.AutoOpenTimeout.Duration)
	} else {
		log.Println("Auto open tasks disabled")
	}

	g.AutoOpen = cfg.Task.AutoOpen
	g.AutoOpenTimeout = cfg.Task.AutoOpenTimeout.Duration

	go g.Run()

	infoD := cfg.WebsocketTimeout.Info.Duration
	if infoD != 0 {
		scoreboard.InfoTimeout = infoD
	}
	log.Println("Update info timeout:", scoreboard.InfoTimeout)

	scoreboardD := cfg.WebsocketTimeout.Scoreboard.Duration
	if scoreboardD != 0 {
		scoreboard.ScoreboardTimeout = scoreboardD
	}
	log.Println("Update scoreboard timeout:", scoreboard.ScoreboardTimeout)

	tasksD := cfg.WebsocketTimeout.Tasks.Duration
	if tasksD != 0 {
		scoreboard.TasksTimeout = tasksD
	}
	log.Println("Update tasks timeout:", scoreboard.TasksTimeout)

	flagSendD := cfg.Flag.SendTimeout.Duration
	if flagSendD != 0 {
		scoreboard.FlagTimeout = flagSendD
	}
	log.Println("Flag timeout:", scoreboard.FlagTimeout)

	scoreboardRecalcD := cfg.Scoreboard.RecalcTimeout.Duration
	if scoreboardRecalcD != 0 {
		scoreboard.ScoreboardRecalcTimeout = scoreboardRecalcD
	}

	log.Println("Score recalc timeout:", scoreboard.ScoreboardRecalcTimeout)

	log.Println("Use html files from", cfg.Scoreboard.WwwPath)
	log.Println("Listen at", cfg.Scoreboard.Addr)
	err = scoreboard.Scoreboard(database, &g,
		cfg.Scoreboard.WwwPath,
		cfg.Scoreboard.TemplatePath,
		cfg.Scoreboard.Addr)

	return
}

func main() {

	if len(CommitID) > 7 {
		CommitID = CommitID[:7] // abbreviated commit hash
	}

	version := BuildDate + " " + CommitID +
		" (Mikhail Klementyev <jollheef@riseup.net>)"

	kingpin.Version(version)

	kingpin.Parse()

	fmt.Println(version)

	cfg, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatalln("Cannot open config:", err)
	}

	logFile, err := os.OpenFile(cfg.LogFile,
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Cannot open file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println(version)

	var rlim syscall.Rlimit
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim)
	if err != nil {
		log.Fatalln("Getrlimit fail:", err)
	}

	log.Println("RLIMIT_NOFILE CUR:", rlim.Cur, "MAX:", rlim.Max)

	var database *sql.DB

	if *dbReinit {

		if cfg.Database.SafeReinit {
			if time.Now().After(cfg.Game.Start.Time) {
				log.Fatalln("Reinit after start not allowed")
			}
		}

		database, err = db.InitDatabase(cfg.Database.Connection)
		if err != nil {
			log.Fatalln("Error:", err)
		}

		err = db.CleanDatabase(database)
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

	err = initGame(database, cfg)
	if err != nil {
		log.Fatalln("Error:", err)
	}
}
