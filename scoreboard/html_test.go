/**
 * @file html_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date December, 2015
 * @brief test html helpers
 */

package scoreboard

import (
	"github.com/jollheef/henhouse/game"
	"testing"
)

func TestTaskToHTML(*testing.T) {
	html := taskToHTML(1, game.TaskInfo{})
	testMatch("Task is closed", html)

	html = taskToHTML(1, game.TaskInfo{Opened: true})
	testNotMatch("Task is closed", html)
}

func TestCategoryToHTML(*testing.T) {

	cat := game.CategoryInfo{}

	cat.TasksInfo = append(cat.TasksInfo, game.TaskInfo{})

	html := categoryToHTML(0, cat)

	testMatch("Task is closed", html)
}
