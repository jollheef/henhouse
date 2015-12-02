var scoreboard = new WebSocket("ws://" + location.host + "/tasks");

scoreboard.onmessage = function(e) {
    document.getElementById('tasks-table').innerHTML = e.data
}

var info = new WebSocket("ws://" + location.host + "/tasks-info");

info.onmessage = function(e) {
    document.getElementById('info').innerHTML = e.data
}
