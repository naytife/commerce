<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	// This will be the main dashboard that shows user's stores
	// For now, let's redirect to account page where they can manage stores
	onMount(() => {
		goto('/account', { replaceState: true });
	});
</script>

<div class="min-h-screen bg-gradient-to-br from-background via-surface-elevated to-surface-muted relative overflow-hidden flex items-center justify-center">
	<div class="text-center">
		<h1 class="text-2xl font-bold text-foreground mb-4">Loading Dashboard...</h1>
		<p class="text-muted-foreground">Taking you to your account...</p>
	</div>
</div>
