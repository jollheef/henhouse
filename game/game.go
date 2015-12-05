/**
 * @file game.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date November, 2015
 * @brief game api
 *
 * Contain functions for calculate score, check flag etc.
 */

package game

import (
	"database/sql"
	"github.com/jollheef/henhouse/db"
	"log"
	"sort"
	"sync"
	"time"
)

// Game struct
type Game struct {
	db              *sql.DB
	Start           time.Time
	End             time.Time
	OpenTimeout     time.Duration // after solve task
	AutoOpen        bool
	AutoOpenTimeout time.Duration // if task does not solved
	teams           []db.Team
	tasks           []db.Task
	categories      []db.Category
	scoreboardLock  sync.Mutex
}

// TaskInfo provide information about task
type TaskInfo struct {
	ID       int
	Name     string
	Desc     string
	Author   string
	Price    int
	Opened   bool
	SolvedBy []int
}

// CategoryInfo provide information about categories and tasks
type CategoryInfo struct {
	Name      string
	TasksInfo []TaskInfo
}

// TeamScoreInfo provide information about team score
type TeamScoreInfo struct {
	ID    int
	Name  string
	Score int
}

type byScore []TeamScoreInfo

func (tr byScore) Len() int           { return len(tr) }
func (tr byScore) Swap(i, j int)      { tr[i], tr[j] = tr[j], tr[i] }
func (tr byScore) Less(i, j int) bool { return tr[i].Score > tr[j].Score }

// NewGame create new game
func NewGame(database *sql.DB, start, end time.Time) (g Game, err error) {

	g.db = database
	g.Start = start
	g.End = end

	g.teams, err = db.GetTeams(g.db)
	if err != nil {
		return
	}

	g.tasks, err = db.GetTasks(g.db)
	if err != nil {
		return
	}

	g.categories, err = db.GetCategories(g.db)
	if err != nil {
		return
	}

	g.RecalcScoreboard()

	return
}

// Run open first level tasks and start auto open routine
func (g Game) Run() (err error) {

	for time.Now().Before(g.Start) {
		time.Sleep(time.Second)
	}

	for i, task := range g.tasks {
		if task.Level == 1 && !task.Opened {
			err = db.SetOpened(g.db, task.ID, true)
			if err != nil {
				return
			}

			g.tasks[i].Opened = true

			if g.AutoOpen {
				go g.autoOpen(task)
			}
		}
	}

	return
}

func (g Game) autoOpen(task db.Task) {
	time.Sleep(g.AutoOpenTimeout)
	err := g.OpenNextTask(task)
	if err != nil {
		log.Println("Auto open next task fail:", err)
	}
}

func taskPrice(database *sql.DB, taskID int) (price int, err error) {

	count, err := db.GetSolvedCount(database, taskID)

	fprice := float64(count) / 20.0

	if fprice <= 0.1 {
		price = 500
	} else if fprice <= 0.15 {
		price = 400
	} else if fprice <= 0.3 {
		price = 300
	} else if fprice <= 0.5 {
		price = 200
	} else {
		price = 100
	}

	return
}

// Tasks returns categories with tasks
func (g Game) Tasks() (cats []CategoryInfo, err error) {

	for _, category := range g.categories {

		cat := CategoryInfo{Name: category.Name}

		for _, task := range g.tasks {

			if task.CategoryID == category.ID {

				var price int
				price, err = taskPrice(g.db, task.ID)
				if err != nil {
					return
				}

				var solvedBy []int
				solvedBy, err = db.GetSolvedBy(g.db, task.ID)
				if err != nil {
					return
				}

				tInfo := TaskInfo{
					ID:       task.ID,
					Name:     task.Name,
					Desc:     task.Desc,
					Price:    price,
					Opened:   task.Opened,
					SolvedBy: solvedBy,
					Author:   task.Author,
				}

				cat.TasksInfo = append(cat.TasksInfo, tInfo)
			}
		}

		cats = append(cats, cat)
	}

	return
}

// Scoreboard returns sorted scoreboard
func (g Game) Scoreboard() (scores []TeamScoreInfo, err error) {

	g.scoreboardLock.Lock()

	for _, team := range g.teams {

		if team.Test {
			continue
		}

		var s db.Score
		s, err = db.GetLastScore(g.db, team.ID)
		if err != nil {
			return
		}

		scores = append(scores, TeamScoreInfo{team.ID, team.Name, s.Score})
	}

	sort.Sort(byScore(scores))

	g.scoreboardLock.Unlock()

	return
}

// RecalcScoreboard update scoreboard
func (g Game) RecalcScoreboard() (err error) {

	g.scoreboardLock.Lock()

	for _, team := range g.teams {

		if team.Test {
			continue
		}

		score := 0

		for _, task := range g.tasks {

			var price int
			price, err = taskPrice(g.db, task.ID)
			if err != nil {
				return
			}

			var solved bool
			solved, err = db.IsSolved(g.db, team.ID, task.ID)
			if err != nil {
				return
			}

			if solved {
				score += price
			}
		}

		err = db.AddScore(g.db, &db.Score{TeamID: team.ID, Score: score})
		if err != nil {
			return
		}
	}

	g.scoreboardLock.Unlock()

	return
}

// OpenNextTask open next task by level
func (g Game) OpenNextTask(t db.Task) (err error) {

	time.Sleep(g.OpenTimeout)

	for i, task := range g.tasks {
		// If same category and next level
		if t.CategoryID == task.CategoryID && t.Level+1 == task.Level {
			// If not already opened
			if !task.Opened {
				// Open it!
				err = db.SetOpened(g.db, task.ID, true)
				if err != nil {
					return
				}

				g.tasks[i].Opened = true

				if g.AutoOpen {
					go g.autoOpen(task)
				}
			}
		}
	}

	return
}

func (g Game) isTestTeam(teamID int) bool {
	for _, team := range g.teams {
		if team.ID == teamID {
			return team.Test
		}
	}
	return false
}

// Solve check flag for task and recalc scoreboard if flag correct
func (g Game) Solve(teamID, taskID int, flag string) (solved bool, err error) {

	for _, task := range g.tasks {
		if task.ID == taskID {

			if task.Flag == flag { // fix to regex

				solved = true

				if g.isTestTeam(teamID) {
					return
				}

				now := time.Now()

				if now.After(g.Start) && now.Before(g.End) {
					err = db.AddFlag(g.db, &db.Flag{
						TeamID: teamID,
						TaskID: taskID,
						Flag:   flag,
						Solved: solved,
					})
					if err != nil {
						return
					}

					go g.OpenNextTask(task)
				}
			}

			break
		}
	}

	return
}
