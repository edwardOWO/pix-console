// websocket.js

let ws;

function connectWebSocket(path) {
    ws = new WebSocket("ws://localhost:8080/ws?path=" + path);

    ws.onmessage = function (event) {
        const logs = document.getElementById("logs");
        logs.textContent += event.data; // Use += to append new log data
        logs.scrollTop = logs.scrollHeight;
    };

    ws.onclose = function () {
        window.close();
    };

    ws.onerror = function () {
        window.close();
    };
}
