import { create } from "zustand";

interface CurrentProject {
  rootDir: string | null;
  setRootDir: (file: string) => void;
}

export const useCurrentProject = create<CurrentProject>((set) => ({
  rootDir: null,
  setRootDir: (root) => set({ rootDir: root }),
}));
