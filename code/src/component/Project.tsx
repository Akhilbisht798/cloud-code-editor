import { FC, useEffect, useRef } from "react";
import { ROOT_DIR, SOCKET_SERVER } from "../global";
import { File } from "../interface";
import { requestFiles } from "./socket/socketHandler";
import { useProjectFiles } from "../state/projectFilesState";
import FileComponent from "./ProjectFiles/FileComponent";

interface Data {
  files: { [key: string]: File };
}

const Project: FC = () => {
  const ws = useRef<WebSocket | null>(null);
  const { setFiles } = useProjectFiles();
  const rootFile: File = {
    path: "..",
    name: "client",
    isDir: true,
  };
  // const [projectFiles, setProjectFiles] = useState<File[]>([]);

  const handleFilesFromServer = (data: Data) => {
    const files = data["files"];
    setFiles(files);
    console.log(files);
  };

  useEffect(() => {
    ws.current = new WebSocket(SOCKET_SERVER);

    ws.current.onopen = () => {
      console.log("hello world");
      // requestFiles(ws.current, ROOT_DIR);
      setFiles({
        [rootFile.path + "/" + rootFile.name]: rootFile,
      });
    };

    ws.current.onmessage = function socketEventHandler(event) {
      const response = JSON.parse(event.data);
      switch (response.event) {
        case "server-send-files":
          handleFilesFromServer(response.data);
          break;
        // case "command-response":
        //   handleCommandResponse(response.data);
        //   break;
        default:
          console.log("event not detected");
      }
    };

    ws.current.onerror = (error) => {
      console.log("WebSocket error ", error);
    };

    return () => {
      if (ws.current) {
        ws.current.close();
        console.log("WebSocket disconnected");
      }
    };
  }, []);
  return (
    <div>
      <FileComponent {...rootFile} />
    </div>
  );
};

export default Project;
