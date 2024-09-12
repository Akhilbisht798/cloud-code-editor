import { create } from "zustand";

interface UserAuthInterface {
  email: string;
  setEmail: (email: string) => void;
}

export const useAuth = create<UserAuthInterface>((set) => ({
  email: "",
  setEmail: (email) => set({ email: email }),
}));
