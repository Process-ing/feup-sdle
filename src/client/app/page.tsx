"use client";

import { useRouter } from "next/navigation";
import { ShoppingListHome } from "@/components/shopping-list-home";

export default function HomePage() {
	const router = useRouter();

	return (
		<ShoppingListHome onSelect={(list) => router.push(`/list/${list.getListId()}`)} />
	);
}
