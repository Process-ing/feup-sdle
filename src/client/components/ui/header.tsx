import { ShoppingCart } from "lucide-react";

export default function Header() {
	return (
		<div className="text-center mb-8">
			<div className="flex items-center justify-center gap-3 mb-4">
				<ShoppingCart className="w-12 h-12 text-indigo-600 dark:text-indigo-400" />
				<h1 className="text-4xl font-bold text-gray-900 dark:text-white">
					SDLE Shopping Lists
				</h1>
			</div>
			<p className="text-gray-600 dark:text-gray-300">
				Create and manage your shopping lists with ease
			</p>
		</div>
	);
}
