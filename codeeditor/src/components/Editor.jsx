import AceEditor from "react-ace";

import "ace-builds/src-noconflict/mode-java";
import "ace-builds/src-noconflict/theme-github";
import "ace-builds/src-noconflict/ext-language_tools";
import { fileChanged } from "./socket/socketHandler";
import useFile from "../state/file";

export default function Editor() {
  const { file } = useFile()

  function onChange(change) {
    console.log("change:", change);
  }

  function getFileMode() {
    const fileName = file?.name;
    if (fileName === undefined) {
      return ""
    }
    const splitedArray = fileName.split(".");

    const extension = splitedArray[splitedArray.length - 1];

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
        return "golang";
      case "cs":
        return "csharp";
      case "litcoffee":
        return "coffee";
      case "css":
        return "css";
      default:
        return "";
    }
  };

  return (
    <div className=" w-full h-full font-sans ">
      <AceEditor
        mode={getFileMode(file)}
        value={file?.content}
        fontSize="16"
        width="100%"
        height="100%"
        style={{fontSize: "16px"}}
        theme="github"
        onChange={onChange}
        name="editor"
        editorProps={{ $blockScrolling: true }}
      />
    </div>
  );
}
