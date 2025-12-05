import { ShoppingList as ShoppingListProto } from "@/lib/proto/global";

export default class ShoppingList implements ProtocolEntity {
    id: string;
    name: string;
    itemIds: string[];

    constructor(id: string, name: string, itemIds: string[] = []) {
        this.id = id;
        this.name = name;
        this.itemIds = itemIds;
    }

    addItem(itemId: string): void {
        if (!this.itemIds.includes(itemId)) {
            this.itemIds.push(itemId);
        }
    }

    removeItem(itemId: string): void {
        this.itemIds = this.itemIds.filter((id) => id !== itemId);
    }

    serialize(): Uint8Array {
        const proto = new ShoppingListProto({
            id: this.id,
            name: this.name,
            itemIds: this.itemIds,
        });

        return ShoppingListProto.encode(proto).finish();
    }
}