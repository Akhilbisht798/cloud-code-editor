import { FC } from "react";
import { useProjectFiles } from "../../state/projectFilesState";

interface FileExplorereInterface {
  ws: WebSocket | null;
}

const FileExplorer: FC = () => {
  const { files } = useProjectFiles();
  return (
    <>
      {console.log(files)}
      <div>Hello world</div>
    </>
  );
};

export default FileExplorer;
