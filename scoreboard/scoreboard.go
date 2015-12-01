/**
 * @file scoreboard.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date November, 2015
 * @brief web scoreboard
 *
 * Contain web ui and several helpers for convert results to table
 */

package scoreboard

import (
	"fmt"
	"github.com/jollheef/henhouse/game"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"time"
)

const (
	contestStateNotAvailable = "state n/a"
	contestNotStarted        = "not started"
	contestRunning           = "running"
	contestCompleted         = "completed"
)

var (
	startTime             time.Time
	endTime               time.Time
	currentResult         string
	currentTasks          string
	contestStatus         string
	lastScoreboardUpdated string
	lastTasksUpdated      string
)

func durationToHMS(d time.Duration) string {

	sec := int(d.Seconds())

	var h, m, s int

	h = sec / 60 / 60
	m = (sec / 60) % 60
	s = sec % 60

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func getInfo(lastUpdated string) string {

	var left time.Duration
	var btnType string

	now := time.Now()

	if now.Before(startTime) {

		contestStatus = contestNotStarted
		left = startTime.Sub(now)
		btnType = "warning"

	} else if now.Before(endTime) {

		contestStatus = contestRunning
		left = endTime.Sub(now)
		btnType = "success"

	} else {
		contestStatus = contestCompleted
		left = 0
		btnType = "primary"
	}

	info := fmt.Sprintf(`<span class="btn btn-%s">contest %s</span>`,
		btnType, contestStatus)

	if left != 0 {
		info += fmt.Sprintf(`<span class="btn btn-info">Left %s</span>`+
			`<span class="btn btn-info">Updated at %s</span>`,
			durationToHMS(left), lastUpdated)
	}

	return info
}

func infoHandler(ws *websocket.Conn) {

	defer ws.Close()
	for {
		_, err := fmt.Fprint(ws, getInfo(lastScoreboardUpdated))
		if err != nil {
			log.Println("Socket closed:", err)
			return
		}

		time.Sleep(time.Second)
	}
}

func tasksInfoHandler(ws *websocket.Conn) {

	defer ws.Close()
	for {
		_, err := fmt.Fprint(ws, getInfo(lastTasksUpdated))
		if err != nil {
			log.Println("Socket closed:", err)
			return
		}

		time.Sleep(time.Second)
	}
}

func scoreboardHandler(ws *websocket.Conn) {

	defer ws.Close()

	fmt.Fprint(ws, currentResult)
	sendedResult := currentResult
	lastUpdate := time.Now()

	for {
		if sendedResult != currentResult ||
			time.Now().After(lastUpdate.Add(time.Minute)) {

			sendedResult = currentResult
			lastUpdate = time.Now()

			_, err := fmt.Fprint(ws, currentResult)
			if err != nil {
				log.Println("Socket closed:", err)
				return
			}
		}

		time.Sleep(time.Second)
	}
}

func scoreboardHTMLUpdater(game *game.Game, updateTimeout time.Duration) {

	head := "<thead><th>#</th><th>Team</th><th>Score</th></thead>"

	for {
		result := head

		result += "<tbody>"

		scores, err := game.Scoreboard()
		if err != nil {
			log.Println("Get scoreboard fail:", err)
			time.Sleep(updateTimeout)
			continue
		}

		for n, teamScore := range scores {
			result += fmt.Sprintf(
				"<tr><td>%d</td><td>%s</td><td>%d</td><tr>",
				n, teamScore.Name, teamScore.Score)

		}

		result += "</tbody>"

		currentResult = result

		now := time.Now()
		lastScoreboardUpdated = fmt.Sprintf("%02d:%02d:%02d", now.Hour(),
			now.Minute(), now.Second())

		time.Sleep(updateTimeout)
	}
}

func scoreboardUpdater(game *game.Game, updateTimeout time.Duration) {

	for {
		time.Sleep(updateTimeout)

		err := game.RecalcScoreboard()
		if err != nil {
			log.Println("Recalc scoreboard fail:", err)
		}
	}
}

func tasksHTMLUpdater(game *game.Game, updateTimeout time.Duration) {

	for {
		cats, err := game.Tasks()
		if err != nil {
			log.Println("Get tasks fail:", err)
		}

		var result string

		for _, cat := range cats {
			result += categoryToHTML(cat)
		}

		currentTasks = result

		now := time.Now()
		lastTasksUpdated = fmt.Sprintf("%02d:%02d:%02d", now.Hour(),
			now.Minute(), now.Second())

		time.Sleep(updateTimeout)
	}
}

func tasksHandler(ws *websocket.Conn) {

	defer ws.Close()

	fmt.Fprint(ws, currentTasks)
	sendedTasks := currentTasks
	lastUpdate := time.Now()

	for {
		if sendedTasks != currentTasks ||
			time.Now().After(lastUpdate.Add(time.Minute)) {

			sendedTasks = currentTasks
			lastUpdate = time.Now()

			_, err := fmt.Fprint(ws, currentTasks)
			if err != nil {
				log.Println("Socket closed:", err)
				return
			}
		}

		time.Sleep(time.Second)
	}
}

// Scoreboard implements web scoreboard
func Scoreboard(game *game.Game, wwwPath, addr string) (err error) {

	contestStatus = contestStateNotAvailable
	startTime = game.Start
	endTime = game.End

	go scoreboardHTMLUpdater(game, time.Second)
	go tasksHTMLUpdater(game, time.Second)

	go scoreboardUpdater(game, time.Second)

	http.Handle("/scoreboard", websocket.Handler(scoreboardHandler))
	http.Handle("/scoreboard-info", websocket.Handler(infoHandler))
	http.Handle("/tasks", websocket.Handler(tasksHandler))
	http.Handle("/tasks-info", websocket.Handler(tasksInfoHandler))
	http.Handle("/", http.FileServer(http.Dir(wwwPath)))

	log.Println("Launching scoreboard at", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		return
	}

	return
}
