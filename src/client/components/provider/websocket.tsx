'use client';

import WebProtocolSocket from "@/lib/protocol/web-protocol-socket";
import { createContext, useContext, useEffect, useRef, useState } from "react";

export const WebProtocolSocketContext = createContext<WebProtocolSocket | null>(null);

export const WebProtocolSocketProvider = ({ children }: { children: React.ReactNode }) => {
  const [socket, setSocket] = useState<WebProtocolSocket | null>(null);

  useEffect(() => {
    const websocketUrl = "ws://localhost:8080/ws";
    const websocket = new WebSocket(websocketUrl);

    const protocolSocket = new WebProtocolSocket(websocket);

    setSocket(protocolSocket);
  }, []);

  return (
    <WebProtocolSocketContext.Provider value={socket}>
      {children}
    </WebProtocolSocketContext.Provider>
  )
}

export const useWebProtocolSocket = (): WebProtocolSocket | null => {
  const socket = useContext(WebProtocolSocketContext);
  return socket;
}
