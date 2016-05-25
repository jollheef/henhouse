/**
 * @file task_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2015
 * @brief test parse task xml
 */

package config

import (
	"fmt"
	"testing"
)

func TestParseXML(*testing.T) {

	name := "bar"
	description := "fofofo"
	category := "test"
	level := 1
	flag := "justflag"

	xml := fmt.Sprintf(`
	<Task>
	  <Name>%s</Name>
	  <Description>%s</Description>
	  <Category>%s</Category>
	  <Level>%d</Level>
	  <Flag>%s</Flag>
	</Task>`, name, description, category, level, flag)

	task, err := ParseXMLTask([]byte(xml))
	if err != nil {
		panic(err)
	}

	if task.Name != name {
		panic("invalid parse")
	}

	if task.Description != description {
		panic("invalid parse")
	}

	if task.Category != category {
		panic("invalid parse")
	}

	if task.Level != level {
		panic("invalid parse")
	}

	if task.Flag != flag {
		panic("invalid parse")
	}
}
