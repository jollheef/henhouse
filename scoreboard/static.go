/**
 * @file static.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
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

	fmt.Fprintf(w, `<!DOCTYPE html>
<html class="full" lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Juniors CTF</title>

    <link rel="stylesheet" href="https://bootswatch.com/yeti/bootstrap.min.css">
    <link rel="stylesheet" href="css/style.css">

    <script type="text/javascript" src="js/scoreboard.js"></script>

  </head>
  <body>
    <ul class="nav nav-tabs h4">
      <li class="active">
        <a href="#">Scoreboard</a>
      </li>
      <li><a href="tasks.html">Tasks</a></li>
      <li><a href="news.html">News</a></li>
    </ul>
    <div style="padding: 15px;">
      <div id="info">%s</div>
      <br>
      <center><table id="scoreboard-table" class="table table-hover h2">%s</table></center>
      <center><img id="juniorstext" src="/images/juniors_ctf_txt.png"></center>
    </div>
  </body>
</html>`, getInfo(lastScoreboardUpdated), currentResult)
}

func staticTasks(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, `<!DOCTYPE html>
<html class="full" lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Juniors CTF</title>

    <link rel="stylesheet" href="https://bootswatch.com/yeti/bootstrap.min.css">
    <link rel="stylesheet" href="css/style.css">

    <script type="text/javascript" src="js/tasks.js"></script>

  </head>
  <body>
    <ul class="nav nav-tabs h4">
      <li><a href="index.html">Scoreboard</a></li>
      <li class="active">
        <a href="#">Tasks</a>
      </li>
      <li><a href="news.html">News</a></li>
    </ul>
    <div style="padding: 15px;">
      <div id="info">%s</div>
      <br>
      <center><table id="tasks-table" class="table table-hover">%s</table></center>
    </div>
  </body>
</html>`, getInfo(lastTasksUpdated), currentTasks)
}
