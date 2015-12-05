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

func taskToHTML(task game.TaskInfo) (html string) {

	html = `<p><button class="btn btn-primary btn-lg" `

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

func categoryToHTML(category game.CategoryInfo) (html string) {

	html = fmt.Sprintf(`<div class="col-xs-3"> <h1>%s</h1> `,
		category.Name)

	for _, task := range category.TasksInfo {
		html += taskToHTML(task)
	}

	html += `</div>`

	return
}
