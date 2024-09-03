import { create } from "zustand";
import { File } from "../interface";

interface ProjectFiles {
  files: { [key: string]: File };
  setFiles: (file: { [key: string]: File }) => void;
  updateFile: (path: string, updates: Partial<File>) => void;
}

export const useProjectFiles = create<ProjectFiles>((set) => ({
  files: {},
  setFiles: (newfiles) =>
    set((state) => ({
      files: { ...state.files, ...newfiles },
    })),
  updateFile: (path, updates) =>
    set((state) => ({
      ...state.files,
      [path]: {
        ...state.files[path],
        ...updates,
      },
    })),
}));
