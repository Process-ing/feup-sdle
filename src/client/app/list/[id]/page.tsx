"use client";

import { useParams, useRouter } from "next/navigation";
import { ShoppingListDetail } from "@/components/shopping-list-detail";
import { WebProtocolSocketProvider } from "@/components/provider/websocket";

export default function ListPage() {
	const router = useRouter();
	const { id }: { id: string } = useParams();

	return (
		<WebProtocolSocketProvider>
			<ShoppingListDetail listId={id} onBack={() => router.push("/")} />
		</WebProtocolSocketProvider>
	);
}
