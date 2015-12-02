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

	head := `<p><a class="btn btn-primary btn-lg"`

	if task.Opened {
		html = fmt.Sprintf(head+`title="%s">%d. %s</a></p>`,
			task.Name, task.Price, task.Name)
	} else {
		html = fmt.Sprintf(head+
			`disabled="disabled" title="Task is closed">%d. %s</a></p>`,
			task.Price, task.Name)
	}

	return
}

func categoryToHTML(category game.CategoryInfo) (html string) {

	html = fmt.Sprintf(`<div class="col-xs-3"> <h1>%s</h1>`,
		category.Name)

	for _, task := range category.TasksInfo {
		html += taskToHTML(task)
	}

	html += `</div>`

	return
}
