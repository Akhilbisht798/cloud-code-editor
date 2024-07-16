import { useEffect } from "react";
import Editor from "./components/Editor";
import WS from "./components/socket/socket";
import FilesViewer from "./components/FilesViewer";

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
      <h1 className="text-3xl font-bold underline">Hello world!</h1>;
      <Editor />
      <FilesViewer />
      <button onClick={onClickHandler}>Send Command</button>
    </>
  );
}
