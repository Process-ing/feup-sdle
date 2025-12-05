import { ShoppingItem, ShoppingList } from "@/types";

interface ProtocolListener<T> {
    on(shoppingList: ShoppingList): T;
    on(shoppingItem: ShoppingItem): T;
}