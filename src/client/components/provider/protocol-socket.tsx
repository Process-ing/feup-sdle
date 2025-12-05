'use client';

import NullProtocolSocket from "@/lib/protocol/null-protocol-socket";
import WebProtocolSocket from "@/lib/protocol/web-protocol-socket";
import { createContext, useContext, useEffect, useState } from "react";

const ProtocolSocketContext = createContext<ProtocolSocket>(new NullProtocolSocket());

export const WebProtocolSocketProvider = ({ children }: { children: React.ReactNode }) => {
  const [socket, setSocket] = useState<ProtocolSocket>(new NullProtocolSocket());

  useEffect(() => {
    const websocketUrl = "ws://localhost:8080/ws";
    const websocket = new WebSocket(websocketUrl);

    const protocolSocket = new WebProtocolSocket(websocket);

    setSocket(protocolSocket);
  }, []);

  return (
    <ProtocolSocketContext.Provider value={socket}>
      {children}
    </ProtocolSocketContext.Provider>
  )
}

export const useProtocolSocket = (): ProtocolSocket => {
  const socket = useContext(ProtocolSocketContext);
  return socket;
}
