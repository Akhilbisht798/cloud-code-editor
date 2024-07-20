import { useState } from "react";
import { requestFiles } from "../socket/socketHandler";

export default function FileComponent({ file, dictFiles, depth = 0 }) {
  const [childFiles, setChildFiles] = useState([]);

  function getChildFiles() {
    const searchValue = file.path + "/" + file.name;

    const childFilesArray = [];

    for (const key in dictFiles) {
      if (
        dictFiles.hasOwnProperty(key) &&
        dictFiles[key].path === searchValue
      ) {
        childFilesArray.push(dictFiles[key]);
      }
    }

    setChildFiles(childFilesArray);
  }

  function onClickHandler(e) {
    const filePath = e.target.id;

    if (file.isDir && file.files === undefined) {
      console.log("requesting files");
      requestFiles(filePath);
    }
    if (file.isDir && file.files) {
      getChildFiles();
    } else {
      console.log("Is file");
    }
  }
  return (
    <>
      <div
        id={file.path + "/" + file.name}
        data-parent={file.path}
        onClick={onClickHandler}
        style={{ paddingLeft: `${depth * "1rem"}` }}
        className="cursor-pointer"
      >
        {file.name}
      </div>
      <div>
        {childFiles.map((childFile) => (
          <FileComponent
            file={childFile}
            dictFiles={dictFiles}
            key={childFile.path + "/" + childFile.name}
            depth={depth + 1}
          />
        ))}
      </div>
    </>
  );
}
