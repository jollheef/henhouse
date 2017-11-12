if (location.protocol == 'https:')
    var protocol = "wss://";
else
    var protocol = "ws://";

var scoreboard = new WebSocket(protocol + location.host + "/scoreboard");

scoreboard.onmessage = function(e) {
    document.getElementById('updated-scoreboard-content').innerHTML = e.data
}

var info = new WebSocket(protocol + location.host + "/info");

info.onmessage = function(e) {
    document.getElementById('info').innerHTML = e.data
}
