<script lang="ts">
	import { deploymentStore } from '$lib/stores/deployment';
	import { Loader2, CheckCircle, XCircle, ExternalLink, RefreshCw } from 'lucide-svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { goto } from '$app/navigation';

	$: ({ isDeploying, deployments } = $deploymentStore);
	$: activeDeployments = deployments.filter(d => d.status === 'deploying');
	$: completedDeployments = deployments.filter(d => d.status === 'deployed');
	$: failedDeployments = deployments.filter(d => d.status === 'failed');
	
	function formatTimeAgo(date: Date): string {
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffMins = Math.floor(diffMs / 60000);
		
		if (diffMins < 1) return 'Just now';
		if (diffMins < 60) return `${diffMins}m ago`;
		
		const diffHours = Math.floor(diffMins / 60);
		if (diffHours < 24) return `${diffHours}h ago`;
		
		const diffDays = Math.floor(diffHours / 24);
		return `${diffDays}d ago`;
	}
	
	function goToShop(url: string) {
		window.open(url, '_blank');
	}
	
	function dismissDeployment(shopId: string) {
		deploymentStore.removeDeployment(shopId);
	}
	
	function refreshDeploymentStatus(shopId: string) {
		deploymentStore.checkDeploymentStatus(shopId);
	}
</script>

{#if isDeploying || completedDeployments.length > 0 || failedDeployments.length > 0}
	<div class="fixed top-4 right-4 z-50 w-80 space-y-2">
		<!-- Active Deployments -->
		{#each activeDeployments as deployment (deployment.shopId)}
			<Card.Root class="glass border-border/50 backdrop-blur-xl shadow-lg">
				<Card.Content class="p-4">
					<div class="flex items-center gap-3">
						<div class="w-8 h-8 bg-primary/20 rounded-full flex items-center justify-center">
							<Loader2 class="w-4 h-4 text-primary animate-spin" />
						</div>
						<div class="flex-1 min-w-0">
							<p class="font-medium text-sm text-foreground">
								Deploying {deployment.subdomain}
							</p>
							<p class="text-xs text-muted-foreground">
								Your store is being set up...
							</p>
						</div>
						<div class="flex items-center gap-2">
							<Button 
								size="sm" 
								variant="ghost" 
								class="h-7 w-7 p-0"
								on:click={() => refreshDeploymentStatus(deployment.shopId)}
								title="Check status now"
							>
								<RefreshCw class="w-3 h-3" />
							</Button>
							<Badge variant="secondary" class="text-xs">
								{formatTimeAgo(deployment.startedAt)}
							</Badge>
						</div>
					</div>
					<div class="mt-3 w-full bg-muted rounded-full h-1.5">
						<div class="bg-primary h-1.5 rounded-full animate-pulse" style="width: 60%"></div>
					</div>
				</Card.Content>
			</Card.Root>
		{/each}

		<!-- Completed Deployments -->
		{#each completedDeployments as deployment (deployment.shopId)}
			<Card.Root class="glass border-border/50 backdrop-blur-xl shadow-lg border-green-200/50">
				<Card.Content class="p-4">
					<div class="flex items-center gap-3">
						<div class="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
							<CheckCircle class="w-4 h-4 text-green-600" />
						</div>
						<div class="flex-1 min-w-0">
							<p class="font-medium text-sm text-foreground">
								{deployment.subdomain} is ready!
							</p>
							<p class="text-xs text-muted-foreground">
								Your store has been deployed successfully
							</p>
						</div>
						<div class="flex items-center gap-2">
							<Button 
								size="sm" 
								variant="outline" 
								class="h-7 px-2 text-xs"
								on:click={() => goToShop(deployment.url)}
							>
								<ExternalLink class="w-3 h-3 mr-1" />
								Open
							</Button>
							<Button 
								size="sm" 
								variant="ghost" 
								class="h-7 w-7 p-0"
								on:click={() => dismissDeployment(deployment.shopId)}
							>
								<XCircle class="w-3 h-3" />
							</Button>
						</div>
					</div>
				</Card.Content>
			</Card.Root>
		{/each}

		<!-- Failed Deployments -->
		{#each failedDeployments as deployment (deployment.shopId)}
			<Card.Root class="glass border-border/50 backdrop-blur-xl shadow-lg border-red-200/50">
				<Card.Content class="p-4">
					<div class="flex items-center gap-3">
						<div class="w-8 h-8 bg-red-100 rounded-full flex items-center justify-center">
							<XCircle class="w-4 h-4 text-red-600" />
						</div>
						<div class="flex-1 min-w-0">
							<p class="font-medium text-sm text-foreground">
								{deployment.subdomain} deployment failed
							</p>
							<p class="text-xs text-muted-foreground">
								{deployment.message || 'Something went wrong during deployment'}
							</p>
						</div>
						<Button 
							size="sm" 
							variant="ghost" 
							class="h-7 w-7 p-0"
							on:click={() => dismissDeployment(deployment.shopId)}
						>
							<XCircle class="w-3 h-3" />
						</Button>
					</div>
				</Card.Content>
			</Card.Root>
		{/each}
	</div>
{/if}
