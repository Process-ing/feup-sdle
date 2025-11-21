"use client";

import { useState } from "react";
import { ShoppingListDetail } from "@/components/shopping-list-detail";
import { ShoppingListHome } from "@/components/shopping-list-home";

export default function HomePage() {
	const [selectedListId, setSelectedListId] = useState<string | null>(null);

	if (selectedListId) {
		return (
			<ShoppingListDetail
				listId={selectedListId}
				onBack={() => setSelectedListId(null)}
			/>
		);
	}

	return <ShoppingListHome onSelect={(list) => setSelectedListId(list.id)} />;
}
