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
        // Serialization logic here
        return new Uint8Array();
    }
}