/**
 * @file scoreboard_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date November, 2015
 * @brief test scoreboard
 */

package scoreboard

import (
	"fmt"
	"github.com/jollheef/henhouse/db"
	"github.com/jollheef/henhouse/game"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"
	"time"
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
			"e", "d", "l", "p"}

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

	originURL := "http://localhost/"

	cats, err := game.Tasks()
	if err != nil {
		panic(err)
	}

	authEnabled = false

	matchBody("http://"+addr+"/task?id=1", cats[0].TasksInfo[0].Desc)

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

	var msg = make([]byte, 1024)
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

	testMatch("<tr><td>0</td><td>team0</td><td>500</td><tr>", string(msg))

	ws.Close()

	for i := 10; i < 15; i++ {
		solved, err = game.Solve(i, i, validFlag)
		if err != nil {
			panic(err)
		}
		if !solved {
			panic("solve task failed")
		}
		time.Sleep(time.Second)
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

	// Check availablity after close database

	database.Close()

	time.Sleep(time.Second * 2)

	ws, err = websocket.Dial(scoreboardURL, "", originURL)
	if err != nil {
		panic(err)
	}

	if _, err = ws.Read(msg); err != nil {
		panic(err)
	}

	testMatch("Team", string(msg))

	ws.Close()

	ws, err = websocket.Dial(infoURL, "", originURL)
	if err != nil {
		panic(err)
	}

	if _, err = ws.Read(msg); err != nil {
		panic(err)
	}

	testMatch(contestRunning, string(msg))

	ws.Close()
}
