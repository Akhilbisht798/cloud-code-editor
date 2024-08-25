import { useEffect } from "react";
import WS from "./components/socket/socket";
import Project from "./components/Project";
import { useState } from "react";

export default function App() {
  const [render, setRender] = useState(true);
  useEffect(() => {
    WS;
  }, []);

  function onClickHandler() {
    setRender(!render);
  }

  return (
    <>
      <button onClick={onClickHandler}>Show</button>
      {render ? <Project /> : <>Nothing to show you.</>}
    </>
  );
}
