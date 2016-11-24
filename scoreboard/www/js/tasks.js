if (location.protocol == 'https:')
    var protocol = "wss://";
else
    var protocol = "ws://";

var scoreboard = new WebSocket(protocol + location.host + "/tasks");

scoreboard.onmessage = function(e) {
    document.getElementById('tasks-table').innerHTML = e.data
}

var info = new WebSocket(protocol + location.host + "/info");

info.onmessage = function(e) {
    document.getElementById('info').innerHTML = e.data
}
