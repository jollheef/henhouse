/**
 * @file marshal.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date November, 2015
 * @brief test unmarshaling
 */

package config

import (
	"errors"
	"testing"
)

func TestFailDurationUnmarshal(*testing.T) {

	var d _duration

	err := d.UnmarshalTOML([]byte("testing"))
	if err == nil {
		panic(errors.New("Invalid unmarshaling"))
	}
}

func TestFailTimeUnmarshal(*testing.T) {

	var t _time

	err := t.UnmarshalTOML([]byte("testing"))
	if err == nil {
		panic(errors.New("Invalid unmarshaling"))
	}
}
