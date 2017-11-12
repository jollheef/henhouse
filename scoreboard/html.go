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

	if taskSolvedBy(task, teamID) {
		buttonClass = "success"
	} else if len(task.SolvedBy) > 0 {
		buttonClass = "solved"
	} else if task.Opened {
		buttonClass = "opened"
	}

	if task.Opened {
		html = fmt.Sprintf(`
			<div class="jctf-task mdl-cell mdl-cell--3-col">
				<a href="/task?id=%d" `+
				`class="jctf-task-block %s">`,
			task.ID, buttonClass)
	} else {
		html = fmt.Sprintf(`
			<div class="jctf-task mdl-cell mdl-cell--3-col">
				<a class="jctf-task-block %s">`, 
			buttonClass)
	}

	var name string
	if ru {
		name = task.Name
	} else {
		name = task.NameEn
	}

	if task.Opened {
		html += fmt.Sprintf(`
			<div class="jctf-task-block__header">
				<span class="jctf-task-block__name">%s</span>
			</div>
			<div class="jctf-task-block__body">%d</div>
			<div class="jctf-task-block__footer">
				<span class="jctf-task-block__tags">%s</span>
			</div>
		</a></div>`, name, task.Price, task.Tags)
	} else {
		html += fmt.Sprintf(`
			<div class="jctf-task-block__header">
				<span class="jctf-task-block__name"></span>
			</div>
			<div class="jctf-task-block__body">
				<img class="jctf-image-closed-task" src="images/closed_task.png">
			</div>
			<div class="jctf-task-block__footer">
				<span class="jctf-task-block__tags"></span>
			</div>
		</a></div>`)
	}

	return
}

func categoryToHTML(teamID int, category game.CategoryInfo,
	ru bool) (html string) {

	html = `<div class="jctf-content__taskline mdl-cell mdl-cell--12-col mdl-grid">`

	for _, task := range category.TasksInfo {
		html += taskToHTML(teamID, task, ru)
	}

	html += `</div>`

	return
}
