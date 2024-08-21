import { useState } from "react";
import { useEffect } from "react";
import FileComponent from "./FileComponent";
import { RootDir } from "../socket/socket";
import useFiles from "../../state/files";

export default function FileExplorer() {
  const {files, setFiles } = useFiles()

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
  }, [files])
  //TODO: Fix the files not rendering in list form properly.
  //TODO: inconsistent file explorer.
  return (
    <>
        {files && Object.keys(files).length > 0 ? (
          Object.entries(files).map(([key, value]) => {
            if (value.path === RootDir) {
              return (
                  <FileComponent key={key} file={value} />
              );
            }
          })
        ) : (
          <div>No file available</div>
        )}
    </>
  );
}
