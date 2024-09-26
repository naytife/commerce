<script>
	import { page } from '$app/stores';
	import NaytifeLogo from '$lib/images/naytifelogo.svg';

	let pages = [
		{ page: 'Home', url: '/', current: false },
		{ page: 'Solutions', url: '/solutions', current: false },
		{ page: 'Resources', url: '/resources', current: false },
		{ page: 'Pricing', url: '/pricing', current: false },
		{ page: `What's New`, url: '/whatsnew', current: false }
	];

	// Update the current page based on the current path
	$: currentPath = $page.url.pathname;
	pages = pages.map((p) => ({ ...p, current: p.url === currentPath }));

	// Toggle mobile menu visibility
</script>

<nav class="bg-gray-200 bg-opacity-90 z-50 top-0 left-0 right-0 sticky">
	<div class="mx-auto max-w-7xl px-2 sm:px-6 lg:px-8">
		<div class="relative flex items-center justify-between h-16">
			<!-- Logo Section -->
			<div class="flex items-center space-x-4">
				<a href="https://dashboard.Naytife.com" class="flex items-center">
					<img src={NaytifeLogo} alt="Naytife" class="h-8 w-auto" />
				</a>
			</div>

			<!-- Mobile Menu Button -->
			<div class="block lg:hidden">
				<button id="menu-toggle" class="text-gray-700 focus:outline-none">
					<svg
						class="w-6 h-6"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
						xmlns="http://www.w3.org/2000/svg"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M4 6h16M4 12h16M4 18h16"
						></path>
					</svg>
				</button>
			</div>

			<!-- Navigation Links -->
			<div class="hidden lg:flex flex-grow justify-center">
				<ul class="flex space-x-4">
					{#each pages as p}
						<li
							class={`relative ${p.current ? 'text-green-500' : 'text-gray-700'} hover:border-b-2 p-2 rounded-md hover:bg-gray-700/10 text-sm`}
						>
							<a
								href={p.url}
								class="hover:text-[#1d1d1d] text-gray-600 font-medium transition-colors hover:no-underline"
							>
								{p.page}
							</a>
							{#if p.current}
								<div
									class="absolute bottom-0 right-0 w-1.5 h-1.5 border-t-4 border-teal-500 border-transparent border-solid"
								></div>
							{/if}
						</li>
					{/each}
				</ul>
			</div>

			<!-- Action Buttons -->
			<div class="hidden lg:flex items-center space-x-4">
				<a
					href="https://github.com/Emereh-boop/Naytife"
					class="text-gray-700 text-xs hover:no-underline">Log in</a
				>
				<a
					href="https://github.com/Emereh-boop/Naytife"
					class="bg-green-500 font-bold text-xs text-white rounded-full py-2 px-4 hover:no-underline"
					>Get started</a
				>
			</div>
		</div>
	</div>

	<!-- Mobile Menu -->
	<div
		id="mobile-menu"
		class="lg:hidden absolute top-16 left-0 right-0 bg-gray-200 bg-opacity-90 p-4 hidden"
	>
		<ul class="flex flex-col space-y-2">
			{#each pages as p}
				<li
					class={`relative ${p.current ? 'text-[#00fa9a]' : 'text-gray-700'} hover:border-b-2 p-2 rounded-md hover:bg-gray-700/10 text-sm`}
				>
					<a
						href={p.url}
						class="hover:text-[#1d1d1d] font-medium transition-colors hover:no-underline"
					>
						{p.page}
					</a>
					{#if p.current}
						<div
							class="absolute bottom-0 right-0 w-1.5 h-1.5 border-t-4 border-teal-500 border-transparent border-solid"
						></div>
					{/if}
				</li>
			{/each}
			<li>
				<a
					href="https://github.com/Emereh-boop/Naytife"
					class="block text-gray-700 text-xs hover:no-underline">Log in</a
				>
			</li>
			<li>
				<a
					href="https://github.com/Emereh-boop/Naytife"
					class="block bg-[#00fa9a] font-bold text-xs text-black rounded-full py-2 px-4 hover:no-underline"
					>Get started</a
				>
			</li>
		</ul>
	</div>
</nav>
