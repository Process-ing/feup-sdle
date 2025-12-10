import { ErrorCode } from "../proto/client";
import ProtocolRequest from "./protocol-entity";

export default interface ProtocolSocket {
    send(entity: ProtocolRequest, onError: (error: ErrorCode) => void): void;
}