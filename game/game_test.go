/**
 * @file game_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date November, 2015
 * @brief test game package
 */

package game

import (
	"fmt"
	"github.com/jollheef/henhouse/db"
	"testing"
	"time"
)

const dbPath string = "user=postgres dbname=henhouse_test sslmode=disable"

// TestNewGame test new game
func TestNewGame(*testing.T) {

	database, err := db.InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer database.Close()

	_, err = NewGame(database, time.Now(), time.Now().Add(time.Hour))
	if err != nil {
		panic(err)
	}
}

// TestNewGameFail test new game with closed database
func TestNewGameFail(*testing.T) {

	database, err := db.InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	database.Close()

	_, err = NewGame(database, time.Now(), time.Now().Add(time.Hour))
	if err == nil {
		panic("work at closed database")
	}
}

func TestTasks(*testing.T) {

	database, err := db.InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer database.Close()

	ntasks := 30

	testCategory := db.Category{Name: "test"}

	err = db.AddCategory(database, &testCategory)
	if err != nil {
		panic(err)
	}

	for i := 0; i < ntasks; i++ {

		task := db.Task{
			ID:         255,
			Name:       fmt.Sprintf("%d", i),
			CategoryID: testCategory.ID,
		}

		err = db.AddTask(database, &task)
		if err != nil {
			panic(err)
		}
	}

	game, err := NewGame(database, time.Now(), time.Now().Add(time.Hour))
	if err != nil {
		panic(err)
	}

	cats, err := game.Tasks()

	ntasksReal := 0

	for _, catInfo := range cats {

		if catInfo.Name != testCategory.Name {
			panic("Get invalid category")
		}

		for n, taskInfo := range catInfo.TasksInfo {

			if taskInfo.Name != fmt.Sprintf("%d", n) {
				panic("Get invalid task")
			}

			ntasksReal++
		}
	}

	if ntasks != ntasksReal {
		panic("Mismatch get tasks length")
	}
}

func TestScoreboard(*testing.T) {

	database, err := db.InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer database.Close()

	validFlag := "testflag"

	nteams := 20

	for i := 0; i < nteams; i++ {

		team := db.Team{255, fmt.Sprintf("team%d", i),
			"e", "d", "l", "p", "s"}

		err = db.AddTeam(database, &team)
		if err != nil {
			panic(err)
		}
	}

	ncategories := 5

	for i := 0; i < ncategories; i++ {

		category := db.Category{Name: fmt.Sprintf("category%d", i)}

		err = db.AddCategory(database, &category)
		if err != nil {
			panic(err)
		}

		ntasks := 5

		for i := 0; i < ntasks; i++ {

			task := db.Task{
				Name:          fmt.Sprintf("task%d", i),
				Flag:          validFlag,
				CategoryID:    category.ID,
				Price:         500,
				MaxSharePrice: 500,
				MinSharePrice: 100,
				Shared:        true,
				Opened:        false,
			}

			err = db.AddTask(database, &task)
			if err != nil {
				panic(err)
			}
		}
	}

	game, err := NewGame(database, time.Now(), time.Now().Add(time.Hour))
	if err != nil {
		panic(err)
	}

	scores, err := game.Scoreboard()
	if err != nil {
		panic(err)
	}

	for _, teamScoreInfo := range scores {
		if teamScoreInfo.Score != 0 {
			panic("score at game start not zero")
		}
	}

	for teamID := 1; teamID <= nteams; teamID++ {
		solved, err := game.Solve(teamID, 1, validFlag)
		if err != nil {
			panic(err)
		}
		if !solved {
			panic("solve task failed")
		}
	}

	err = game.RecalcScoreboard()
	if err != nil {
		panic(err)
	}

	scores, err = game.Scoreboard()
	if err != nil {
		panic(err)
	}

	for _, teamScoreInfo := range scores {
		if teamScoreInfo.Score != 100 {
			panic("invalid score")
		}
	}
}
