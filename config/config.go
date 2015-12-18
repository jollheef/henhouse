/**
 * @file config.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date November, 2015
 * @brief read configuration
 *
 * Contain functions for read configuration file
 */

package config

import (
	"github.com/naoina/toml"
	"io/ioutil"
	"os"
)

// Config is a configuration of henhouse
type Config struct {
	// All logs redirected to file
	LogFile string

	// Path to directory contains task xml files
	TaskDir string

	Database struct {
		Connection     string
		MaxConnections int
		SafeReinit     bool
	}

	Scoreboard struct {
		WwwPath       string
		Addr          string
		RecalcTimeout _duration
	}

	WebsocketTimeout struct {
		Info       _duration
		Scoreboard _duration
		Tasks      _duration
	}

	TaskPrice struct {
		UseTeamsBase           bool
		TeamsBase              int
		P500, P400, P300, P200 int
	}

	Game struct {
		Start _time
		End   _time
	}

	Flag struct {
		// Timeout between send flags
		SendTimeout _duration
	}

	Task struct {
		// Timeout after send correct flag before open next task
		OpenTimeout _duration
		// Auto open task after previous solved
		AutoOpen        bool
		AutoOpenTimeout _duration
	}

	Teams []struct {
		Name        string
		Description string
		Token       string
		Test        bool
	}
}

// ReadConfig read file and return configuration
func ReadConfig(path string) (cfg Config, err error) {

	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	err = toml.Unmarshal(buf, &cfg)
	if err != nil {
		return
	}

	return
}
