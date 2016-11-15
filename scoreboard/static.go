/**
 * @file static.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date December, 2015
 * @brief non-dynamic html results
 *
 * Generate static html page
 */

package scoreboard

import (
	"fmt"
	"log"
	"net/http"
)

func staticScoreboard(w http.ResponseWriter, r *http.Request) {

	teamID := getTeamID(r)

	tmpl, err := getTmpl("scoreboard")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Fprintf(w, tmpl, getInfo(), scoreboardHTML(teamID))
}

func staticTasks(w http.ResponseWriter, r *http.Request) {

	teamID := getTeamID(r)

	tmpl, err := getTmpl("tasks")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Fprintf(w, tmpl, getInfo(), tasksHTML(teamID))
}
