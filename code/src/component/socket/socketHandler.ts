import { File } from "../../interface";

export function requestFiles(ws: WebSocket | null | undefined, path: string) {
  if (ws === null || ws === undefined) return;
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

export function sendFileChanges(
  ws: WebSocket | null | undefined,
  file: File,
  changes: string | undefined,
) {
  if (ws === null || ws === undefined) return;
  if (changes === undefined) return;
  const t = {
    event: "file-changes",
    data: {
      file: file.path + "/" + file.name,
      content: changes,
    },
  };

  ws.send(JSON.stringify(t));
}

export function newFileOrDirChanges(
  ws: WebSocket | null | undefined,
  file: File,
) {
  if (ws === null || ws === undefined) return;
  const t = {
    event: "new-file-or-dir",
    data: {
      file: file,
    },
  };
  ws.send(JSON.stringify(t));
}

export function deleteFileOrDir(ws: WebSocket | null | undefined, file: File) {
  if (ws === null || ws === undefined) return;
  const t = {
    event: "delete-file-or-dir",
    data: {
      file: file,
    },
  };
  ws.send(JSON.stringify(t));
}
