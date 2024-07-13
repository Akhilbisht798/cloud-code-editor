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
  console.log(files);
}

function commandResponseHandler(data) {
  const response = data["response"];
  console.log(response);
}
