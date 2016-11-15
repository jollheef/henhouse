/**
 * @file html.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
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

func taskToHTML(teamID int, task game.TaskInfo,
	ru bool) (html string) {

	buttonClass := "closed"

	if len(task.SolvedBy) == 0 && task.Opened {
		buttonClass = "opened"
	} else if taskSolvedBy(task, teamID) {
		buttonClass = "success"
	}

	if task.Opened {
		html = fmt.Sprintf(`<a href="/task?id=%d" `+
			`class="task_block task_block-%s">`,
			task.ID, buttonClass)
	} else {
		html = fmt.Sprintf(`<a class="task_block task_block-%s">`, buttonClass)
	}

	var name string
	if ru {
		name = task.Name
	} else {
		name = task.NameEn
	}

	html += fmt.Sprintf(`
          <div class="task_block-header">
	    <span class="task_block-name">%s</span>
	  </div>
	  <div class="task_block-body">%d</div>
	  <div class="task_block-footer">
	    <span class="task_block-tags">%s</span>
	  </div>
	</a>`, name, task.Price, task.Tags)

	return
}

func categoryToHTML(teamID int, category game.CategoryInfo,
	ru bool) (html string) {

	html = `<div class="col-xs-3">`

	for _, task := range category.TasksInfo {
		html += taskToHTML(teamID, task, ru)
	}

	html += `</div>`

	return
}
