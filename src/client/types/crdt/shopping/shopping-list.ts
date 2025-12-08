import DotContext from "../generic/dot-context";
import ORMap from "../generic/ormap";
import ShoppingItem from "./shopping-item";

export default class ShoppingList {
    private replicaId: string;
    private listId: string;
    private name: string;
    private items: ORMap<string, ShoppingItem>
    private dotContext: DotContext;

    constructor(replicaId: string, listId: string, name: string) {
        this.replicaId = replicaId;
        this.listId = listId;
        this.name = name;
        this.dotContext = new DotContext();

        this.items = new ORMap<string, ShoppingItem>(
            this.replicaId,
            (replicaId: string) => new ShoppingItem(replicaId, "", ""),
            this.dotContext
        );
    }

    public getReplicaId(): string {
        return this.replicaId;
    }

    public getListId(): string {
        return this.listId;
    }

    public getName(): string {
        return this.name;
    }

    public getContext(): DotContext {
        return this.dotContext;
    }

    public setContext(context: DotContext): void {
        this.dotContext = context;
        this.items.setContext(context);
    }

    public getItems(): ShoppingItem[] {
        return Array.from(this.items.keys())
            .map((itemId: string) => this.items.get(itemId).clone())
            .filter((item: ShoppingItem) => !item.isNull());
    }

    public getItem(itemId: string): ShoppingItem | undefined {
        const item = this.items.get(itemId);
        if (item.isNull()) {
            return undefined;
        }
        return item.clone();
    }

    public putItem(itemId: string, name: string, quantity: number, acquired: number): ShoppingList {
        const delta = new ShoppingList(this.replicaId, this.listId, this.name);

        const itemsDelta = this.items.apply(itemId, (item: ShoppingItem) => {
            item.setItemId(itemId);
            item.setName(name);

            const itemDelta = item.incQuantity(quantity);
            itemDelta.join(item.incAcquired(acquired));

            return itemDelta;
        });

        delta.items = itemsDelta;
        delta.setContext(itemsDelta.getContext());

        return delta;
    }

    public removeItem(itemId: string): ShoppingList {
        const delta = new ShoppingList(this.replicaId, this.listId, this.name);

        const itemsDelta = this.items.remove(itemId);
        delta.items = itemsDelta;
        delta.setContext(itemsDelta.getContext());

        return delta;
    }

    public join(other: ShoppingList): void {
        this.dotContext.join(other.getContext());
        this.items.join(other.items);
    }

    public clone(): ShoppingList {
        const clone = new ShoppingList(this.replicaId, this.listId, this.name);
        clone.setContext(this.dotContext);
        clone.items = this.items.clone();
        return clone;
    }
}