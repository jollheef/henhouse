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
	"errors"
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
	tasksLock       sync.Mutex
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
	Level    int
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
	Desc  string
	Score int
}

type byScore []TeamScoreInfo

func (tr byScore) Len() int           { return len(tr) }
func (tr byScore) Swap(i, j int)      { tr[i], tr[j] = tr[j], tr[i] }
func (tr byScore) Less(i, j int) bool { return tr[i].Score > tr[j].Score }

type byLevel []TaskInfo

func (ti byLevel) Len() int           { return len(ti) }
func (ti byLevel) Swap(i, j int)      { ti[i], ti[j] = ti[j], ti[i] }
func (ti byLevel) Less(i, j int) bool { return ti[i].Level < ti[j].Level }

// NewGame create new game
func NewGame(database *sql.DB, start, end time.Time) (g Game, err error) {

	g.db = database
	g.Start = start
	g.End = end

	err = g.RecalcScoreboard()
	if err != nil {
		return
	}

	return
}

func (g Game) findTaskByID(id int, tasks []db.Task) (t db.Task, err error) {

	for _, task := range tasks {
		if task.ID == id {
			t = task
			return
		}
	}

	err = errors.New("task no found")

	return
}

// Run open first level tasks and start auto open routine
func (g Game) Run() (err error) {

	for time.Now().Before(g.Start) {
		time.Sleep(time.Second)
	}

	cats, err := g.Tasks()
	if err != nil {
		return
	}

	for _, c := range cats {
		for _, t := range c.TasksInfo {
			err = db.SetOpened(g.db, t.ID, true)
			if err != nil {
				return
			}

			break
		}
	}

	return
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

	tasks, err := db.GetTasks(g.db)
	if err != nil {
		return
	}

	categories, err := db.GetCategories(g.db)
	if err != nil {
		return
	}

	for _, category := range categories {

		cat := CategoryInfo{Name: category.Name}

		for _, task := range tasks {

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
					Level:    task.Level,
				}

				cat.TasksInfo = append(cat.TasksInfo, tInfo)
			}
		}

		sort.Sort(byLevel(cat.TasksInfo))

		cats = append(cats, cat)
	}

	return
}

// Scoreboard returns sorted scoreboard
func (g Game) Scoreboard() (scores []TeamScoreInfo, err error) {

	g.scoreboardLock.Lock()

	teams, err := db.GetTeams(g.db)
	if err != nil {
		return
	}

	for _, team := range teams {

		if team.Test {
			continue
		}

		var s db.Score
		s, err = db.GetLastScore(g.db, team.ID)
		if err != nil {
			return
		}

		scores = append(scores,
			TeamScoreInfo{team.ID, team.Name, team.Desc, s.Score})
	}

	sort.Sort(byScore(scores))

	g.scoreboardLock.Unlock()

	return
}

// RecalcScoreboard update scoreboard
func (g Game) RecalcScoreboard() (err error) {

	g.scoreboardLock.Lock()

	teams, err := db.GetTeams(g.db)
	if err != nil {
		return
	}

	tasks, err := db.GetTasks(g.db)
	if err != nil {
		return
	}

	for _, team := range teams {

		if team.Test {
			continue
		}

		score := 0

		for _, task := range tasks {

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

	tasks, err := db.GetTasks(g.db)
	if err != nil {
		return
	}

	for _, task := range tasks {
		// If same category and next level
		if t.CategoryID == task.CategoryID && t.Level+1 == task.Level {
			// If not already opened
			if !task.Opened {
				// Open it!
				err = db.SetOpened(g.db, task.ID, true)
				if err != nil {
					return
				}
			}
		}
	}

	return
}

func (g Game) isTestTeam(teamID int) bool {

	teams, err := db.GetTeams(g.db)
	if err != nil {
		log.Println("Get teams fail:", err)
		return true
	}

	for _, team := range teams {
		if team.ID == teamID {
			return team.Test
		}
	}

	return false
}

// Solve check flag for task and recalc scoreboard if flag correct
func (g Game) Solve(teamID, taskID int, flag string) (solved bool, err error) {

	tasks, err := db.GetTasks(g.db)
	if err != nil {
		return
	}

	for _, task := range tasks {
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
