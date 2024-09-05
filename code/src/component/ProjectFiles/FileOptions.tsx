import { FC } from "react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@radix-ui/react-dropdown-menu";
import { Button } from "@/components/ui/button";
import { File } from "@/interface";
import { newFileOrDirChanges } from "../socket/socketHandler";

interface FileOptionsInterFace {
  file: File;
  ws: WebSocket | null | undefined;
}

const FileOptions: FC<FileOptionsInterFace> = ({ file, ws }) => {
  const { path, name } = file;

  function newFileHandler(e: React.MouseEvent<HTMLButtonElement>) {
    e.stopPropagation();
    const filename = window.prompt("Enter file name: ");
    if (filename === null) {
      window.alert("Filename should be provided");
      return;
    }
    const dirPath = path + "/" + name;
    const newFile: File = {
      path: dirPath,
      isDir: false,
      name: filename,
      content: "",
    };
    newFileOrDirChanges(ws, newFile);
    //setFiles({
    //  [dirPath + "/" + filename] : newFile
    //})
  }

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <button>⋮</button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-56">
          <DropdownMenuItem>
            <Button
              className="w-56"
              variant="outline"
              id={path + "/" + name}
              onClick={newFileHandler}
            >
              new file
            </Button>
          </DropdownMenuItem>
          <DropdownMenuItem>
            <Button id={path + "/" + name} className="w-56" variant="outline">
              new folder
            </Button>
          </DropdownMenuItem>
          <DropdownMenuItem>
            <Button id={path + "/" + name} className="w-56" variant="outline">
              delete
            </Button>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </>
  );
};

export default FileOptions;
