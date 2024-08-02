import AceEditor from "react-ace";

import "ace-builds/src-noconflict/mode-java";
import "ace-builds/src-noconflict/theme-github";
import "ace-builds/src-noconflict/ext-language_tools";
import { useState } from "react";
import { fileChanged } from "./socket/socketHandler";

export default function Editor() {
  //context apis
  const [file, setFile] = useState(null);

  function onChange(change) {
    console.log("change:", change);
  }

  return (
    <AceEditor
      mode="java"
      value="hello"
      fontSize="16"
      theme="github"
      onChange={onChange}
      name="editor"
      editorProps={{ $blockScrolling: true }}
    />
  );
}
