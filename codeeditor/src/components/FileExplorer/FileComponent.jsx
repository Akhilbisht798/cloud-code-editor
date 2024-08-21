import { useState } from "react";
import { requestFiles } from "../socket/socketHandler";
import useFile from "../../state/file";
import useFiles from "../../state/files";

export default function FileComponent({ file }) {
  const [childFiles, setChildFiles] = useState([]);
  const [isExpanded, setIsExpanded] = useState(false);
  const { setFile } = useFile()
  const { files } = useFiles()

  function getChildFiles() {
    const searchValue = file.path + "/" + file.name;

    const childFilesArray = [];

    for (const key in files) {
      if (
        files.hasOwnProperty(key) &&
        files[key].path === searchValue
      ) {
        childFilesArray.push(files[key]);
      }
    }

    setChildFiles(childFilesArray);
  }

  function onClickHandler(e) {
    const filePath = e.target.id;

    if (file.isDir && file.files === undefined) {
      console.log("requesting files");
      requestFiles(filePath);
      if (!isExpanded) {
        getChildFiles();
        setIsExpanded(true);
      }
    }
    if (file.isDir && file.files) {
      if (isExpanded) {
        // delete all files below it.
        return;
      }
      // else expand the files.
      getChildFiles();
    } else {
      //File handling.
      setFile(file)
      console.log("Is file");
    }
  }
  return (
    <>
      <div
        id={file.path + "/" + file.name}
        data-parent={file.path}
        onClick={onClickHandler}
        className="cursor-pointer"
      >
        {file.name}
      </div>
      {file.isDir ? (
        <ul>
          <li>
            {childFiles.map((childFile) => (
              <FileComponent
                file={childFile}
                key={childFile.path + "/" + childFile.name}
              />
            ))}
          </li>
        </ul>
      ) : null}
    </>
  );
}
