import { create } from "zustand";

interface ProjectContainerIP {
  ip: string | null;
  setProjectContainerIP: (ip: string) => void;
}

export const useProjectContainerIP = create<ProjectContainerIP>((set) => ({
  ip: null,
  setProjectContainerIP: (ip) => set({ ip: ip}),
}));
