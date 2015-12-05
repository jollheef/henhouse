/**
 * @file auth.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date December, 2015
 * @brief auth helpers and middleware
 */

package scoreboard

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"github.com/jollheef/henhouse/db"
	"net/http"
)

const sessionCookieName = "session"

func genSession() (s string, err error) {

	session_len := 256

	rand_buf := make([]byte, session_len)

	_, err = rand.Read(rand_buf)

	if err != nil {
		return
	}

	s = fmt.Sprintf("%x", rand_buf)

	return
}

func getSessionTeamID(database *sql.DB, r *http.Request) (teamID int,
	err error) {

	session, err := r.Cookie(sessionCookieName)
	if err != nil {
		return
	}

	teamID, err = db.GetSessionTeam(database, session.Value)
	if err != nil {
		return
	}

	return
}

func setSessionTeamID(database *sql.DB, w http.ResponseWriter,
	teamID int) (err error) {

	session, err := genSession()
	if err != nil {
		return
	}

	cookie := http.Cookie{Name: sessionCookieName, Value: session}

	err = db.AddSession(database, &db.Session{
		TeamID:  teamID,
		Session: session,
	})
	if err != nil {
		return
	}

	http.SetCookie(w, &cookie)

	return
}

func authorized(database *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := getSessionTeamID(database, r)
		if err != nil {
			http.Redirect(w, r, "/auth.html", 307)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
