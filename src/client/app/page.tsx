"use client";

import { useRouter } from "next/navigation";
import { WebSocketContext } from "@/components/provider/websocket";
import { ShoppingListHome } from "@/components/shopping-list-home";
import { useContext } from "react";

export default function HomePage() {
  const router = useRouter();

	return (
		<ShoppingListHome onSelect={(list) => router.push(`/list/${list.id}`)} />
	);
}
