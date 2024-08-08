import { term } from "../Terminal";
import WS from "./socket";

export default function socketEventHandler(event) {
  const response = JSON.parse(event.data);
  switch (response.event) {
    case "server-send-files":
      recivingFilesFromServer(response.data);
      break;
    case "command-response":
      commandResponseHandler(response.data);
      break;
    default:
      console.log("event not detected");
  }
}

function recivingFilesFromServer(data) {
  const files = data["files"];

  try {
    let localFiles = localStorage.getItem("files");
    if (localFiles === null) {
      localStorage.setItem("files", JSON.stringify(files));
      document.dispatchEvent(new Event("storage"));
      return;
    }

    localFiles = JSON.parse(localFiles);

    const tempFile = Object.values(files)[0];
    if (localFiles[String(tempFile.path)] !== undefined) {
      localFiles[String(tempFile.path)].files = true;
    }

    for (const key in files) {
      if (!localFiles.hasOwnProperty(key)) {
        localFiles[key] = files[key];
      }
    }

    localStorage.setItem("files", JSON.stringify(localFiles));
    document.dispatchEvent(new Event("storage"));
  } catch (err) {
    console.log("error getting localfiles: ", err);
  }
}

export function requestFiles(path) {
  const t = {
    event: "send-files",
    data: {
      path,
    },
  };
  WS.send(JSON.stringify(t));
  console.log("Send Files");
}

export function fileChanged(file) {
  const t = {
    event: "file-changes",
    data: {
      file: file.path,
      content: file.content,
    },
  };

  WS.send(JSON.stringify(t));
  console.log("Changes send");
}

export function sendCommand(command) {
  const t = {
    event: "command",
    data: {
      command: command,
    }
  }
  console.log("sending ", t)

  WS.send(JSON.stringify(t))
}

//TODO: handle reciving command from server.
export function commandResponseHandler(data) {
  const response = data["response"];
  term.write(response + " ")
}
