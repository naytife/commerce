<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { User, Package, MapPin, CreditCard, Heart, Settings } from 'lucide-svelte';
</script>

<svelte:head>
	<title>My Account</title>
	<meta name="description" content="Manage your account, orders, and preferences" />
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-6xl">
	<div class="mb-8">
		<h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">My Account</h1>
		<p class="text-gray-600 dark:text-gray-400">Manage your account settings and preferences</p>
	</div>

	<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
		<!-- Profile Overview -->
		<Card.Root class="col-span-full lg:col-span-1">
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<User class="h-5 w-5" />
					Profile
				</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="space-y-4">
					<div class="flex items-center gap-4">
						<div class="h-16 w-16 rounded-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
							<User class="h-8 w-8 text-gray-500" />
						</div>
						<div>
							<h3 class="font-semibold text-gray-900 dark:text-white">Welcome Back!</h3>
							<p class="text-sm text-gray-600 dark:text-gray-400">Manage your account</p>
						</div>
					</div>
					<Button class="w-full" variant="outline">
						Edit Profile
					</Button>
				</div>
			</Card.Content>
		</Card.Root>

		<!-- Quick Actions -->
		<div class="col-span-full lg:col-span-2 grid gap-4 md:grid-cols-2">
			<!-- Orders -->
			<Card.Root class="hover:shadow-lg transition-shadow cursor-pointer">
				<Card.Header>
					<Card.Title class="flex items-center gap-2">
						<Package class="h-5 w-5" />
						My Orders
					</Card.Title>
				</Card.Header>
				<Card.Content>
					<p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
						View and track your order history
					</p>
					<Button>
						<a href="/orders">View Orders</a>
					</Button>
				</Card.Content>
			</Card.Root>

			<!-- Addresses -->
			<Card.Root class="hover:shadow-lg transition-shadow cursor-pointer">
				<Card.Header>
					<Card.Title class="flex items-center gap-2">
						<MapPin class="h-5 w-5" />
						Addresses
					</Card.Title>
				</Card.Header>
				<Card.Content>
					<p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
						Manage your delivery addresses
					</p>
					<Button variant="outline">
						<a href="/addresses">Manage Addresses</a>
					</Button>
				</Card.Content>
			</Card.Root>

			<!-- Payment Methods -->
			<Card.Root class="hover:shadow-lg transition-shadow cursor-pointer">
				<Card.Header>
					<Card.Title class="flex items-center gap-2">
						<CreditCard class="h-5 w-5" />
						Payment Methods
					</Card.Title>
				</Card.Header>
				<Card.Content>
					<p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
						Manage your billing information
					</p>
					<Button variant="outline">
						<a href="/billing">Manage Billing</a>
					</Button>
				</Card.Content>
			</Card.Root>

			<!-- Favourites -->
			<Card.Root class="hover:shadow-lg transition-shadow cursor-pointer">
				<Card.Header>
					<Card.Title class="flex items-center gap-2">
						<Heart class="h-5 w-5" />
						Favourites
					</Card.Title>
				</Card.Header>
				<Card.Content>
					<p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
						View your favourite products
					</p>
					<Button variant="outline">
						<a href="/favourites">View Favourites</a>
					</Button>
				</Card.Content>
			</Card.Root>
		</div>
	</div>

	<!-- Account Settings -->
	<Card.Root class="mt-8">
		<Card.Header>
			<Card.Title class="flex items-center gap-2">
				<Settings class="h-5 w-5" />
				Account Settings
			</Card.Title>
		</Card.Header>
		<Card.Content>
			<div class="grid gap-4 md:grid-cols-2">
				<Button variant="outline" class="justify-start">
					<a href="/settings">
						<Settings class="h-4 w-4 mr-2" />
						General Settings
					</a>
				</Button>
				<Button variant="outline" class="justify-start">
					<a href="/signout">
						Sign Out
					</a>
				</Button>
			</div>
		</Card.Content>
	</Card.Root>
</div>
