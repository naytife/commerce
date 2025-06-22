<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import { ArrowLeft, User, Mail, Lock, Bell } from 'lucide-svelte';
</script>

<svelte:head>
	<title>Settings</title>
	<meta name="description" content="Manage your account settings and preferences" />
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-4xl">
	<div class="mb-8">
		<Button variant="ghost" class="mb-4">
			<a href="/account">
				<ArrowLeft class="h-4 w-4 mr-2" />
				Back to Account
			</a>
		</Button>
		<h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">Settings</h1>
		<p class="text-gray-600 dark:text-gray-400">Manage your account settings and preferences</p>
	</div>

	<div class="space-y-6">
		<!-- Profile Settings -->
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<User class="h-5 w-5" />
					Profile Information
				</Card.Title>
			</Card.Header>
			<Card.Content>
				<form class="space-y-4">
					<div class="grid gap-4 md:grid-cols-2">
						<div>
							<Label htmlFor="firstName">First Name</Label>
							<Input id="firstName" placeholder="Enter your first name" />
						</div>
						<div>
							<Label htmlFor="lastName">Last Name</Label>
							<Input id="lastName" placeholder="Enter your last name" />
						</div>
					</div>
					<div>
						<Label htmlFor="phone">Phone Number</Label>
						<Input id="phone" type="tel" placeholder="Enter your phone number" />
					</div>
					<Button type="submit">Save Changes</Button>
				</form>
			</Card.Content>
		</Card.Root>

		<!-- Email Settings -->
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<Mail class="h-5 w-5" />
					Email Settings
				</Card.Title>
			</Card.Header>
			<Card.Content>
				<form class="space-y-4">
					<div>
						<Label htmlFor="email">Email Address</Label>
						<Input id="email" type="email" placeholder="Enter your email address" />
					</div>
					<Button type="submit">Update Email</Button>
				</form>
			</Card.Content>
		</Card.Root>

		<!-- Security Settings -->
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<Lock class="h-5 w-5" />
					Security
				</Card.Title>
			</Card.Header>
			<Card.Content>
				<form class="space-y-4">
					<div>
						<Label htmlFor="currentPassword">Current Password</Label>
						<Input id="currentPassword" type="password" placeholder="Enter current password" />
					</div>
					<div>
						<Label htmlFor="newPassword">New Password</Label>
						<Input id="newPassword" type="password" placeholder="Enter new password" />
					</div>
					<div>
						<Label htmlFor="confirmPassword">Confirm New Password</Label>
						<Input id="confirmPassword" type="password" placeholder="Confirm new password" />
					</div>
					<Button type="submit">Change Password</Button>
				</form>
			</Card.Content>
		</Card.Root>

		<!-- Notification Settings -->
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<Bell class="h-5 w-5" />
					Notifications
				</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="space-y-4">
					<div class="flex items-center justify-between">
						<div>
							<Label htmlFor="orderUpdates">Order Updates</Label>
							<p class="text-sm text-gray-600 dark:text-gray-400">
								Receive notifications about your order status
							</p>
						</div>
						<Switch id="orderUpdates" />
					</div>
					<div class="flex items-center justify-between">
						<div>
							<Label htmlFor="promotions">Promotions & Offers</Label>
							<p class="text-sm text-gray-600 dark:text-gray-400">
								Get notified about special offers and promotions
							</p>
						</div>
						<Switch id="promotions" />
					</div>
					<div class="flex items-center justify-between">
						<div>
							<Label htmlFor="newsletter">Newsletter</Label>
							<p class="text-sm text-gray-600 dark:text-gray-400">
								Subscribe to our weekly newsletter
							</p>
						</div>
						<Switch id="newsletter" />
					</div>
				</div>
			</Card.Content>
		</Card.Root>

		<!-- Account Actions -->
		<Card.Root>
			<Card.Header>
				<Card.Title>Account Actions</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="space-y-4">
					<Button variant="outline" class="w-full justify-start">
						Download My Data
					</Button>
					<Button variant="destructive" class="w-full justify-start">
						Delete Account
					</Button>
				</div>
			</Card.Content>
		</Card.Root>
	</div>
</div>
