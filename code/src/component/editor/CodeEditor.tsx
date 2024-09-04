import { FC } from "react";
import { Editor } from "@monaco-editor/react";
import { useCurrentFile } from "../../state/currentFile";
import { useProjectFiles } from "../../state/projectFilesState";

const CodeEditor: FC = () => {
  const { file } = useCurrentFile();
  const { updateFile } = useProjectFiles();

  function onChangeHandler(change: string | undefined, e: any) {
    const path = file?.path + "/" + file?.name;
    updateFile(path, { content: change });
  }

  console.log(file);
  return (
    <>
      <div className=" w-full h-full font-sans ">
        <Editor
          value={file?.content}
          height="75vh"
          width="75vw"
          theme="vs-dark"
          defaultLanguage="javascript"
          defaultValue="console.log(`hello world`)"
          onChange={onChangeHandler}
        />
      </div>
    </>
  );
};

export default CodeEditor;
