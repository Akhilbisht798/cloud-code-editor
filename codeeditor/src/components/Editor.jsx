import AceEditor from "react-ace";

import "ace-builds/src-noconflict/mode-java";
import "ace-builds/src-noconflict/theme-github";
import "ace-builds/src-noconflict/ext-language_tools";

export default function Editor() {
  function onChange(change) {
    console.log("change:", change);
  }

  return (
    <AceEditor
      mode="java"
      value="public static void main"
      fontSize="16"
      theme="github"
      onChange={onChange}
      name="editor"
      editorProps={{ $blockScrolling: true }}
    />
  );
}
