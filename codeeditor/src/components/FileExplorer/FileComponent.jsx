import { useState } from "react";
import { requestFiles } from "../socket/socketHandler";
import useFile from "../../state/file";

export default function FileComponent({ file, dictFiles }) {
  const [childFiles, setChildFiles] = useState([]);
  const [isExpanded, setIsExpanded] = useState(false);
  const { setFile } = useFile()

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
    <div>
      <p
        id={file.path + "/" + file.name}
        data-parent={file.path}
        onClick={onClickHandler}
        className="cursor-pointer"
      >
        {file.name}
      </p>
      {file.isDir ? (
        <ul>
          <li>
            {childFiles.map((childFile) => (
              <FileComponent
                file={childFile}
                dictFiles={dictFiles}
                key={childFile.path + "/" + childFile.name}
              />
            ))}
          </li>
        </ul>
      ) : null}
    </div>
  );
}
