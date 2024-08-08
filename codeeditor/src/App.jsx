import { useEffect } from "react";
import Editor from "./components/Editor";
import WS from "./components/socket/socket";
import FileExplorer from "./components/FileExplorer/FileExplorer";
import Terminal from "./components/Terminal";

export default function App() {
  useEffect(() => {
    WS;
  }, []);

  function onClickHandler() {
    console.log("sending");
    //let command = "pwd";
  }

  return (
    <>
      <div className=" flex gap-3 resize-y overflow-auto min-h-10">
        <div className="w-1/4 resize-x overflow-auto">
          <FileExplorer />
        </div>
        <div className="w-3/4 ">
          <Editor  />
        </div>
      </div>
      <div className=" resize-y overflow-auto p-4 m-3">
        <Terminal />
      </div>
    </>
  );
}
