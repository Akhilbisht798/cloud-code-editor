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
  const dir = data["dir"];
  console.log("directory From server: " + dir);
  try {
    let localFiles = localStorage.getItem("files");
    if (localFiles === null) {
      localStorage.setItem("files", JSON.stringify(files));
      document.dispatchEvent(new Event("storage"));
      return;
    }

    localFiles = JSON.parse(localFiles);
    for (let i = 0; i < localFiles.length; i++) {
      const file = localFiles[i];
      const filePath = file.path + "/" + file.name;
      if (filePath === dir) {
        console.log("changed localfiles");
        file["files"] = files;
        break;
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
}

//TODO: handle reciving command from server.
function commandResponseHandler(data) {
  const response = data["response"];
  console.log(response);
}
