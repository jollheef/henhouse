/**
 * @file scoreboard_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2015
 * @brief test scoreboard
 */

package scoreboard

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/jollheef/henhouse/db"
	"github.com/jollheef/henhouse/game"
	"golang.org/x/net/websocket"
)

const dbPath = "user=postgres dbname=henhouse_test sslmode=disable"

func TestDurationToHMS(*testing.T) {

	t := time.Hour + time.Minute + time.Second
	real := "01:01:01"
	result := durationToHMS(t)
	if result != real {
		panic(fmt.Sprintf("Invalid result: %s instead %s", result, real))
	}

	t = time.Hour + time.Minute + time.Second + 100*time.Nanosecond
	real = "01:01:01"
	result = durationToHMS(t)
	if result != real {
		panic(fmt.Sprintf("Invalid result: %s instead %s", result, real))
	}

	t = 23*time.Hour + 13*time.Minute + 100*time.Nanosecond
	real = "23:13:00"
	result = durationToHMS(t)
	if result != real {
		panic(fmt.Sprintf("Invalid result: %s instead %s", result, real))
	}

	t = 0
	real = "00:00:00"
	result = durationToHMS(t)
	if result != real {
		panic(fmt.Sprintf("Invalid result: %s instead %s", result, real))
	}

	t = 15*time.Second + 100*time.Nanosecond
	real = "00:00:15"
	result = durationToHMS(t)
	if result != real {
		panic(fmt.Sprintf("Invalid result: %s instead %s", result, real))
	}

	t = 15 * time.Hour
	real = "15:00:00"
	result = durationToHMS(t)
	if result != real {
		panic(fmt.Sprintf("Invalid result: %s instead %s", result, real))
	}

	t = 15 * time.Minute
	real = "00:15:00"
	result = durationToHMS(t)
	if result != real {
		panic(fmt.Sprintf("Invalid result: %s instead %s", result, real))
	}
}

func testMatch(pattern, s string) {
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		panic(err)
	}
	if !matched {
		panic(s)
	}
}

func testNotMatch(pattern, s string) {
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		panic(err)
	}
	if matched {
		panic(s)
	}
}

func TestGetInfo(*testing.T) {

	database, err := db.InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer database.Close()

	startTime := time.Now().Add(time.Second)
	endTime := startTime.Add(time.Second)

	game, err := game.NewGame(database, startTime, endTime)
	if err != nil {
		panic(err)
	}

	gameShim = &game

	info := getInfo()

	testMatch(contestNotStarted, info)

	time.Sleep(time.Second)

	info = getInfo()

	testMatch(contestRunning, info)

	time.Sleep(time.Second)

	info = getInfo()

	testMatch(contestCompleted, info)
}

func matchBody(url, pattern string) {

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	testMatch(pattern, string(body))
}

func addTestData(database *sql.DB, nteams, ncategories, ntasks int,
	validFlag string) (err error) {

	for i := 0; i < nteams; i++ {

		team := db.Team{255, fmt.Sprintf("team%d", i),
			"e", "d", "l", false}

		err = db.AddTeam(database, &team)
		if err != nil {
			panic(err)
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
			}

			err = db.AddTask(database, &task)
			if err != nil {
				return
			}
		}
	}

	return
}

func checkAvailability(database *sql.DB, scoreboardURL, originURL,
	infoURL string) (err error) {

	var msg = make([]byte, 4096)

	ws, err := websocket.Dial(scoreboardURL, "", originURL)
	if err != nil {
		return
	}

	if _, err = ws.Read(msg); err != nil {
		return
	}

	testMatch("Team", string(msg))

	ws.Close()

	ws, err = websocket.Dial(infoURL, "", originURL)
	if err != nil {
		return
	}

	if _, err = ws.Read(msg); err != nil {
		return
	}

	testMatch(contestRunning, string(msg))

	ws.Close()

	return
}

func solveTasks(game *game.Game, validFlag string, start, end int) (err error) {
	var solved bool
	for i := start; i < end; i++ {
		solved, err = game.Solve(i, i, validFlag)
		if err != nil {
			return
		}
		if !solved {
			err = errors.New("solve task failed")
			return
		}
		time.Sleep(time.Second)
	}
	return
}

func checkScoreboard(database *sql.DB, game *game.Game, addr, validFlag string,
	nteams, ncategories, ntasks int) (err error) {

	originURL := "http://localhost/"

	authEnabled = false

	// Invalid id => get the fu^W scoreboard
	matchBody("http://"+addr+"/task?id=kekeke", "Scoreboard")

	matchBody("http://"+addr+"/", "DOCTYPE html")

	matchBody("http://"+addr+"/tasks.html", "DOCTYPE html")

	// Scoreboard wwwPath is "" => must be not found
	matchBody("http://"+addr+"/news.html", "not found")

	infoURL := "ws://" + addr + "/info"

	ws, err := websocket.Dial(infoURL, "", originURL)
	if err != nil {
		panic(err)
	}

	var msg = make([]byte, 4096)
	if _, err = ws.Read(msg); err != nil {
		panic(err)
	}

	testMatch(contestRunning, string(msg))

	ws.Close()

	time.Sleep(3 * time.Second)

	scoreboardURL := "ws://" + addr + "/scoreboard"

	ws, err = websocket.Dial(scoreboardURL, "", originURL)
	if err != nil {
		panic(err)
	}

	if _, err = ws.Read(msg); err != nil {
		panic(err)
	}

	testMatch("Team", string(msg))

	for i := 0; i < nteams; i++ {
		testMatch(fmt.Sprintf("team%d", i), string(msg))
	}

	solved, err := game.Solve(1, 1, validFlag)
	if err != nil {
		panic(err)
	}
	if !solved {
		panic("solve task failed")
	}

	time.Sleep(time.Second)

	if _, err = ws.Read(msg); err != nil {
		panic(err)
	}

	testMatch("Team", string(msg))

	for i := 1; i < nteams; i++ {
		testMatch(fmt.Sprintf("<td>team%d</td><td>0</td>", i),
			string(msg))
	}

	testMatch("<td>1</td><td>team0</td><td>500</td>", string(msg))

	ws.Close()

	err = solveTasks(game, validFlag, 10, 15)
	if err != nil {
		panic(err)
	}

	// tasks page
	tasksInfoURL := "ws://" + addr + "/tasks"

	ws, err = websocket.Dial(tasksInfoURL, "", originURL)
	if err != nil {
		panic(err)
	}

	if _, err = ws.Read(msg); err != nil {
		panic(err)
	}

	testMatch("category", string(msg))

	ws.Close()

	tasksURL := "ws://" + addr + "/info"

	ws, err = websocket.Dial(tasksURL, "", originURL)
	if err != nil {
		panic(err)
	}

	if _, err = ws.Read(msg); err != nil {
		panic(err)
	}

	testMatch(contestRunning, string(msg))

	ws.Close()

	// Check availability after close database

	database.Close()

	time.Sleep(time.Second * 2)

	err = checkAvailability(database, scoreboardURL, originURL, infoURL)
	if err != nil {
		panic(err)
	}

	return
}

func TestScoreboard(*testing.T) {

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

	start := time.Now()
	end := start.Add(time.Hour)

	game, err := game.NewGame(database, start, end)
	if err != nil {
		panic(err)
	}

	addr := "localhost:8080"

	err = game.Run()
	if err != nil {
		panic(err)
	}

	go func() {
		err = Scoreboard(database, &game, "", addr)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second) // wait for start listening

	err = checkScoreboard(database, &game,
		addr, validFlag, nteams, ncategories, ntasks)
	if err != nil {
		panic(err)
	}
}
