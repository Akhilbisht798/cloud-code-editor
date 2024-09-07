import { FC } from "react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@radix-ui/react-dropdown-menu";
import { Button } from "@/components/ui/button";
import { File } from "@/interface";
import { deleteFileOrDir, newFileOrDirChanges } from "../socket/socketHandler";
import { useProjectFiles } from "@/state/projectFilesState";

interface FileOptionsInterFace {
  file: File;
  ws: WebSocket | null | undefined;
}

const FileOptions: FC<FileOptionsInterFace> = ({ file, ws }) => {
  const { path, name, isDir } = file;
  const { files, deleteFile } = useProjectFiles();

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
  }

  function newFolderHandler(e: React.MouseEvent<HTMLButtonElement>) {
    e.stopPropagation();
    const dirName = window.prompt("Enter folder name: ");
    if (dirName === null) {
      window.alert("folder name should be provided");
      return;
    }
    const dirPath = path + "/" + name;
    const newFolder: File = {
      path: dirPath,
      isDir: true,
      name: dirName,
      content: "",
    };
    newFileOrDirChanges(ws, newFolder);
  }

  function deleteFileOrFolder(e: React.MouseEvent<HTMLButtonElement>) {
    e.stopPropagation();
    const p = path + "/" + name;
    if (!isDir) {
      deleteFileOrDir(ws, file);
      deleteFile(p);
      return;
    }
    deleteFileOrDir(ws, file);
    deleteFile(p);

    const parent = p + "/";
    Object.keys(files).forEach((filePath) => {
      if (filePath.startsWith(parent)) {
        deleteFile(filePath);
      }
    });
  }

  if (isDir) {
    return (
      <>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button>⋮</button>
          </DropdownMenuTrigger>
          <DropdownMenuContent className="w-56 bg-[#252526] border-[#3c3c3c] text-[#cccccc]" align="end">
            <DropdownMenuItem className="hover:bg-[#2a2d2e]" >
              <Button
                className="w-56"
                variant="ghost"
                id={path + "/" + name}
                onClick={newFileHandler}
              >
                Create a File
              </Button>
            </DropdownMenuItem>
            <DropdownMenuItem className="hover:bg-[#2a2d2e]">
              <Button
                id={path + "/" + name}
                className="w-56"
                variant="ghost"
                onClick={newFolderHandler}
              >
                Create a Folder
              </Button>
            </DropdownMenuItem>
            <DropdownMenuItem className="hover:bg-[#2a2d2e]">
              <Button
                id={path + "/" + name}
                className="w-56"
                variant="ghost"
                onClick={deleteFileOrFolder}
              >
                delete
              </Button>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </>
    );
  }

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <button>⋮</button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-56 bg-[#252526] border-[#3c3c3c] text-[#cccccc]" align="end">
          <DropdownMenuItem className="hover:bg-[#2a2d2e]">
            <Button
              id={path + "/" + name}
              className="w-56"
              variant="ghost"
              onClick={deleteFileOrFolder}
            >
              delete
            </Button>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </>
  );
};

export default FileOptions;
