/**
 * @file config_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2015
 * @brief test config package
 */

package config

import (
	"errors"
	"io/ioutil"
	"testing"
)

func bugOnInvalid(real, parsed string) {
	if real != parsed {
		panic(errors.New("Parsed invalid value"))
	}
}

func TestReadConfig(*testing.T) {

	cfg, err := ReadConfig("henhouse.toml")
	if err != nil {
		panic(err)
	}

	bugOnInvalid("2015-11-17 10:00:00 +0300 MSK", cfg.Game.Start.String())

	bugOnInvalid("2015-12-31 23:59:00 +0300 MSK", cfg.Game.End.String())

	bugOnInvalid("1m0s", cfg.Task.OpenTimeout.String())

	bugOnInvalid("6h0m0s", cfg.Task.AutoOpenTimeout.String())

	// other values has built-in types
}

// Test read config with invalid path
func TestFailReadConfig(*testing.T) {

	_, err := ReadConfig("/dev/ololo/pewpew")
	if err == nil {
		panic(errors.New("Ok read non exist config"))
	}
}

func TestReadInvalidConfig(*testing.T) {

	configPath := "/tmp/invalid-config"

	err := ioutil.WriteFile(configPath, []byte("non toml blabla%%["), 0644)
	if err != nil {
		panic(err)
	}

	_, err = ReadConfig(configPath)
	if err == nil {
		panic(errors.New("Ok read invalid config"))
	}
}
