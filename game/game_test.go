/**
 * @file game_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2015
 * @brief test game package
 */

package game

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jollheef/henhouse/db"
)

const dbPath string = "user=postgres dbname=henhouse_test sslmode=disable"

// TestNewGame test new game
func TestNewGame(*testing.T) {

	database, err := db.InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer database.Close()

	teams, err := db.GetTeams(database)
	if err != nil {
		return
	}

	_, err = NewGame(database, time.Now(),
		time.Now().Add(time.Hour), float64(len(teams)))
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

	teams, err := db.GetTeams(database)
	if err != nil {
		return
	}

	_, err = NewGame(database, time.Now(),
		time.Now().Add(time.Hour), float64(len(teams)))
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

	teams, err := db.GetTeams(database)
	if err != nil {
		return
	}

	game, err := NewGame(database, time.Now(),
		time.Now().Add(time.Hour), float64(len(teams)))
	if err != nil {
		panic(err)
	}

	cats, err := game.Tasks()
	if err != nil {
		panic(err)
	}

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

func addTestData(database *sql.DB, nteams, ncategories, ntasks int,
	validFlag string) (err error) {

	for i := 0; i < nteams; i++ {

		team := db.Team{255, fmt.Sprintf("team%d", i),
			"e", "d", "l", false}

		err = db.AddTeam(database, &team)
		if err != nil {
			return (err)
		}
	}

	for i := 0; i < ncategories; i++ {

		category := db.Category{Name: fmt.Sprintf("category%d", i)}

		err = db.AddCategory(database, &category)
		if err != nil {
			return
		}

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
				Level:         i,
			}

			err = db.AddTask(database, &task)
			if err != nil {
				return
			}
		}
	}

	return
}

func checkSumScores(scores []TeamScoreInfo)(sumScores int){
	for _, teamScoreInfo := range scores {
		if teamScoreInfo.Score != 0 {
			sumScores = 1
			return
		}
	}
	sumScores = 0
	return
}
func TestScoreboard(*testing.T) {

	var sumCheck int

	database, err := db.InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer database.Close()

	validFlag := "testflag"

	nteams := 20
	ncategories := 5
	ntasks := 5

	err = addTestData(database, nteams, ncategories, ntasks, validFlag)
	if err != nil {
		panic(err)
	}

	teams, err := db.GetTeams(database)
	if err != nil {
		return
	}

	game, err := NewGame(database, time.Now().Add(time.Second),
		time.Now().Add(time.Hour), float64(len(teams)))
	if err != nil {
		panic(err)
	}

	err = game.Run()
	if err != nil {
		panic(err)
	}

	scores, err := game.Scoreboard()
	if err != nil {
		panic(err)
	}

	sumCheck = checkSumScores(scores)

	if sumCheck == 1 {
		panic("score at game start not zero")
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

func fillTestDB(database *sql.DB, validFlag string) (err error) {

	nteams := 4

	for i := 0; i < nteams; i++ {

		team := db.Team{255, fmt.Sprintf("team%d", i),
			"e", "d", "l", false}

		err = db.AddTeam(database, &team)
		if err != nil {
			panic(err)
		}
	}

	ncategories := 4

	for i := 0; i < ncategories; i++ {

		category := db.Category{Name: fmt.Sprintf("category%d", i)}

		err = db.AddCategory(database, &category)
		if err != nil {
			panic(err)
		}

		ntasks := 4

		for i := 0; i < ntasks; i++ {

			task := db.Task{
				Name:          fmt.Sprintf("task%d", i),
				Flag:          validFlag,
				CategoryID:    category.ID,
				Price:         500,
				Level:         i + 1,
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

	return
}

func testSolveTask(database *sql.DB, game *Game, teamID, taskID int,
	validFlag string) (err error) {

	solved, err := game.Solve(teamID, taskID, validFlag)
	if err != nil {
		return
	}
	if !solved {
		err = errors.New("solve task failed")
		return
	}

	solved, err = db.IsSolved(database, teamID, taskID)
	if !solved {
		err = errors.New("is solved task check failed: unsolved")
		return
	}

	return
}

func TestSolve(*testing.T) {

	database, err := db.InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer database.Close()

	validFlag := "testflag"

	err = fillTestDB(database, validFlag)
	if err != nil {
		panic(err)
	}

	start := time.Now().Add(time.Second)
	end := start.Add(time.Second)

	teams, err := db.GetTeams(database)
	if err != nil {
		return
	}

	game, err := NewGame(database, start, end, float64(len(teams)))
	if err != nil {
		panic(err)
	}

	// Try to solve task before game start
	err = testSolveTask(database, &game, 1, 1, validFlag)
	if err == nil {
		panic("task solved before game start")
	}
	time.Sleep(time.Second)

	// Try to solve task after game start
	err = testSolveTask(database, &game, 2, 2, validFlag)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second)

	// Try to solve task after game end
	err = testSolveTask(database, &game, 3, 3, validFlag)
	if err == nil {
		panic("task solved after game end")
	}
	time.Sleep(time.Second)
}

func TestFirstOpen(*testing.T) {

	database, err := db.InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer database.Close()

	validFlag := "testflag"

	err = fillTestDB(database, validFlag)
	if err != nil {
		panic(err)
	}

	start := time.Now().Add(time.Second)
	end := start.Add(time.Second)

	teams, err := db.GetTeams(database)
	if err != nil {
		return
	}

	game, err := NewGame(database, start, end, float64(len(teams)))
	if err != nil {
		panic(err)
	}

	game.Run()

	cats, err := game.Tasks()
	if err != nil {
		panic(err)
	}

	for _, c := range cats {
		if !c.TasksInfo[0].Opened {
			log.Fatalln("Err: first task not opened")
		}

		for i := 1; i < len(c.TasksInfo); i++ {
			if c.TasksInfo[i].Opened {
				log.Fatalln("Err: not first task is opened")
			}
		}
	}
}

func initGame(teamID, taskID int, flag string) (database *sql.DB, game Game) {
	database, err := db.InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	err = fillTestDB(database, flag)
	if err != nil {
		panic(err)
	}

	start := time.Now().Add(time.Second)
	end := start.Add(time.Second)

	teams, err := db.GetTeams(database)
	if err != nil {
		return
	}

	game, err = NewGame(database, start, end, float64(len(teams)))
	if err != nil {
		panic(err)
	}

	game.SetTeamsBase(5)

	return
}

func TestTaskPriceDefaultValues(*testing.T) {
	teamID := 1
	taskID := 1
	validFlag := "testflag"

	database, game := initGame(teamID, taskID, validFlag)
	defer database.Close()

	game.Run()

	game.Solve(teamID, taskID, validFlag)
	game.Solve(teamID, taskID, validFlag)
	game.Solve(teamID, taskID, validFlag)

	price, err := game.taskPrice(database, taskID)
	if err != nil {
		panic(err)
	}

	if price != 300 {
		panic("price mismatch")
	}
}

func TestTaskPriceFixedValues(*testing.T) {
	teamID := 1
	taskID := 1
	validFlag := "testflag"

	database, game := initGame(teamID, taskID, validFlag)
	defer database.Close()

	game.SetTaskPrice(100, 100, 100, 100)

	game.Run()

	game.Solve(teamID, taskID, validFlag)
	game.Solve(teamID, taskID, validFlag)
	game.Solve(teamID, taskID, validFlag)

	price, err := game.taskPrice(database, taskID)
	if err != nil {
		panic(err)
	}

	if price != 500 { // max
		panic("price mismatch")
	}
}

func TestAutoOpenTimeoutDisabled(*testing.T) {
	database, game := initGame(0, 0, "")
	defer database.Close()

	game.AutoOpen = false
	game.AutoOpenTimeout = time.Nanosecond

	game.Run()

	time.Sleep(time.Second * 2)

	tasks, err := db.GetTasks(database)
	if err != nil {
		return
	}

	for _, t := range tasks {
		if t.Level > 1 && t.Opened {
			panic("auto open disabled but level >1 opened")
		} else if t.Level == 1 && !t.Opened {
			panic("level one not opened")
		}
	}
}

func TestAutoOpenTimeoutEnabled(*testing.T) {
	database, game := initGame(0, 0, "")
	defer database.Close()

	game.AutoOpen = true
	game.AutoOpenTimeout = time.Nanosecond

	game.Run()

	time.Sleep(time.Second * 2)

	tasks, err := db.GetTasks(database)
	if err != nil {
		return
	}

	for _, t := range tasks {
		if (t.Level == 1 || t.Level == 2) && !t.Opened {
			panic("levels 1, 2 not opened")
		}
	}
}

func TestCalcTeamsBase(*testing.T) {
	database, _ := initGame(0, 0, "")
	defer database.Close()

	n, err := CalcTeamsBase(database)
	if err != nil {
		panic(err)
	}

	if n != 21 {
		panic("Default abount of teams not equal to 21")
	}
}
