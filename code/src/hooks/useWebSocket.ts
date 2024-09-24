import { useEffect, useRef } from "react";
import { SOCKET_SERVER } from "../global";
import { File } from "../interface";
import { useProjectFiles } from "../state/projectFilesState";
import { term } from "@/component/terminal/terminal";

interface CommandResponse {
  response: string;
}

const useWebSocket = () => {
  const ws = useRef<WebSocket | null>(null);
  const { setFiles } = useProjectFiles();

  const rootFile: File = {
    // for production if i think it will be "."
    // and name will be the project name
    path: "..",
    name: "client",
    isDir: true,
  };

  function handleCommandResponse(data: CommandResponse) {
    let res = data["response"];
    res += " ";
    term.write(res);
  }

  useEffect(() => {
    ws.current = new WebSocket(SOCKET_SERVER);

    ws.current.onopen = () => {
      console.log("web socket connected");
      setFiles({
        [rootFile.path + "/" + rootFile.name]: rootFile,
      });
    };

    ws.current.onmessage = function socketEventHandler(event) {
      const response = JSON.parse(event.data);
      switch (response.event) {
        case "server-send-files":
          setFiles(response.data.files);
          console.log("New file added to library.", response.data.files);
          break;
        case "command-response":
          handleCommandResponse(response.data);
          break;
        default:
          console.log("event not detected");
      }
    };

    ws.current.onerror = (error) => {
      console.log("websocket err: ", error);
      if (ws.current) {
        console.log("Websocket ready state: ", ws.current.readyState);
      }
    };

    return () => {
      if (ws.current) {
        ws.current.close();
        console.log("WebSocket disconnected");
      }
    };
  }, []);
  return ws.current;
};

export default useWebSocket;
