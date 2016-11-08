/**
 * @file static.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date December, 2015
 * @brief non-dynamic html results
 *
 * Generate static html page
 */

package scoreboard

import (
	"fmt"
	"net/http"
)

func staticScoreboard(w http.ResponseWriter, r *http.Request) {

	teamID := getTeamID(r)

	fmt.Fprintf(w, `<!DOCTYPE html>
<html class="full" lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="shortcut icon" href="images/favicon.png" type="image/png">
    <title>Juniors CTF</title>

    <link rel="stylesheet" href="css/style.css" class="--apng-checked">

    <script type="text/javascript" src="js/scoreboard.js"></script>

  </head>
  <body>
    <ul id="header">
      <li class="header_link active"><a href="#">Scoreboard</a></li>
      <li class="header_link"><a href="tasks.html">Tasks</a></li>
      <li class="header_link"><a href="news.html">News</a></li>
      <li class="header_link"><a href="sponsors.html">Sponsors</a></li>
      <li id="info">%s</li>
    </ul>
    <div id="content">
      <table id="scoreboard-table">%s</table>
    </div>
    <div class="center"><img id="juniorstext" src="./images/juniors_ctf_txt.png"></div>
  </body>
</html>`, getInfo(), scoreboardHTML(teamID))
}

func staticTasks(w http.ResponseWriter, r *http.Request) {

	teamID := getTeamID(r)

	fmt.Fprintf(w, `<!DOCTYPE html>
<html class="full" lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="shortcut icon" href="images/favicon.png" type="image/png">
    <title>Juniors CTF</title>

    <link rel="stylesheet" href="css/style.css">

    <script type="text/javascript" src="js/tasks.js"></script>

  </head>
  <body>
    <ul id="header">
      <li class="header_link"><a href="scoreboard.html">Scoreboard</a></li>
      <li class="header_link active"><a href="#">Tasks</a></li>
      <li class="header_link"><a href="news.html">News</a></li>
      <li class="header_link"><a href="sponsors.html">Sponsors</a></li>
      <li id="info">%s</li>
    </ul>

    <div id="content">
      <div id="tasks-table">%s</div> 
    </div>
  </body>
</html>`, getInfo(), tasksHTML(teamID))
}
