import ProtocolEntity from "./protocol-entity";

export default interface ProtocolSocket {
    send(entity: ProtocolEntity): void;
}