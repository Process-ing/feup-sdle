interface ProtocolSocket {
    sendTo(entity: ProtocolEntity, recipientId: string): void;  // TODO: Change to use proper node identification
}