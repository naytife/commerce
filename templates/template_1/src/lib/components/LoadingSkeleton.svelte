<script lang="ts">
	export let variant: 'product-card' | 'product-list' | 'product-detail' | 'general' = 'general';
	export let count = 1;
	export let className = '';
</script>

{#if variant === 'product-card'}
	<div class="space-y-6">
		{#each Array(count) as _}
			<div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden animate-pulse {className}">
				<!-- Image skeleton -->
				<div class="w-full h-48 bg-gray-300 dark:bg-gray-600"></div>
				
				<div class="p-4 space-y-3">
					<!-- Title skeleton -->
					<div class="h-4 bg-gray-300 dark:bg-gray-600 rounded w-3/4"></div>
					
					<!-- Price skeleton -->
					<div class="h-5 bg-gray-300 dark:bg-gray-600 rounded w-1/3"></div>
					
					<!-- Description skeleton -->
					<div class="space-y-2">
						<div class="h-3 bg-gray-300 dark:bg-gray-600 rounded w-full"></div>
						<div class="h-3 bg-gray-300 dark:bg-gray-600 rounded w-2/3"></div>
					</div>
					
					<!-- Button skeleton -->
					<div class="h-9 bg-gray-300 dark:bg-gray-600 rounded w-full"></div>
				</div>
			</div>
		{/each}
	</div>
{:else if variant === 'product-list'}
	<div class="space-y-4">
		{#each Array(count) as _}
			<div class="flex items-center space-x-4 p-4 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 animate-pulse {className}">
				<!-- Image skeleton -->
				<div class="w-16 h-16 bg-gray-300 dark:bg-gray-600 rounded"></div>
				
				<div class="flex-1 space-y-2">
					<!-- Title skeleton -->
					<div class="h-4 bg-gray-300 dark:bg-gray-600 rounded w-1/2"></div>
					<!-- Price skeleton -->
					<div class="h-3 bg-gray-300 dark:bg-gray-600 rounded w-1/4"></div>
				</div>
				
				<!-- Action skeleton -->
				<div class="h-8 w-20 bg-gray-300 dark:bg-gray-600 rounded"></div>
			</div>
		{/each}
	</div>
{:else if variant === 'product-detail'}
	<div class="animate-pulse {className}">
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
			<!-- Image section skeleton -->
			<div class="space-y-4">
				<div class="w-full h-96 bg-gray-300 dark:bg-gray-600 rounded-lg"></div>
				<div class="flex space-x-2">
					{#each Array(4) as _}
						<div class="w-16 h-16 bg-gray-300 dark:bg-gray-600 rounded"></div>
					{/each}
				</div>
			</div>
			
			<!-- Product info skeleton -->
			<div class="space-y-6">
				<!-- Title -->
				<div class="h-8 bg-gray-300 dark:bg-gray-600 rounded w-3/4"></div>
				
				<!-- Price -->
				<div class="h-6 bg-gray-300 dark:bg-gray-600 rounded w-1/3"></div>
				
				<!-- Description -->
				<div class="space-y-2">
					<div class="h-4 bg-gray-300 dark:bg-gray-600 rounded w-full"></div>
					<div class="h-4 bg-gray-300 dark:bg-gray-600 rounded w-5/6"></div>
					<div class="h-4 bg-gray-300 dark:bg-gray-600 rounded w-4/6"></div>
				</div>
				
				<!-- Options -->
				<div class="space-y-4">
					<div class="h-4 bg-gray-300 dark:bg-gray-600 rounded w-1/4"></div>
					<div class="flex space-x-2">
						{#each Array(3) as _}
							<div class="w-12 h-12 bg-gray-300 dark:bg-gray-600 rounded"></div>
						{/each}
					</div>
				</div>
				
				<!-- Quantity and Add to Cart -->
				<div class="space-y-4">
					<div class="h-12 bg-gray-300 dark:bg-gray-600 rounded w-1/4"></div>
					<div class="h-12 bg-gray-300 dark:bg-gray-600 rounded w-full"></div>
				</div>
			</div>
		</div>
	</div>
{:else}
	<!-- General skeleton -->
	<div class="animate-pulse {className}">
		{#each Array(count) as _}
			<div class="space-y-3 mb-6">
				<div class="h-4 bg-gray-300 dark:bg-gray-600 rounded w-3/4"></div>
				<div class="h-4 bg-gray-300 dark:bg-gray-600 rounded w-1/2"></div>
				<div class="h-4 bg-gray-300 dark:bg-gray-600 rounded w-5/6"></div>
			</div>
		{/each}
	</div>
{/if}
