/**
 * @file html_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date December, 2015
 * @brief test html helpers
 */

package scoreboard

import (
	"testing"

	"github.com/jollheef/henhouse/game"
)

func TestTaskToHTML(*testing.T) {
	html := taskToHTML(1, game.TaskInfo{}, true)
	testMatch("closed", html)

	html = taskToHTML(1, game.TaskInfo{Opened: true}, true)
	testNotMatch("closed", html)
}

func TestCategoryToHTML(*testing.T) {

	cat := game.CategoryInfo{}

	cat.TasksInfo = append(cat.TasksInfo, game.TaskInfo{})

	html := categoryToHTML(0, cat, true)

	testMatch("closed", html)
}
