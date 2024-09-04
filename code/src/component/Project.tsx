import { FC } from "react";
import { File } from "../interface";
import FileComponent from "./ProjectFiles/FileComponent";
import useWebSocket from "../hooks/useWebSocket";
import SocketProvider from "../context/socketContextProvider";

const Project: FC = () => {
  const ws = useWebSocket()
  const rootFile: File = {
    path: "..",
    name: "client",
    isDir: true,
  };

  return (
    <div>
      <SocketProvider.Provider value={{ ws }}>
        <FileComponent {...rootFile} />
      </SocketProvider.Provider>
    </div>
  );
};

export default Project;
