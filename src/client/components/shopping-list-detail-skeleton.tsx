"use client";

import { Card, CardContent, CardHeader } from "@/components/ui/card";

export function ShoppingListDetailSkeleton() {
	return (
		<div>
			<div className="mb-4 h-10 w-32 bg-gray-200 dark:bg-gray-700 rounded-md animate-pulse" />

			<Card className="mb-6">
				<CardHeader>
					<div className="h-8 w-48 bg-gray-200 dark:bg-gray-700 rounded-md animate-pulse mb-2" />
					<div className="h-4 w-64 bg-gray-200 dark:bg-gray-700 rounded-md animate-pulse" />
				</CardHeader>
				<CardContent>
					<div className="h-3 w-full bg-gray-200 dark:bg-gray-700 rounded-full animate-pulse" />
				</CardContent>
			</Card>

			<Card className="mb-6">
				<CardHeader>
					<div className="h-6 w-32 bg-gray-200 dark:bg-gray-700 rounded-md animate-pulse mb-2" />
					<div className="h-4 w-48 bg-gray-200 dark:bg-gray-700 rounded-md animate-pulse" />
				</CardHeader>
				<CardContent>
					<div className="flex gap-2">
						<div className="h-10 flex-1 bg-gray-200 dark:bg-gray-700 rounded-md animate-pulse" />
						<div className="h-10 w-20 bg-gray-200 dark:bg-gray-700 rounded-md animate-pulse" />
						<div className="h-10 w-24 bg-gray-200 dark:bg-gray-700 rounded-md animate-pulse" />
					</div>
				</CardContent>
			</Card>

			<Card>
				<CardHeader>
					<div className="h-6 w-24 bg-gray-200 dark:bg-gray-700 rounded-md animate-pulse mb-2" />
					<div className="h-4 w-40 bg-gray-200 dark:bg-gray-700 rounded-md animate-pulse" />
				</CardHeader>
				<CardContent>
					<div className="space-y-3">
						<div className="h-16 w-full bg-gray-200 dark:bg-gray-700 rounded-lg animate-pulse" />
						<div className="h-16 w-full bg-gray-200 dark:bg-gray-700 rounded-lg animate-pulse" />
						<div className="h-16 w-full bg-gray-200 dark:bg-gray-700 rounded-lg animate-pulse" />
					</div>
				</CardContent>
			</Card>
		</div>
	);
}
