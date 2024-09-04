import { create } from "zustand";
import { File } from "../interface";

interface CurrentFile {
  file: File | null;
  setCurrentFile: (file: File) => void;
}

export const useCurrentFile = create<CurrentFile>((set) => ({
  file: null,
  setCurrentFile: (newFile) => set({ file: newFile }),
}));
