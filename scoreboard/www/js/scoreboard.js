var scoreboard = new WebSocket("ws://" + location.host + "/scoreboard");

scoreboard.onmessage = function(e) {
    document.getElementById('scoreboard-table').innerHTML = e.data
}

var info = new WebSocket("ws://" + location.host + "/info");

info.onmessage = function(e) {
    document.getElementById('info').innerHTML = e.data
}
