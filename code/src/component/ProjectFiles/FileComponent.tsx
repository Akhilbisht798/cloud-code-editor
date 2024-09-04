import React, { FC, useEffect, useState } from "react";
import { File } from "../../interface";
import { useProjectFiles } from "../../state/projectFilesState";
import { requestFiles } from "../socket/socketHandler";
import { useContext } from "react";
import SocketProvider from "../../context/socketContextProvider";

const FileComponent: FC<File> = (props) => {
  const { ws }= useContext(SocketProvider);
  const { path, name, isDir } = props;
  const [expand, setExpand] = useState(false);
  const [childFiles, setChildFiles] = useState<File[]>([]);
  const { files, updateFile } = useProjectFiles();

  function getChildFiles() {
    const searchValue = path + "/" + name;

    const childFilesArray: File[] = [];

    for (const key in files) {
      if (files[key].path === searchValue) {
        childFilesArray.push(files[key]);
      }
    }

    for (let i = 0; i < childFilesArray.length; i++) {
      console.log(childFilesArray[i]);
    }

    setChildFiles(childFilesArray);
  }

  function onClickHandler(e: React.MouseEvent<HTMLDivElement>) {
    const target = e.target as HTMLDivElement;
    const filePath = target.id;
    const file = files[filePath];
    if (isDir && file.hasFiles) {
      // console.log("Child Files: ", childFiles);
      setExpand(!expand);
    } else if (isDir) {
      requestFiles(ws, filePath);
      updateFile(filePath, { hasFiles: true });
    }
  }

  useEffect(() => {
    if (isDir && !childFiles.length) {
      getChildFiles();
    }
  }, [files]);
  // console.log("All Files: ", files);

  if (isDir) {
    return (
      <>
        <div id={path + "/" + name} data-parent={path} onClick={onClickHandler}>
          <span id={path + "/" + name}>🖿 {name}</span>
        </div>
        <div
          style={{ display: expand ? "block" : "none", paddingLeft: "22px" }}
        >
          {childFiles.map((file) => (
            <FileComponent
              key={file.path + "/" + file.name}
              name={file.name}
              path={file.path}
              isDir={file.isDir}
              hasFiles={file?.hasFiles}
            />
          ))}
        </div>
      </>
    );
  }

  return (
    <div id={path + "/" + name} data-parent={path} onClick={onClickHandler}>
      <span>📄 {name}</span>
    </div>
  );
};

export default FileComponent;
