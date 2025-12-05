export default class ShoppingList implements ProtocolEntity {
    id: string;
    name: string;
    items: string[];
    createdAt: Date;

    constructor(id: string, name: string, items: string[] = []) {
        this.id = id;
        this.name = name;
        this.items = items;
        this.createdAt = new Date();
    }

    serialize(): Uint8Array {
        // Serialization logic here
        return new Uint8Array();
    }
}