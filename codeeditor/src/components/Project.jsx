import Editor from "./Editor";
import Terminal from "./Terminal";
import FileExplorer from "./FileExplorer/FileExplorer";
import { useLayoutEffect } from "react";

export default function Project() {
  useLayoutEffect(() => {
    function close() {
      console.log("close the docker");
    }
    return () => close();
  }, []);
  return (
    <>
      <div className=" flex gap-3 resize-y overflow-auto min-h-[50vh]">
        <div className="w-1/4 resize-x overflow-auto">
          <FileExplorer />
        </div>
        <div className="w-3/4">
          <Editor />
        </div>
      </div>
      <div className=" resize-y overflow-auto p-4 m-3">
        <Terminal />
      </div>
    </>
  );
}
