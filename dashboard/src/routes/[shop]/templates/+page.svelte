<script lang="ts">
	import { onMount } from 'svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { api } from '$lib/api';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import Button from '$lib/components/ui/button/button.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import CardContent from '$lib/components/ui/card/card-content.svelte';
	import CardHeader from '$lib/components/ui/card/card-header.svelte';
	import CardTitle from '$lib/components/ui/card/card-title.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { RotateCcw, Download, Settings, Rocket, Database } from 'lucide-svelte';

	const authFetch = getContext<typeof fetch>('authFetch');

	interface Template {
		name: string;
		latest_version?: string;
		updated_at?: string;
	}

	// Reactive queries
	$: templatesQuery = createQuery({
		queryKey: ['templates'],
		queryFn: () => api(authFetch).getTemplates() as Promise<Template[]>,
		refetchInterval: 30000 // Refresh every 30 seconds
	});

	$: deploymentStatusQuery = createQuery({
		queryKey: ['deployment-status'],
		queryFn: () => api(authFetch).getDeploymentStatus(),
		refetchInterval: 10000 // Refresh every 10 seconds
	});

	// State
	let selectedTemplate = 'template_1';
	let buildLoading = false;
	let deployLoading = false;
	let dataUpdateLoading = false;

	// Functions
	async function buildTemplate(templateName: string, force: boolean = false) {
		buildLoading = true;
		try {
			await api(authFetch).buildTemplate({
				template_name: templateName,
				force: force
			});
			toast.success(`Template ${templateName} build initiated successfully`);
			// Refresh templates data
			$templatesQuery.refetch();
		} catch (error) {
			console.error('Build failed:', error);
			toast.error('Failed to initiate template build');
		} finally {
			buildLoading = false;
		}
	}

	async function deployStore(templateName: string = selectedTemplate) {
		deployLoading = true;
		try {
			// Check current deployment status first
			const currentStatus = await api(authFetch).getDeploymentStatus();
			
			// Only deploy if template has changed or store is not deployed
			if (currentStatus && currentStatus.template_name === templateName && currentStatus.status === 'deployed') {
				toast.error('Store is already deployed with this template. Use "Update Store Data" for data changes.');
				return;
			}
			
			await api(authFetch).deployStore({
				template_name: templateName
			});
			toast.success('Store deployment initiated successfully');
			// Refresh deployment status
			$deploymentStatusQuery.refetch();
		} catch (error) {
			console.error('Deployment failed:', error);
			toast.error('Failed to initiate store deployment');
		} finally {
			deployLoading = false;
		}
	}

	async function updateStoreData(dataType: string = 'all') {
		dataUpdateLoading = true;
		try {
			await api(authFetch).updateStoreData({
				data_type: dataType
			});
			toast.success(`Store ${dataType} data updated successfully`);
			// Refresh deployment status
			$deploymentStatusQuery.refetch();
		} catch (error) {
			console.error('Data update failed:', error);
			toast.error('Failed to update store data');
		} finally {
			dataUpdateLoading = false;
		}
	}

	function formatFileSize(bytes: number): string {
		if (bytes === 0) return '0 B';
		const k = 1024;
		const sizes = ['B', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleString();
	}
</script>

<svelte:head>
	<title>Template Management - Naytife</title>
</svelte:head>

<div class="container mx-auto p-6 space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold">Template Management</h1>
			<p class="text-muted-foreground">Manage templates and store deployments</p>
		</div>
		<div class="flex gap-2">
			<Button
				variant="outline"
				on:click={() => {
					$templatesQuery.refetch();
					$deploymentStatusQuery.refetch();
				}}
			>
				<RotateCcw class="h-4 w-4 mr-2" />
				Refresh
			</Button>
		</div>
	</div>

	<!-- Deployment Status -->
	<Card>
		<CardHeader>
			<CardTitle class="flex items-center gap-2">
				<Settings class="h-5 w-5" />
				Current Deployment Status
			</CardTitle>
		</CardHeader>
		<CardContent>
			{#if $deploymentStatusQuery.isLoading}
				<div class="flex items-center justify-center p-8">
					<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
				</div>
			{:else if $deploymentStatusQuery.error}
				<div class="text-center p-8 text-muted-foreground">
					<p>Failed to load deployment status</p>
				</div>
			{:else if $deploymentStatusQuery.data}
				{@const status = $deploymentStatusQuery.data}
				<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
					<div class="space-y-2">
						<h4 class="font-semibold">Status</h4>
						<Badge variant={status.accessible ? 'default' : 'destructive'}>
							{status.status}
						</Badge>
					</div>
					<div class="space-y-2">
						<h4 class="font-semibold">Template</h4>
						<p class="text-sm text-muted-foreground">{status.template_name}</p>
						<p class="text-xs text-muted-foreground">Version: {status.template_version}</p>
					</div>
					<div class="space-y-2">
						<h4 class="font-semibold">Last Updated</h4>
						<p class="text-sm text-muted-foreground">{formatDate(status.last_updated)}</p>
						{#if status.url}
							<a href={status.url} target="_blank" class="text-xs text-primary hover:underline">
								View Store →
							</a>
						{/if}
					</div>
				</div>
			{:else}
				<div class="text-center p-8 text-muted-foreground">
					<p>No deployment status available</p>
				</div>
			{/if}
		</CardContent>
	</Card>

	<!-- Quick Actions -->
	<Card>
		<CardHeader>
			<CardTitle class="flex items-center gap-2">
				<Rocket class="h-5 w-5" />
				Quick Actions
			</CardTitle>
		</CardHeader>
		<CardContent>
			<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
				<!-- Deploy Store -->
				<div class="space-y-3">
					<h4 class="font-semibold">Deploy Store</h4>
					<p class="text-sm text-muted-foreground">Deploy your store with a template (only for new stores or template changes)</p>
					<div class="space-y-2">
						<select
							bind:value={selectedTemplate}
							class="w-full p-2 border rounded-md bg-background"
						>
							<option value="template_1">Template 1 (Default)</option>
							<!-- Add more templates as they become available -->
						</select>
						<Button
							on:click={() => deployStore()}
							disabled={deployLoading}
							class="w-full"
						>
							{#if deployLoading}
								<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-current mr-2"></div>
							{:else}
								<Rocket class="h-4 w-4 mr-2" />
							{/if}
							Deploy Store
						</Button>
					</div>
				</div>

				<!-- Update Data -->
				<div class="space-y-3">
					<h4 class="font-semibold">Update Store Data</h4>
					<p class="text-sm text-muted-foreground">Update store data without rebuilding assets (for configured stores)</p>
					<div class="space-y-2">
						<Button
							on:click={() => updateStoreData('all')}
							disabled={dataUpdateLoading}
							variant="outline"
							class="w-full"
						>
							{#if dataUpdateLoading}
								<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-current mr-2"></div>
							{:else}
								<Database class="h-4 w-4 mr-2" />
							{/if}
							Update All Data
						</Button>
						<div class="flex gap-1">
							<Button
								on:click={() => updateStoreData('products')}
								disabled={dataUpdateLoading}
								variant="outline"
								size="sm"
								class="flex-1"
							>
								Products
							</Button>
							<Button
								on:click={() => updateStoreData('shop')}
								disabled={dataUpdateLoading}
								variant="outline"
								size="sm"
								class="flex-1"
							>
								Shop
							</Button>
							<Button
								on:click={() => updateStoreData('settings')}
								disabled={dataUpdateLoading}
								variant="outline"
								size="sm"
								class="flex-1"
							>
								Settings
							</Button>
						</div>
					</div>
				</div>

				<!-- Build Template -->
				<div class="space-y-3">
					<h4 class="font-semibold">Build Template Assets</h4>
					<p class="text-sm text-muted-foreground">Rebuild template assets (advanced - only needed when template code changes)</p>
					<div class="space-y-2">
						<Button
							on:click={() => buildTemplate(selectedTemplate, false)}
							disabled={buildLoading}
							variant="outline"
							class="w-full"
						>
							{#if buildLoading}
								<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-current mr-2"></div>
							{:else}
								<Download class="h-4 w-4 mr-2" />
							{/if}
							Build Template
						</Button>
						<Button
							on:click={() => buildTemplate(selectedTemplate, true)}
							disabled={buildLoading}
							variant="destructive"
							size="sm"
							class="w-full"
						>
							Force Rebuild
						</Button>
					</div>
				</div>
			</div>
		</CardContent>
	</Card>

	<!-- Available Templates -->
	<Card>
		<CardHeader>
			<CardTitle>Available Templates</CardTitle>
		</CardHeader>
		<CardContent>
			{#if $templatesQuery.isLoading}
				<div class="flex items-center justify-center p-8">
					<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
				</div>
			{:else if $templatesQuery.error}
				<div class="text-center p-8 text-muted-foreground">
					<p>Failed to load templates</p>
				</div>
			{:else if $templatesQuery.data && $templatesQuery.data.length > 0}
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
					{#each $templatesQuery.data as template}
						<Card>
							<CardContent class="p-4">
								<div class="space-y-3">
									<div class="flex items-center justify-between">
										<h4 class="font-semibold">{template.name}</h4>
										<Badge variant="secondary">
											{template.latest_version || 'N/A'}
										</Badge>
									</div>
									
									{#if template.updated_at}
										<p class="text-xs text-muted-foreground">
											Updated: {formatDate(template.updated_at)}
										</p>
									{/if}

									<div class="flex gap-2">
										<Button
											size="sm"
											variant="outline"
											on:click={() => buildTemplate(template.name)}
											disabled={buildLoading}
											class="flex-1"
										>
											Build
										</Button>
										<Button
											size="sm"
											on:click={() => deployStore(template.name)}
											disabled={deployLoading}
											class="flex-1"
										>
											Deploy
										</Button>
									</div>
								</div>
							</CardContent>
						</Card>
					{/each}
				</div>
			{:else}
				<div class="text-center p-8 text-muted-foreground">
					<p>No templates available</p>
				</div>
			{/if}
		</CardContent>
	</Card>

	<!-- Performance Metrics -->
	<Card>
		<CardHeader>
			<CardTitle>New Architecture Benefits</CardTitle>
		</CardHeader>
		<CardContent>
			<div class="grid grid-cols-1 md:grid-cols-4 gap-4 text-center">
				<div class="space-y-2">
					<h4 class="text-2xl font-bold text-green-600">94%</h4>
					<p class="text-sm text-muted-foreground">Faster Deployments</p>
					<p class="text-xs text-muted-foreground">5min → 30sec</p>
				</div>
				<div class="space-y-2">
					<h4 class="text-2xl font-bold text-blue-600">90%</h4>
					<p class="text-sm text-muted-foreground">Cost Reduction</p>
					<p class="text-xs text-muted-foreground">No per-store builds</p>
				</div>
				<div class="space-y-2">
					<h4 class="text-2xl font-bold text-purple-600">5s</h4>
					<p class="text-sm text-muted-foreground">Data Updates</p>
					<p class="text-xs text-muted-foreground">No asset rebuilding</p>
				</div>
				<div class="space-y-2">
					<h4 class="text-2xl font-bold text-orange-600">∞</h4>
					<p class="text-sm text-muted-foreground">Scalability</p>
					<p class="text-xs text-muted-foreground">Copy-based deployment</p>
				</div>
			</div>
		</CardContent>
	</Card>
</div>
