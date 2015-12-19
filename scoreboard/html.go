/**
 * @file html.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date December, 2015
 * @brief html helpers
 */

package scoreboard

import (
	"fmt"
	"github.com/jollheef/henhouse/game"
)

func taskSolvedBy(task game.TaskInfo, teamID int) bool {
	for _, t := range task.SolvedBy {
		if t == teamID {
			return true
		}
	}
	return false
}

func taskToHTML(teamID int, task game.TaskInfo) (html string) {

	buttonClass := "primary"

	if len(task.SolvedBy) == 0 && task.Opened {
		buttonClass = "warning"
	} else if taskSolvedBy(task, teamID) {
		buttonClass = "default"
	}

	html = fmt.Sprintf(`<p><button class="btn btn-%s"`, buttonClass)

	if task.Opened {
		html += fmt.Sprintf(`title="%s" onclick="window.location=`+
			`'task?id=%d';">%d. %s `,
			task.Name, task.ID, task.Price, task.Name)
	} else {
		html += fmt.Sprintf(`disabled="disabled" `+
			`title="Task is closed">%d. %s`,
			task.Price, task.Name)
	}

	html += "</button></p>"

	return
}

func categoryToHTML(teamID int, category game.CategoryInfo) (html string) {

	html = fmt.Sprintf(`<div class="col-xs-3"> <h1>%s</h1> `,
		category.Name)

	for _, task := range category.TasksInfo {
		html += taskToHTML(teamID, task)
	}

	html += `</div>`

	return
}
