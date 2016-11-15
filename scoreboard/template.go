/**
 * @file template.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2016
 * @brief simple printf templates
 */

package scoreboard

import "io/ioutil"

var templatePath string

var cache = make(map[string]string)

func getTmpl(name string) (s string, err error) {
	s, ok := cache[name]
	if ok {
		return
	}

	b, err := ioutil.ReadFile(templatePath + "/" + name + ".htmlf")
	if err != nil {
		return
	}

	s = string(b)
	cache[name] = s
	return
}

func getTmplWoCache(name string) (s string, err error) {
	b, err := ioutil.ReadFile(templatePath + "/" + name + ".htmlf")
	if err != nil {
		return
	}
	s = string(b)
	return
}
