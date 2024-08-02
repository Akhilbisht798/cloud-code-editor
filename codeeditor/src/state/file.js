import { create } from  'zustand'

const useFile = create((set) => ({
    file: null,
    setFile: (file) => set({ file: file})
}));

export default useFile;