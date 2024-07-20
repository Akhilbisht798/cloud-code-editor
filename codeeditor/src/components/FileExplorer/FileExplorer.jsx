import { useState } from "react";
import { useEffect } from "react";
import { requestFiles } from "../socket/socketHandler";
import FileComponent from "./FileComponent";

export default function FileExplorer() {
  const [files, setFiles] = useState(null);

  useEffect(() => {
    function handleFilesChanges() {
      let localfiles = localStorage.getItem("files");
      localfiles = JSON.parse(localfiles);
      setFiles(localfiles);
    }

    document.addEventListener("storage", handleFilesChanges);

    return () => {
      document.removeEventListener("storage", handleFilesChanges);
    };
  }, []);

  useEffect(() => {
    console.log("files changed");
  }, [files]);

  return (
    <>
      {files && Object.keys(files).length > 0 ? (
        Object.entries(files).map(([key, value]) => (
          <FileComponent key={key} file={value} dictFiles={files} />
        ))
      ) : (
        <div>No file available</div>
      )}
    </>
  );
}
