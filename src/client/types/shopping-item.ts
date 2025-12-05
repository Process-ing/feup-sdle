import { ShoppingItem as ShoppingItemProto } from "@/lib/proto/global";

export default class ShoppingItem implements ProtocolEntity {
    id: string;
    name: string;
    totalQuantity: number;
    acquiredQuantity: number;

    constructor(id: string, name: string, totalQuantity: number, acquiredQuantity: number) {
        this.id = id;
        this.name = name;
        this.totalQuantity = totalQuantity;
        this.acquiredQuantity = acquiredQuantity;
    }

    serialize(): Uint8Array {
        const proto = new ShoppingItemProto({
            id: this.id,
            name: this.name,
            totalQuantity: this.totalQuantity,
            acquiredQuantity: this.acquiredQuantity,
        });

        return ShoppingItemProto.encode(proto).finish();
    }
}