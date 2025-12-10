import { ShoppingList } from "@/types";
import { ClientRequest, ErrorCode, ServerResponse, ShoppingList as ShoppingListProto } from "../proto/client";
import ProtocolRequest from "./protocol-entity";
import ProtocolSocket from "./protocol-socket";

class WebProtocolSocket implements ProtocolSocket {
    private socket: WebSocket;
    private isOpen: Promise<void>;
    private onShoppingList: (list: ShoppingList) => void = () => { };
    private errorHandlers: Map<string, (error: ErrorCode) => void> = new Map();

    constructor(ws: WebSocket) {
        this.socket = ws;

        this.isOpen = new Promise((resolve, reject) => {
            this.socket.onopen = (event: Event) => {
                console.log("WebSocket connection established on " + this.socket.url);
                resolve();
            };

            this.socket.onerror = (event: Event) => {
                console.error("WebSocket error on " + this.socket.url, event);
                reject(new Error("WebSocket connection failed"));
            };
        });

        // Create a Promise that resolves when the WebSocket is open
        this.socket.onmessage = this.onMessage.bind(this);
        this.socket.onclose = this.onClose.bind(this);
    }

    private async onMessage(event: MessageEvent): Promise<void> {
        const buffer = await (event.data as Blob).arrayBuffer();
        const response = ServerResponse.decode(new Uint8Array(buffer));

        console.log("Received message via WebSocket:", response);

        const errorHandler = this.errorHandlers.get(response.messageId);
        this.errorHandlers.delete(response.messageId);

        switch (response.responseType) {
            case "shoppingList":
                const shoppingList = ShoppingList.fromProto(response.shoppingList as ShoppingListProto)
                this.onShoppingList(shoppingList);
                break;

            case "error":
                if (errorHandler) {
                    errorHandler(response.error!);
                }
                break;
        }
    }

    private onClose(event: CloseEvent): void {
        console.log("WebSocket connection closed on " + this.socket.url, event);
    }

    async send(request: ProtocolRequest, onError: (error: ErrorCode) => void): Promise<void> {
        // Wait for the WebSocket to be open
        await this.isOpen;

        if (this.socket.readyState !== WebSocket.OPEN) {
            console.error("WebSocket is not open. Current state:", this.socket.readyState);
            return;
        }

        const clientReq = request.toClientRequest();
        const buffer = ClientRequest.encode(clientReq).finish();

        console.log("Sending entity via WebSocket:", buffer);

        this.errorHandlers.set(clientReq.messageId!, onError);
        this.socket.send(buffer);
    }

    setOnShoppingListCallback(callback: (list: ShoppingList) => void): void {
        this.onShoppingList = callback;
    }
}

export default WebProtocolSocket;
