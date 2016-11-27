/**
 * @file main.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date December, 2015
 * @brief contest checking system CLI
 *
 * Entry point for contest checking system CLI
 */

package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/jollheef/henhouse/config"
	"github.com/jollheef/henhouse/db"
	"github.com/jollheef/henhouse/game"
	"github.com/olekukonko/tablewriter"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	configPath = kingpin.Flag("config", "Path to configuration file.").String()

	// Task
	task = kingpin.Command("task", "Work with tasks.")

	taskList        = task.Command("list", "List tasks.")
	taskListWOFlags = taskList.Flag("without-flags", "Do not include flags in output.").Bool()

	taskAdd    = task.Command("add", "Add task.")
	taskAddXML = taskAdd.Arg("xml", "Path to xml.").Required().String()

	taskUpdate    = task.Command("update", "Update task.")
	taskUpdateID  = taskUpdate.Arg("id", "ID of task.").Required().Int()
	taskUpdateXML = taskUpdate.Arg("xml", "Path to xml.").Required().String()

	taskOpen   = task.Command("open", "Open task.")
	taskOpenID = taskOpen.Arg("id", "ID of task").Required().Int()

	taskClose   = task.Command("close", "Close task.")
	taskCloseID = taskClose.Arg("id", "ID of task").Required().Int()

	taskDump   = task.Command("dump", "Dump task to xml.")
	taskDumpID = taskDump.Arg("id", "ID of task").Required().Int()

	// Category
	category = kingpin.Command("category", "Work with categories.")

	categoryList = category.Command("list", "List categories.")

	categoryAdd  = category.Command("add", "Add category.")
	categoryName = categoryAdd.Arg("name", "Name.").Required().String()

	// Team
	team = kingpin.Command("team", "Work with teams.")

	teamList   = team.Command("list", "List teams.")
	teamInfo   = team.Command("info", "Information about team.")
	teamInfoID = teamInfo.Arg("id", "ID of task").Required().Int()

	// Export
	export               = kingpin.Command("export", "Export scoreboard for ctftime.")
	exportWithLastAccept = export.Flag("with-last-accept", "Add last-accept field.").Bool()
)

var (
	// CommitID fill in ./build.sh
	CommitID string
	// BuildDate fill in ./build.sh
	BuildDate string
	// BuildTime fill in ./build.sh
	BuildTime string
)

func getCategoryByID(categoryID int, categories []db.Category) string {
	for _, cat := range categories {
		if cat.ID == categoryID {
			return cat.Name
		}
	}
	return "Unknown"
}

func getCategoryByName(name string, categories []db.Category) (id int, err error) {
	for _, cat := range categories {
		if cat.Name == name {
			return cat.ID, nil
		}
	}

	return 0, errors.New("Category " + name + " not found")
}

type byID []db.Task

func (t byID) Len() int           { return len(t) }
func (t byID) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t byID) Less(i, j int) bool { return t[i].ID < t[j].ID }

func parseTask(path string, categories []db.Category) (t db.Task, err error) {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	task, err := config.ParseXMLTask(content)
	if err != nil {
		return
	}

	t.Name = task.Name
	t.Desc = task.Description
	t.NameEn = task.NameEn
	t.DescEn = task.DescriptionEn
	t.CategoryID, err = getCategoryByName(task.Category, categories)
	if err != nil {
		log.Println("Cant find category")
		return
	}

	t.Level = task.Level
	t.Flag = task.Flag
	t.Price = 500         // TODO support non-shared task
	t.Shared = true       // TODO support non-shared task
	t.MaxSharePrice = 500 // TODO support value from xml
	t.MinSharePrice = 100 // TODO support value from xml
	t.Opened = false      // by default task is closed
	t.Author = task.Author
	t.Tags = task.Tags
	t.ForceClosed = task.ForceClosed

	return
}

var cfgFiles = []string{"/etc/henhouse/cli.toml", "/etc/henhouse.toml",
	"cli.toml", "henhouse.toml"}

func taskUpdateCmd(database *sql.DB, categories []db.Category) (err error) {
	task, err := db.GetTask(database, *taskUpdateID)
	if err != nil {
		return
	}

	id := task.ID

	task, err = parseTask(*taskUpdateXML, categories)
	if err != nil {
		return
	}

	task.ID = id

	err = db.UpdateTask(database, &task)
	if err != nil {
		return
	}

	return
}

func taskAddCmd(database *sql.DB, categories []db.Category) (err error) {
	t, err := parseTask(*taskAddXML, categories)
	if err != nil {
		return
	}

	err = db.AddTask(database, &t)
	if err != nil {
		return
	}

	return
}

func taskListCmd(database *sql.DB, categories []db.Category) (err error) {
	tasks, err := db.GetTasks(database)
	if err != nil {
		return
	}

	sort.Sort(byID(tasks))

	table := tablewriter.NewWriter(os.Stdout)
	var header []string
	if *taskListWOFlags {
		header = []string{"ID", "Name", "Category", "Opened", "Solved by"}
	} else {
		header = []string{"ID", "Name", "Category", "Flag", "Opened", "Solved by"}
	}
	table.SetHeader(header)

	for _, task := range tasks {
		var solvedBy int
		solvedBy, err = db.GetSolvedCount(database, task.ID)
		if err != nil {
			return
		}

		var row []string

		row = append(row, fmt.Sprintf("%d", task.ID))
		row = append(row, task.Name)
		row = append(row, getCategoryByID(task.CategoryID, categories))
		if !*taskListWOFlags {
			row = append(row, task.Flag)
		}
		row = append(row, fmt.Sprintf("%v", task.Opened))
		row = append(row, fmt.Sprintf("%d", solvedBy))

		table.Append(row)
	}

	table.Render()

	return
}

func taskDumpCmd(database *sql.DB, categories []db.Category) (err error) {
	task, err := db.GetTask(database, *taskDumpID)
	if err != nil {
		return
	}

	xmlTask := config.Task{
		Name:          task.Name,
		NameEn:        task.NameEn,
		Description:   task.Desc,
		DescriptionEn: task.DescEn,
		Category:      getCategoryByID(task.CategoryID, categories),
		Level:         task.Level,
		Flag:          task.Flag,
		Author:        task.Author,
		ForceClosed:   task.ForceClosed,
		Tags:          task.Tags,
	}

	output, err := xml.MarshalIndent(xmlTask, "", "	")
	if err != nil {
		return

	}

	fmt.Fprintln(os.Stdout, string(output))

	return
}

func categoryListCmd(database *sql.DB) (err error) {
	categories, err := db.GetCategories(database)
	if err != nil {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name"})

	for _, cat := range categories {
		row := []string{fmt.Sprintf("%d", cat.ID), cat.Name}
		table.Append(row)
	}

	table.Render()

	return
}

type byScore [][]string

func (t byScore) Len() int      { return len(t) }
func (t byScore) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t byScore) Less(i, j int) bool {
	var in, jn int
	fmt.Sscanf(t[i][2], "%d", &in)
	fmt.Sscanf(t[j][2], "%d", &jn)
	return in > jn
}

func teamListCmd(database *sql.DB) (err error) {
	teams, err := db.GetTeams(database)
	if err != nil {
		return
	}

	flags, err := db.GetFlags(database)
	if err != nil {
		return
	}

	tasks, err := db.GetTasks(database)
	if err != nil {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Score", "Solved"})

	rows := [][]string{}

	for _, t := range teams {
		row := []string{fmt.Sprintf("%d", t.ID), t.Name}

		var score db.Score
		score, err = db.GetLastScore(database, t.ID)
		if err != nil {
			continue
		}
		row = append(row, fmt.Sprintf("%d", score.Score))

		solvedCount := 0
		for _, f := range flags {
			if f.TeamID == t.ID {
				solvedCount++
			}
		}

		row = append(row, fmt.Sprintf("%d/%d", solvedCount, len(tasks)))

		rows = append(rows, row)
	}

	sort.Sort(byScore(rows))

	for _, row := range rows {
		table.Append(row)
	}

	table.Render()

	return
}

func teamInfoCmd(database *sql.DB) (err error) {
	teams, err := db.GetTeams(database)
	if err != nil {
		return
	}

	flags, err := db.GetFlags(database)
	if err != nil {
		return
	}

	tasks, err := db.GetTasks(database)
	if err != nil {
		return
	}

	for _, t := range teams {
		if t.ID == *teamInfoID {
			fmt.Println("ID:", t.ID)
			fmt.Println("Name:", t.Name)
			fmt.Println("Email:", t.Email)
			fmt.Println("Description:", t.Desc)
			fmt.Println("Token:", t.Token)
			fmt.Println("Test:", t.Test)
			fmt.Print("Solved: ")
			solvedCount := 0
			for _, f := range flags {
				if f.TeamID == t.ID {
					var task db.Task
					task, err = db.GetTask(database, f.TaskID)
					if err != nil {
						return
					}
					solvedCount++
					fmt.Printf("%s (%d), ", task.Name, task.ID)
				}
			}
			fmt.Println()
			fmt.Printf("Solved: %d/%d\n", solvedCount, len(tasks))
		}
	}

	return
}

func exportScoreboard(database *sql.DB) (err error) {

	scores := []game.TeamScoreInfo{}

	teams, err := db.GetTeams(database)
	if err != nil {
		return
	}

	flags, err := db.GetFlags(database)
	if err != nil {
		return
	}

	for _, team := range teams {

		if team.Test {
			continue
		}

		var s db.Score
		s, err = db.GetLastScore(database, team.ID)
		if err != nil {
			return
		}

		var bName []byte

		bName, err = json.Marshal(team.Name)
		if err != nil {
			return
		}

		scoreInfo := game.TeamScoreInfo{
			ID:         team.ID,
			Name:       string(bName),
			Score:      s.Score,
			LastAccept: game.LastAccept(team.ID, flags),
		}

		scores = append(scores, scoreInfo)
	}

	sort.Sort(game.ByScoreAndLastAccept(scores))

	fmt.Println("{\n\t\"standings\": [")
	for i, s := range scores {
		if *exportWithLastAccept {
			fmt.Printf("\t\t{ \"pos\": %d, \"team\": %s,"+
				" \"score\": %d, \"lastAccept\" : %d }",
				i+1, s.Name, s.Score, s.LastAccept)
		} else {
			fmt.Printf("\t\t{ \"pos\": %d, \"team\": %s,"+
				" \"score\": %d }", i+1, s.Name, s.Score)
		}
		if i != len(scores)-1 {
			fmt.Printf(",\n")
		} else {
			fmt.Printf("\n")
		}

	}
	fmt.Println("\t]\n}")

	return
}

func runCommandLine(database *sql.DB, categories []db.Category) (err error) {
	switch kingpin.Parse() {
	case "task add":
		err = taskAddCmd(database, categories)
	case "task update":
		err = taskUpdateCmd(database, categories)
	case "task list":
		err = taskListCmd(database, categories)
	case "task open":
		err = db.SetOpened(database, *taskOpenID, true)
	case "task close":
		err = db.SetOpened(database, *taskCloseID, false)
	case "task dump":
		err = taskDumpCmd(database, categories)
	case "category add":
		err = db.AddCategory(database, &db.Category{Name: *categoryName})
	case "category list":
		err = categoryListCmd(database)
	case "team list":
		err = teamListCmd(database)
	case "team info":
		err = teamInfoCmd(database)
	case "export":
		err = exportScoreboard(database)
	}

	return
}

func main() {

	if len(CommitID) > 7 {
		CommitID = CommitID[:7] // abbreviated commit hash
	}

	kingpin.Version(BuildDate + " " + CommitID +
		" (Mikhail Klementyev <jollheef@riseup.net>)")

	kingpin.Parse()

	var cfgPath string

	if *configPath != "" {
		cfgPath = *configPath
	} else {

		for _, cfgFile := range cfgFiles {
			_, err := os.Stat(cfgFile)
			if err == nil {
				cfgPath = cfgFile
				break
			}
		}
	}

	if cfgPath == "" {
		log.Fatalln("Config not found")
	}

	cfg, err := config.ReadConfig(cfgPath)
	if err != nil {
		log.Fatalln("Cannot open config:", err)
	}

	database, err := db.OpenDatabase(cfg.Database.Connection)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	defer database.Close()

	database.SetMaxOpenConns(cfg.Database.MaxConnections)

	categories, err := db.GetCategories(database)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	err = runCommandLine(database, categories)
	if err != nil {
		log.Fatalln("Error:", err)
	}
}
