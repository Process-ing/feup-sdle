"use client";

import { useParams, useRouter } from "next/navigation";
import { ShoppingListDetail } from "@/components/shopping-list-detail";
import { useProtocolSocket } from "@/components/provider/protocol-socket";
import GetShoppingListRequest from "@/lib/protocol/get-shopping-list-request";
import { useEffect } from "react";
import WebProtocolSocket from "@/lib/protocol/web-protocol-socket";

export default function ListPage() {
	const router = useRouter();
	const { id }: { id: string } = useParams();

	const socket = useProtocolSocket();

	useEffect(() => {
		if (socket instanceof WebProtocolSocket) {
			const request = new GetShoppingListRequest(id);
			socket.send(request);
		}
	}, [id, socket]);

	return <ShoppingListDetail listId={id} onBack={() => router.push("/")} />
}
