import { createContext } from "react";

export interface SocketContextType {
  ws: WebSocket | null;
}

const SocketProvider = createContext<Partial<SocketContextType>>({});

export default SocketProvider;
