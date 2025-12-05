import { type DBSchema, type IDBPDatabase, openDB } from "idb";
import type { ShoppingList } from "@/types";

const DB_NAME = "shopping-lists-db";
const DB_VERSION = 1;
const STORE_NAME = "shopping-lists";

interface ShoppingListDB extends DBSchema {
	[STORE_NAME]: {
		key: string;
		value: ShoppingList;
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
				if (!db.objectStoreNames.contains(STORE_NAME)) {
					db.createObjectStore(STORE_NAME, { keyPath: "id" });
				}
			},
		});
	}

	async getAllLists(): Promise<ShoppingList[]> {
		if (!this.dbPromise) return [];
		const db = await this.dbPromise;
		return db.getAll(STORE_NAME);
	}

	async getList(id: string): Promise<ShoppingList | undefined> {
		if (!this.dbPromise) return undefined;
		const db = await this.dbPromise;
		return db.get(STORE_NAME, id);
	}

	async createList(list: ShoppingList): Promise<void> {
		if (!this.dbPromise) return;
		const db = await this.dbPromise;
		await db.add(STORE_NAME, list);
	}

	async updateList(list: ShoppingList): Promise<void> {
		if (!this.dbPromise) return;
		const db = await this.dbPromise;
		await db.put(STORE_NAME, list);
	}

	async deleteList(id: string): Promise<void> {
		if (!this.dbPromise) return;
		const db = await this.dbPromise;
		await db.delete(STORE_NAME, id);
	}
}

export const db = new DB();
