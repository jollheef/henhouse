/**
 * @file template_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2016
 */

package scoreboard

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestGetTmpl(*testing.T) {
	_, err := getTmpl("")
	if err == nil {
		log.Fatal("Get template that not really exists")
	}

	content := []byte("test")
	templatePath = "/tmp"

	err = ioutil.WriteFile(templatePath+"/test.htmlf", content, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	s, err := getTmpl("test")
	if err != nil {
		log.Fatal(err)
	}

	if s != string(content) {
		log.Fatalln("Content is wrong")
	}

	os.Remove(templatePath + "/test.htmlf")

	_, err = getTmpl("test")
	if err != nil {
		log.Fatal("Cache not work")
	}
}
