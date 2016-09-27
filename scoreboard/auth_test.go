/**
 * @file auth.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2016
 * @brief test auth
 */

package scoreboard

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jollheef/henhouse/db"
)

func TestGenSession(*testing.T) {
	s, err := genSession()
	if err != nil {
		panic(err)
	}

	if len(s) != 512 {
		panic("Session length mismatch")
	}
}

func TestSessionTeamID(*testing.T) {

	database, err := db.InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	defer database.Close()

	err = addTestData(database, 20, 5, 5, "testFlag")
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()

	realTeamID := 1

	err = setSessionTeamID(database, w, realTeamID)
	if err != nil {
		panic(err)
	}

	resp := http.Response{Header: w.Header()}
	cookies := resp.Cookies()
	if len(cookies) == 0 {
		panic(err)
	}

	req := &http.Request{Header: http.Header{
		"Cookie": w.HeaderMap["Set-Cookie"]}}

	teamID, err := getSessionTeamID(database, req)
	if err != nil {
		panic(err)
	}

	if teamID != realTeamID {
		panic("teamID != realTeamID")
	}
}
