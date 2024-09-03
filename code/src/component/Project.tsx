import { FC } from "react";
import { File } from "../interface";
import FileComponent from "./ProjectFiles/FileComponent";

const Project: FC = () => {
  const rootFile: File = {
    path: "..",
    name: "client",
    isDir: true,
  };

  return (
    <div>
      <FileComponent {...rootFile} />
    </div>
  );
};

export default Project;
