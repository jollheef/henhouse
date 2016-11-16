/**
 * @file scoreboard.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2015
 * @brief web scoreboard
 *
 * Contain web ui and several helpers for convert results to table
 */

package scoreboard

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fiam/gounidecode/unidecode"
	"github.com/jollheef/henhouse/game"
	"golang.org/x/net/websocket"
)

const (
	contestStateNotAvailable = "state n/a"
	contestNotStarted        = "not started"
	contestRunning           = "running"
	contestCompleted         = "completed"
)

var (
	gameShim      *game.Game
	contestStatus string
)

var (
	// InfoTimeout timeout between update info through websocket
	InfoTimeout = time.Second
	// ScoreboardTimeout timeout between update scoreboard through websocket
	ScoreboardTimeout = time.Second
	// TasksTimeout timeout between update tasks through websocket
	TasksTimeout = time.Second
	// FlagTimeout timeout between send flags
	FlagTimeout = time.Second
	// ScoreboardRecalcTimeout timeout between update scoreboard
	ScoreboardRecalcTimeout = time.Second
)

func durationToHMS(d time.Duration) string {

	sec := int(d.Seconds())

	var h, m, s int

	h = sec / 60 / 60
	m = (sec / 60) % 60
	s = sec % 60

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func getInfo() string {

	var left time.Duration
	var btnType string

	now := time.Now()

	if now.Before(gameShim.Start) {

		contestStatus = contestNotStarted
		left = gameShim.Start.Sub(now)
		btnType = "stop"

	} else if now.Before(gameShim.End) {

		contestStatus = contestRunning
		left = gameShim.End.Sub(now)
		btnType = "run"

	} else {
		contestStatus = contestCompleted
		left = 0
		btnType = "stop"
	}

	info := fmt.Sprintf(`<span id="game_status-%s">contest %s</span>`,
		btnType, contestStatus)

	if left != 0 {
		info += fmt.Sprintf(`<span id="timer">Left %s</span>`,
			durationToHMS(left))
	}

	return info
}

func infoHandler(ws *websocket.Conn) {

	defer ws.Close()
	for {
		_, err := fmt.Fprint(ws, l10n(ws.Request(), getInfo()))
		if err != nil {
			//log.Println("Socket closed:", err)
			return
		}

		time.Sleep(InfoTimeout)
	}
}

func scoreboardHandler(ws *websocket.Conn) {

	defer ws.Close()

	teamID := getTeamID(ws.Request())

	currentResult := scoreboardHTML(teamID)

	fmt.Fprint(ws, l10n(ws.Request(), currentResult))

	sendedResult := currentResult

	lastUpdate := time.Now()

	for {
		currentResult := scoreboardHTML(teamID)

		if sendedResult != currentResult ||
			time.Now().After(lastUpdate.Add(time.Minute)) {

			sendedResult = currentResult
			lastUpdate = time.Now()

			_, err := fmt.Fprint(ws, l10n(ws.Request(), currentResult))
			if err != nil {
				//log.Println("Socket closed:", err)
				return
			}
		}

		time.Sleep(ScoreboardTimeout)
	}
}

func scoreboardHTML(teamID int) (result string) {

	result = "<thead>" +
		"<th>#</th>" + "<th>Team</th>" +
		"<th>Score</th>" +
		"</thead>"

	result += "<tbody>"

	scores, err := gameShim.Scoreboard()
	if err != nil {
		log.Println("Get scoreboard fail:", err)
		return
	}

	for n, teamScore := range scores {
		if teamScore.ID == teamID {
			result += `<tr class="self-team">`
		} else {
			result += `<tr>`
		}

		result += fmt.Sprintf(
			`<td class="team_index">%d</td>`+
				`<td class="team_name">%s</td>`+
				`<td class="team_score">%d</td></tr>`,
			n+1, teamScore.Name, teamScore.Score)

	}

	result += "</tbody>"

	return
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

func tasksHTML(teamID int, ru bool) (result string) {

	cats, err := gameShim.Tasks()
	if err != nil {
		log.Println("Get tasks fail:", err)
	}

	for _, cat := range cats {
		result += categoryToHTML(teamID, cat, ru)
	}

	return
}

func tasksHandler(ws *websocket.Conn) {

	defer ws.Close()

	teamID := getTeamID(ws.Request())

	currentTasks := tasksHTML(teamID, isAcceptRussian(ws.Request()))

	fmt.Fprint(ws, l10n(ws.Request(), currentTasks))

	sendedTasks := currentTasks

	lastUpdate := time.Now()

	for {
		currentTasks := tasksHTML(teamID, isAcceptRussian(ws.Request()))

		if sendedTasks != currentTasks ||
			time.Now().After(lastUpdate.Add(time.Minute)) {

			sendedTasks = currentTasks
			lastUpdate = time.Now()

			_, err := fmt.Fprint(ws, l10n(ws.Request(), currentTasks))
			if err != nil {
				//log.Println("Socket closed:", err)
				return
			}
		}

		time.Sleep(TasksTimeout)
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println("Atoi fail:", err)
		http.Redirect(w, r, "/", 307)
		return
	}

	cats, err := gameShim.Tasks()
	if err != nil {
		log.Println("Get tasks fail:", err)
		http.Redirect(w, r, "/", 307)
		return
	}

	task := game.TaskInfo{ID: id, Opened: false}

	for _, c := range cats {
		for _, t := range c.TasksInfo {
			if t.ID == id {
				task = t
				break
			}
		}
	}

	if !task.Opened {
		// Try to see closed task -> gtfo
		http.Redirect(w, r, "/", 307)
		return
	}

	teamID := getTeamID(r)

	flagSubmitFormat := `<br>` +
		`<form class="input-group" action="/flag?id=%d" method="post">` +
		`<input class="form-control float-left" name="flag" value="" placeholder="Flag">` +
		`<span class="input-group-btn">` +
		`<button class="btn btn-submit">Submit</button>` +
		`</span>` +
		`</form>`

	var submitForm string
	if taskSolvedBy(task, teamID) {
		submitForm = "Already solved"
	} else {
		submitForm = fmt.Sprintf(flagSubmitFormat, task.ID)
	}

	tmpl, err := getTmpl("task")
	if err != nil {
		log.Println(err)
		return
	}

	var name, desc, author string
	if isAcceptRussian(r) {
		name = task.Name
		desc = task.Desc
		author = task.Author
	} else {
		name = task.NameEn
		desc = task.DescEn
		author = unidecode.Unidecode(task.Author)
	}

	fmt.Fprintf(w, tmpl, name, desc,
		author, submitForm)
}

func flagHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Redirect(w, r, "/", 307)
		return
	}

	taskID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println("Atoi fail:", err)
		http.Redirect(w, r, "/", 307)
		return
	}

	flag := r.FormValue("flag")

	teamID := getTeamID(r)

	solved, err := gameShim.Solve(teamID, taskID, flag)
	if err != nil {
		solved = false
	}

	var solvedMsg string
	if solved {
		solvedMsg = `<div class="flag_status solved">Solved</div>`
	} else {
		solvedMsg = `<div class="flag_status invalid">Invalid flag</div>`
	}

	log.Printf("Team ID: %d, Task ID: %d, Flag: %s, Result: %s\n",
		teamID, taskID, flag, solvedMsg)

	time.Sleep(FlagTimeout)

	tmpl, err := getTmpl("flag")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Fprintf(w, tmpl, l10n(r, solvedMsg))
}

func handleStaticFile(pattern, file string) {
	http.HandleFunc(pattern,
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, file)
		})
}

func handleStaticFileSimple(file, wwwPath string) {
	handleStaticFile(file, wwwPath+file)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := getTmpl("auth")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprint(w, l10n(r, tmpl))
}

func newsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var tmpl string
	if isAcceptRussian(r) {
		tmpl, err = getTmplWoCache("news.ru")
	} else {
		tmpl, err = getTmplWoCache("news.en")
	}
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprint(w, l10n(r, tmpl))
}

func sponsorsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var tmpl string
	if isAcceptRussian(r) {
		tmpl, err = getTmpl("sponsors.ru")
	} else {
		tmpl, err = getTmpl("sponsors.en")
	}
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprint(w, l10n(r, tmpl))
}

// Scoreboard implements web scoreboard
func Scoreboard(database *sql.DB, game *game.Game,
	wwwPath, tmpltsPath, addr string) (err error) {

	contestStatus = contestStateNotAvailable
	gameShim = game
	templatePath = tmpltsPath

	go scoreboardUpdater(game, ScoreboardRecalcTimeout)

	// Static files
	handleStaticFileSimple("/css/style.css", wwwPath)
	handleStaticFileSimple("/js/scoreboard.js", wwwPath)
	handleStaticFileSimple("/js/tasks.js", wwwPath)
	handleStaticFileSimple("/images/bg.jpg", wwwPath)
	handleStaticFileSimple("/images/favicon.ico", wwwPath)
	handleStaticFileSimple("/images/favicon.png", wwwPath)
	handleStaticFileSimple("/images/401.jpg", wwwPath)
	handleStaticFileSimple("/images/juniors_ctf_txt.png", wwwPath)

	// Get
	http.HandleFunc("/auth.html", signinHandler)
	http.HandleFunc("/outer-scoreboard", outerScoreboard)

	// Get only for authenticated
	http.Handle("/", authorized(database, http.HandlerFunc(innerScoreboard)))
	http.Handle("/index.html", authorized(database, http.HandlerFunc(innerScoreboard)))
	http.Handle("/tasks.html", authorized(database, http.HandlerFunc(staticTasks)))
	http.Handle("/logout", authorized(database, http.HandlerFunc(logoutHandler)))
	http.Handle("/news.html", authorized(database, http.HandlerFunc(newsHandler)))
	http.Handle("/sponsors.html", authorized(database, http.HandlerFunc(sponsorsHandler)))

	// Websocket
	http.Handle("/scoreboard", authorized(database, websocket.Handler(scoreboardHandler)))
	http.Handle("/info", authorized(database, websocket.Handler(infoHandler)))
	http.Handle("/tasks", authorized(database, websocket.Handler(tasksHandler)))

	// Post
	http.Handle("/task", authorized(database, http.HandlerFunc(taskHandler)))
	http.Handle("/flag", authorized(database, http.HandlerFunc(flagHandler)))

	http.HandleFunc("/auth.php", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authHandler(database, w, r)
		}))

	log.Println("Launching scoreboard at", addr)

	return http.ListenAndServe(addr, nil)
}
