export interface File {
  path: string;
  name: string;
  isDir: boolean;
  content?: string;
  hasFiles?: boolean;
}
