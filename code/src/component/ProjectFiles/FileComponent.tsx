import React, { FC, useState } from "react";
import { File } from "../../interface";
import { useProjectFiles } from "../../state/projectFilesState";
import { requestFiles } from "../socket/socketHandler";

const FileComponent: FC<File> = (props) => {
  const { path, name, isDir } = props;
  const [expand, setExpand] = useState(false);
  const [childFiles, setChildFiles] = useState<File[]>([]);
  const { files, updateFile } = useProjectFiles();

  function getChildFiles(): File[] {
    const searchValue = path + "/" + name;
    const childFilesArray: File[] = [];

    for (const key in files) {
      if (
        Object.prototype.hasOwnProperty.call(files, key) &&
        files[key].path === searchValue
      ) {
        childFilesArray.push(files[key]);
      }
    }

    return childFilesArray;
  }

  function onClickHandler(e: React.MouseEvent<HTMLDivElement>) {
    const target = e.target as HTMLDivElement;
    const filePath = target.id;
    const file = files[filePath];
    if (isDir && file.hasFiles) {
      setExpand(!expand);
    } else if (isDir) {
      console.log("getting files");
      requestFiles(ws, filePath);
    }
  }

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
