import { ClientRequest } from "../proto/client";
import ProtocolRequest from "./protocol-entity";

export default class GetShoppingListRequest implements ProtocolRequest {
    private listId: string;

    constructor(listId: string) {
        this.listId = listId;
    }

    public getListId(): string {
        return this.listId;
    }

    public toClientRequest(): ClientRequest {
        return ClientRequest.create({
            messageId: crypto.randomUUID(),
            getShoppingList: {
                id: this.listId
            }
        })

    }
}
