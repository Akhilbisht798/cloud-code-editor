import { FC } from "react";
import { File } from "../interface";
import FileComponent from "./ProjectFiles/FileComponent";
import useWebSocket from "../hooks/useWebSocket";
import SocketProvider from "../context/socketContextProvider";
import Editor from "./editor/CodeEditor";
import Terminal from "./terminal/terminal";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";

const Project: FC = () => {
  const ws = useWebSocket();
  const rootFile: File = {
    path: "..",
    name: "client",
    isDir: true,
  };

  return (
    <SocketProvider.Provider value={{ ws }}>
      <ResizablePanelGroup direction="horizontal" className="w-full rounded-lg border border-[#1e1e1e] bg-[#1e1e1e] text-[#cccccc]">
        <ResizablePanel defaultSize={20} minSize={5}>
          <div className="border-b border-[#3c3c3c] p-2 font-semibold">File Explorer</div>
          <div className="overflow-auto h-[calc(100vh-40px)]">
            <FileComponent {...rootFile} />
          </div>
        </ResizablePanel>

        <ResizableHandle withHandle className="bg-[#3c3c3c] hover:bg-[#505050]" />

        <ResizablePanel defaultSize={80}>
          <ResizablePanelGroup direction="vertical">
            <ResizablePanel defaultSize={60} minSize={20}>
              <Editor />
            </ResizablePanel>
            <ResizableHandle withHandle className="bg-[#3c3c3c] hover:bg-[#505050]"/>
            <ResizablePanel defaultSize={40} minSize={20}>
              <Terminal />
            </ResizablePanel>
          </ResizablePanelGroup>
        </ResizablePanel>
      </ResizablePanelGroup>
    </SocketProvider.Provider>
  );
};

export default Project;
