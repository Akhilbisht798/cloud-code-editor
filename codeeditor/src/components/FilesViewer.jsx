import { useState } from "react";
import { useEffect } from "react";
import { requestFiles } from "./socket/socketHandler";

export default function FilesViewer() {
  const [files, setFiles] = useState(null);
  const [expandedFolders, setExpandedFolders] = useState([]);

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

  function onClickHandler(e) {
    const dir = e.currentTarget.dataset.path;
    const name = e.currentTarget.dataset.name;

    const filePath = dir + "/" + name;

    //TODO: don't think this will traverse through entire files.
    for (let i = 0; i < files.length; i++) {
      let obj = files[i];
      const objFilePath = obj.path + "/" + obj.name;

      if (filePath === objFilePath) {
        if (obj.isDir && obj.files === undefined) {
          requestFiles(objFilePath);
        } else if (obj.isDir) {
          //Check if directory is expanded or not
          // and do the opposite of that
          // get all the files and add divs below it.
        } else {
          // this is file show its content on editor.
        }
        break;
      }
    }
  }

  return (
    <>
      {files && files.length > 0 ? (
        files.map((file, index) => (
          <div
            key={index}
            className=""
            data-path={file.path}
            data-name={file.name}
            onClick={onClickHandler}
          >
            {file.name}
          </div>
        ))
      ) : (
        <div>No file available</div>
      )}
    </>
  );
}
