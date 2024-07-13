import { useEffect } from "react";
import Editor from "./components/Editor";
import ws from "./components/socket/socket";

export default function App() {
  useEffect(() => {
    ws;
  }, []);

  function onClickHandler() {
    console.log("sending");
    //let command = "pwd";
    const t = {
      event: "send-files",
      data: {
        path: "../",
      },
    };
    ws.send(JSON.stringify(t));
  }

  return (
    <>
      <h1 className="text-3xl font-bold underline">Hello world!</h1>;
      <Editor />
      <button onClick={onClickHandler}>Send Command</button>
    </>
  );
}
