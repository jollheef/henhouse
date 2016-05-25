/**
 * @file task.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2015
 * @brief task parser
 *
 * Contain functions for parse task xml
 */

package config

import (
	"encoding/xml"
)

// Task is xml task data model
type Task struct {
	Name        string
	Description string
	Category    string
	Level       int
	Flag        string
	Author      string
}

// ParseXMLTask parse xml task
func ParseXMLTask(rawXML []byte) (task Task, err error) {
	err = xml.Unmarshal(rawXML, &task)
	return
}
