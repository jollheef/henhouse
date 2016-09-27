/**
 * @file auth.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2016
 * @brief test auth
 */

package scoreboard

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/jollheef/henhouse/db"
)

func testDB() (database *sql.DB) {
	database, err := db.InitDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	err = addTestData(database, 20, 5, 5, "testFlag")
	if err != nil {
		panic(err)
	}

	return
}

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

	database := testDB()
	defer database.Close()

	w := httptest.NewRecorder()

	realTeamID := 1

	err := setSessionTeamID(database, w, realTeamID)
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

func TestAuthHandlerGet(*testing.T) {

	database := testDB()
	defer database.Close()

	r := httptest.NewRequest("GET", "http://localhost", nil)
	w := httptest.NewRecorder()

	authHandler(database, w, r)

	if w.Code != http.StatusTemporaryRedirect {
		panic("wrong status")
	}
}

func TestAuthHandlerWithoutToken(*testing.T) {

	database := testDB()
	defer database.Close()

	r := httptest.NewRequest("POST", "http://localhost", nil)
	w := httptest.NewRecorder()

	authHandler(database, w, r)

	if w.Code != http.StatusTemporaryRedirect {
		panic("wrong status")
	}
}

func TestAuthHandlerWithWrongToken(*testing.T) {

	database := testDB()
	defer database.Close()

	r := httptest.NewRequest("POST", "http://localhost", nil)
	w := httptest.NewRecorder()

	r.Form = url.Values{}
	r.Form.Set("token", "WRONGTOKEN")

	authHandler(database, w, r)

	if w.Code != http.StatusUnauthorized {
		panic("wrong status")
	}
}

func TestLogonLogout(*testing.T) {

	database := testDB()
	defer database.Close()

	r := httptest.NewRequest("POST", "http://localhost", nil)
	w := httptest.NewRecorder()

	r.Form = url.Values{}
	r.Form.Set("token", "l") // TODO Fix hardcoded valid token

	authHandler(database, w, r)

	if w.Code != http.StatusSeeOther { // success
		panic("wrong status")
	}

	r2 := httptest.NewRequest("POST", "http://localhost", nil)
	r2.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}
	w2 := httptest.NewRecorder()

	logoutHandler(w2, r2)

	if w2.Code != http.StatusTemporaryRedirect {
		panic("wrong status")
	}

	if w2.HeaderMap["Set-Cookie"][0] != "session=" { // empty session
		panic("logout does not remove cookies")
	}
}
