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

  return (
    <AceEditor
      mode="java"
      value={file?.content}
      fontSize="16"
      // height="100%"
      // width="100%"
      //style={{ position: "relative", height: "100%", width: "100%" }}
      theme="github"
      onChange={onChange}
      name="editor"
      editorProps={{ $blockScrolling: true }}
    />
  );
}
