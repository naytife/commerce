<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { LogOut } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	let isSigningOut = false;

	async function handleSignOut() {
		isSigningOut = true;
		
		// Simulate sign out process
		try {
			// In a real app, this would call your authentication API
			await new Promise(resolve => setTimeout(resolve, 1000));
			
			// Clear any stored authentication data
			if (typeof window !== 'undefined') {
				localStorage.removeItem('auth_token');
				sessionStorage.clear();
			}
			
			// Redirect to home page
			goto('/');
		} catch (error) {
			console.error('Sign out error:', error);
			isSigningOut = false;
		}
	}

	function cancelSignOut() {
		goto('/account');
	}
</script>

<svelte:head>
	<title>Sign Out</title>
	<meta name="description" content="Sign out of your account" />
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-md">
	<div class="text-center mb-8">
		<LogOut class="h-12 w-12 mx-auto text-gray-400 mb-4" />
		<h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">Sign Out</h1>
		<p class="text-gray-600 dark:text-gray-400">
			Are you sure you want to sign out of your account?
		</p>
	</div>

	<Card.Root>
		<Card.Content class="pt-6">
			<div class="space-y-4">
				<div class="text-center">
					<p class="text-sm text-gray-600 dark:text-gray-400 mb-6">
						You'll need to sign in again to access your account, orders, and saved items.
					</p>
				</div>
				
				<div class="flex gap-3">
					<Button 
						variant="outline" 
						class="flex-1"
						on:click={cancelSignOut}
						disabled={isSigningOut}
					>
						Cancel
					</Button>
					<Button 
						variant="destructive"
						class="flex-1"
						on:click={handleSignOut}
						disabled={isSigningOut}
					>
						{#if isSigningOut}
							Signing Out...
						{:else}
							Sign Out
						{/if}
					</Button>
				</div>
			</div>
		</Card.Content>
	</Card.Root>
</div>
