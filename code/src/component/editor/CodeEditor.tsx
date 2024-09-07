import { FC, useContext } from "react";
import { Editor } from "@monaco-editor/react";
import { useCurrentFile } from "../../state/currentFile";
import { useProjectFiles } from "../../state/projectFilesState";
import { sendFileChanges } from "../socket/socketHandler";
import SocketProvider from "../../context/socketContextProvider";

const CodeEditor: FC = () => {
  const { file } = useCurrentFile();
  const { updateFile } = useProjectFiles();
  const { ws } = useContext(SocketProvider);

  //TODO: Also send changes to sockekt server.
  function onChangeHandler(change: string | undefined, e: any) {
    const path = file?.path + "/" + file?.name;
    updateFile(path, { content: change });
    if (file === null) return;
    sendFileChanges(ws, file, change);
  }

  function getFileMode() {
    const name = file?.name;
    if (name === undefined) {
      return;
    }

    const splitName = name.split(".");
    const extension = splitName[splitName.length - 1];
    switch (extension) {
      case "js":
        return "javascript";
      case "py":
        return "python";
      case "java":
        return "java";
      case "xml":
        return "xml";
      case "rb":
        return "ruby";
      case "sass":
        return "sass";
      case "md":
        return "markdown";
      case "sql":
        return "mysql";
      case "json":
        return "json";
      case "html":
        return "html";
      case "hbs":
        return "handlebars";
      case "handlebars":
        return "handlebars";
      case "go":
        return "go";
      case "cs":
        return "csharp";
      case "litcoffee":
        return "coffee";
      case "css":
        return "css";
      default:
        return "";
    }
  }

  return (
    <>
      <div>{file?.name}</div>
      <Editor
        value={file?.content}
        height="100%"
        width="100%"
        theme="vs-dark"
        language={getFileMode()}
        defaultValue="console.log(`hello world`)"
        onChange={onChangeHandler}
      />
    </>
  );
};

export default CodeEditor;
