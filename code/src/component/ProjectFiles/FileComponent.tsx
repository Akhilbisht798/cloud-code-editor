import React, { FC, useEffect, useState } from "react";
import { File } from "../../interface";
import { useProjectFiles } from "../../state/projectFilesState";
import { requestFiles } from "../socket/socketHandler";
import { useContext } from "react";
import { useCurrentFile } from "../../state/currentFile";
import SocketProvider from "../../context/socketContextProvider";
import FileOptions from "./FileOptions";

const FileComponent: FC<File> = (props) => {
  const { ws } = useContext(SocketProvider);
  const { path, name, isDir } = props;
  const [expand, setExpand] = useState(false);
  const [childFiles, setChildFiles] = useState<File[]>([]);
  const { files, updateFile } = useProjectFiles();
  const { setCurrentFile } = useCurrentFile();

  function getChildFiles() {
    const searchValue = path + "/" + name;

    const childFilesArray: File[] = [];

    for (const key in files) {
      if (files[key].path === searchValue) {
        childFilesArray.push(files[key]);
      }
    }

    setChildFiles(childFilesArray);
  }

  function onClickHandler(e: React.MouseEvent<HTMLDivElement>) {
    const target = e.target as HTMLDivElement;
    const filePath = target.id;
    const file = files[filePath];
    if (isDir && file.hasFiles) {
      setExpand(!expand);
    } else if (isDir) {
      requestFiles(ws, filePath);
      updateFile(filePath, { hasFiles: true });
    } else if (!isDir) {
      setCurrentFile({ path, isDir, name, content: file.content });
    }
  }

  useEffect(() => {
    if (isDir) {
      getChildFiles();
    }
  }, [files]);

  if (isDir) {
    return (
      <div className="text-[#cccccc] text-sm">
        <div
          id={path + "/" + name}
          onClick={onClickHandler}
          className="flex justify-between cursor-pointer"
        >
          <span id={path + "/" + name}>🖿 {name}</span>
          <FileOptions file={{ path, name, isDir }} ws={ws} />
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
      </div>
    );
  }

  return (
    <div
      id={path + "/" + name}
      data-parent={path}
      onClick={onClickHandler}
      className="flex justify-between cursor-pointer"
    >
      <span id={path + "/" + name}>📄 {name}</span>
      <FileOptions file={{ path, name, isDir }} ws={ws} />
    </div>
  );
};

export default FileComponent;
