import socketEventHandler, { requestFiles } from "./socketHandler";

const WS = new WebSocket("ws://localhost:5000/ws");
export const RootDir = "userId-1/client";

WS.onopen = function () {
  console.log("Connection opened.");
  requestFiles(RootDir);
};

WS.onmessage = socketEventHandler;

WS.onerror = function (event) {
  console.error("WebSocket error:", event);
};

WS.onclose = function () {
  console.log("Connection closed.");
};

export default WS;
