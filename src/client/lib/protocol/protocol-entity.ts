import { IClientRequest } from "../proto/client";

export default interface ProtocolEntity {
    toClientRequest(): IClientRequest;
}
