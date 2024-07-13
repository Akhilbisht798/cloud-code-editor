import socketEventHandler from "./socketHandler";

var ws = new WebSocket("ws://localhost:5000/ws");

ws.onopen = function () {
  console.log("Connection opened.");
};

// Handle events.
ws.onmessage = socketEventHandler;

ws.onerror = function (event) {
  console.error("WebSocket error:", event);
};

ws.onclose = function () {
  console.log("Connection closed.");
};

export default ws;
