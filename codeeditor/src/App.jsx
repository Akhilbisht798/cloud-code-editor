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
      <div className=" flex ">
        <div className="w-1/4">
          <FileExplorer />
        </div>
        <Editor className="w-3/4" />
      </div>
      <div>
        <Terminal />
      </div>
      <div>
        <button onClick={onClickHandler}>Send Command</button>
      </div>
    </>
  );
}
