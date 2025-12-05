import { type DBSchema, type IDBPDatabase, openDB } from "idb";
import { ShoppingItem, ShoppingList } from "@/types";

const DB_NAME = "shopping-lists-db";
const DB_VERSION = 2;
const LIST_STORE_NAME = "shopping-lists";
const ITEM_STORE_NAME = "shopping-items";

interface ShoppingListDB extends DBSchema {
	[LIST_STORE_NAME]: {
		key: string;
		value: {
			id: string;
			name: string;
			itemIds: string[];
		};
	};
	[ITEM_STORE_NAME]: {
		key: string;
		value: {
			id: string;
			name: string;
			totalQuantity: number;
			acquiredQuantity: number;
		};
	};
}

class DB {
	private dbPromise: Promise<IDBPDatabase<ShoppingListDB>> | null = null;

	constructor() {
		if (typeof window !== "undefined") {
			this.dbPromise = this.initDB();
		}
	}

	private async initDB(): Promise<IDBPDatabase<ShoppingListDB>> {
		return openDB<ShoppingListDB>(DB_NAME, DB_VERSION, {
			upgrade(db) {
				if (!db.objectStoreNames.contains(LIST_STORE_NAME)) {
					db.createObjectStore(LIST_STORE_NAME, { keyPath: "id" });
				}
				if (!db.objectStoreNames.contains(ITEM_STORE_NAME)) {
					db.createObjectStore(ITEM_STORE_NAME, { keyPath: "id" });
				}
			},
		});
	}

	private async getDB(): Promise<IDBPDatabase<ShoppingListDB>> {
		if (!this.dbPromise) {
			throw new Error("Database not initialized");
		}
		return this.dbPromise;
	}

	async getAllLists(): Promise<ShoppingList[]> {
		const db = await this.getDB();
		const lists = await db.getAll(LIST_STORE_NAME);

		return lists.map(
			(list) => new ShoppingList(list.id, list.name, list.itemIds)
		);
	}

	async getList(id: string): Promise<ShoppingList | undefined> {
		const db = await this.getDB();
		const list = await db.get(LIST_STORE_NAME, id);
		if (!list) return undefined;

		return new ShoppingList(list.id, list.name, list.itemIds);
	}

	async createList(name: string): Promise<ShoppingList> {
		const db = await this.getDB();

		const list = new ShoppingList(
			Math.random().toString(36).substring(2, 9),
			name,
			[],
		);

		await db.add(LIST_STORE_NAME, list);
		return list;
	}

	async updateList(list: ShoppingList): Promise<void> {
		const db = await this.getDB();
		await db.put(LIST_STORE_NAME, list);
	}

	async deleteList(id: string): Promise<void> {
		const db = await this.getDB();
		await db.delete(LIST_STORE_NAME, id);
	}

	async getItem(id: string): Promise<ShoppingItem | undefined> {
		const db = await this.getDB();

		const item = await db.get(ITEM_STORE_NAME, id);
		if (!item) return undefined;

		return new ShoppingItem(
			item.id,
			item.name,
			item.totalQuantity,
			item.acquiredQuantity
		);
	}

	async createItem(name: string, quantity: number): Promise<ShoppingItem> {
		const db = await this.getDB();

		const item = new ShoppingItem(
			Math.random().toString(36).substring(2, 9),
			name,
			quantity,
			0
		);

		db.add(ITEM_STORE_NAME, item);
		return item;
	}

	async updateItem(item: ShoppingItem): Promise<void> {
		const db = await this.getDB();
		await db.put(ITEM_STORE_NAME, item);
	}

	async deleteItem(id: string): Promise<void> {
		const db = await this.getDB();
		await db.delete(ITEM_STORE_NAME, id);
	}
}

export const db = new DB();
