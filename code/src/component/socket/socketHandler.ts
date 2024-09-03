import { File } from "../../interface";

export function requestFiles(ws: WebSocket | null, path: string) {
  if (ws === null) return;
  const t = {
    event: "send-files",
    data: {
      path,
    },
  };
  ws.send(JSON.stringify(t));
  console.log("file requested");
}

export function sendCommand(ws: WebSocket, command: string): void {
  const t = {
    event: "command",
    data: {
      command: command,
    },
  };
  ws.send(JSON.stringify(t));
}

export function sendFileChanges(ws: WebSocket, file: File, changes: string) {
  const t = {
    event: "file-changes",
    data: {
      file: file.path + "/" + file.name,
      content: changes,
    },
  };

  ws.send(JSON.stringify(t));
}
