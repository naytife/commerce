<script lang="ts">
	import { onMount } from 'svelte';
	import { fly, fade } from 'svelte/transition';
	import { elasticOut } from 'svelte/easing';
	import { mode } from 'mode-watcher';
	
	let isVisible = false;
	let heroRef: HTMLElement;
	let mouseX = 0;
	let mouseY = 0;

	onMount(() => {
		const timer = setTimeout(() => {
			isVisible = true;
		}, 100);

		// Parallax mouse effect
		const handleMouseMove = (e: MouseEvent) => {
			if (heroRef) {
				const rect = heroRef.getBoundingClientRect();
				mouseX = (e.clientX - rect.left) / rect.width;
				mouseY = (e.clientY - rect.top) / rect.height;
			}
		};

		document.addEventListener('mousemove', handleMouseMove);
		
		return () => {
			clearTimeout(timer);
			document.removeEventListener('mousemove', handleMouseMove);
		};
	});
</script>

<section 
	bind:this={heroRef}
	class="relative min-h-screen sm:h-screen bg-gradient-to-br from-background via-surface to-brand-50 dark:from-background dark:via-surface dark:to-brand-900/20 overflow-hidden flex items-center justify-center"
>
	<!-- Enhanced animated background -->
	<div class="absolute inset-0 overflow-hidden">
		<!-- Dynamic gradient orbs with mouse interaction -->
		<div 
			class="absolute -top-40 -right-40 w-80 h-80 bg-gradient-to-br from-brand-primary/20 to-brand-600/20 dark:from-brand-primary/10 dark:to-brand-400/10 rounded-full blur-3xl animate-pulse transition-transform duration-1000 ease-out"
			style="transform: translate({mouseX * 20}px, {mouseY * 20}px)"
		></div>
		<div 
			class="absolute -bottom-40 -left-40 w-80 h-80 bg-gradient-to-tr from-brand-400/20 to-brand-secondary/20 dark:from-brand-400/10 dark:to-brand-secondary/10 rounded-full blur-3xl animate-pulse delay-1000 transition-transform duration-1000 ease-out"
			style="transform: translate({mouseX * -15}px, {mouseY * -15}px)"
		></div>
		<div 
			class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-gradient-to-r from-brand-primary/10 to-brand-500/10 dark:from-brand-primary/5 dark:to-brand-500/5 rounded-full blur-3xl animate-pulse delay-500 transition-transform duration-1000 ease-out"
			style="transform: translate(calc(-50% + {mouseX * 10}px), calc(-50% + {mouseY * 10}px))"
		></div>
		
		<!-- Enhanced floating particles -->
		<div class="absolute top-20 left-20 w-2 h-2 bg-brand-primary/40 dark:bg-brand-primary/30 rounded-full animate-bounce delay-300"></div>
		<div class="absolute top-40 right-32 w-1 h-1 bg-brand-600/40 dark:bg-brand-400/30 rounded-full animate-bounce delay-700"></div>
		<div class="absolute bottom-32 left-1/3 w-1.5 h-1.5 bg-brand-500/40 dark:bg-brand-500/30 rounded-full animate-bounce delay-1000"></div>
		<div class="absolute bottom-20 right-20 w-1 h-1 bg-brand-400/40 dark:bg-brand-300/30 rounded-full animate-bounce delay-500"></div>
		<div class="absolute top-1/3 left-1/4 w-1 h-1 bg-brand-600/40 dark:bg-brand-400/30 rounded-full animate-bounce delay-800"></div>
		<div class="absolute top-2/3 right-1/4 w-1.5 h-1.5 bg-brand-primary/40 dark:bg-brand-primary/30 rounded-full animate-bounce delay-1200"></div>
	</div>

	<!-- Subtle grid pattern -->
	<div class="absolute inset-0 bg-[linear-gradient(to_right,hsl(var(--brand-primary))_1px,transparent_1px),linear-gradient(to_bottom,hsl(var(--brand-primary))_1px,transparent_1px)] bg-[size:6rem_6rem] opacity-[0.02] dark:opacity-[0.01]"></div>

	<div class="relative z-20 container mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-16 lg:py-24 xl:py-32 max-w-7xl">
		<div class="grid lg:grid-cols-2 gap-8 sm:gap-16 lg:gap-24 xl:gap-32 items-center">
			<!-- Left content -->
			<div class="space-y-6 sm:space-y-10 text-center lg:text-left">
				{#if isVisible}
					<!-- Badge -->
					<div 
						class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-gradient-to-r from-brand-100 to-brand-200 dark:from-brand-900/50 dark:to-brand-800/50 border border-brand-300/50 dark:border-brand-700/50 text-brand-700 dark:text-brand-300 text-sm font-medium backdrop-blur-sm"
						in:fly={{ y: 20, duration: 600, delay: 0, easing: elasticOut }}
					>
						<div class="w-2 h-2 bg-gradient-to-r from-brand-primary to-brand-600 rounded-full animate-pulse"></div>
						Nigeria's #1 E-commerce Platform
					</div>

					<!-- Main headline -->
					<div class="space-y-4 sm:space-y-6">
						<h1 class="text-3xl sm:text-4xl md:text-5xl lg:text-6xl xl:text-7xl font-black text-foreground leading-[1.1] tracking-tight">
							<span 
								class="block"
								in:fly={{ y: 30, duration: 700, delay: 100, easing: elasticOut }}
							>
								Scale your
							</span>
							<span 
								class="block"
								in:fly={{ y: 30, duration: 700, delay: 200, easing: elasticOut }}
							>
								<span class="relative inline-block">
									online store
									<svg 
										class="absolute -bottom-2 left-0 w-full h-3 text-brand-primary" 
										viewBox="0 0 200 12" 
										fill="none" 
										xmlns="http://www.w3.org/2000/svg"
										in:fly={{ y: 10, duration: 800, delay: 600, easing: elasticOut }}
									>
										<path d="M2 6C2 6 50 2 100 6C150 10 198 6 198 6" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>
									</svg>
								</span>
							</span>
							<span 
								class="block"
								in:fly={{ y: 30, duration: 700, delay: 300, easing: elasticOut }}
							>
								with <span class="bg-gradient-to-r from-brand-primary via-brand-600 to-brand-700 bg-clip-text text-transparent font-black italic">Naytife</span>
							</span>
						</h1>
						
						<p 
							class="text-base sm:text-lg md:text-xl lg:text-xl xl:text-2xl text-muted-foreground font-medium max-w-2xl mx-auto lg:mx-0 leading-relaxed"
							in:fly={{ y: 20, duration: 600, delay: 400, easing: elasticOut }}
						>
							Nigeria's leading e-commerce platform. Sell, manage, and grow your business with powerful tools designed for success.
						</p>
					</div>

					<!-- Action buttons - Hidden on mobile -->
					<div 
						class="hidden sm:flex flex-col sm:flex-row gap-4 justify-center lg:justify-start"
						in:fly={{ y: 20, duration: 600, delay: 500, easing: elasticOut }}
					>
						<a 
							href="https://dashboard.naytife.com/register" 
							class="group relative inline-flex items-center justify-center px-8 py-4 text-lg font-semibold text-white bg-gradient-to-r from-brand-primary to-brand-600 rounded-2xl shadow-brand hover:shadow-brand-lg transition-all duration-300 hover:scale-105 active:scale-95 overflow-hidden"
						>
							<span class="relative z-10 flex items-center gap-2">
								Start Selling Free
								<svg class="w-5 h-5 transition-transform duration-300 group-hover:translate-x-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6"></path>
								</svg>
							</span>
							<div class="absolute inset-0 bg-gradient-to-r from-brand-primary to-brand-600 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-300 blur-xl"></div>
						</a>
						<a 
							href="#pricing" 
							class="inline-flex items-center justify-center px-8 py-4 text-lg font-semibold text-foreground bg-surface/80 dark:bg-surface-elevated/80 backdrop-blur-sm border-2 border-border rounded-2xl shadow-md hover:shadow-lg transition-all duration-300 hover:scale-105 active:scale-95 hover:border-brand-primary/50 hover:text-brand-primary hover:bg-surface dark:hover:bg-surface-elevated"
						>
							<span class="flex items-center gap-2">
								View Plans
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
								</svg>
							</span>
						</a>
					</div>

					<!-- Enhanced trust indicators - Single row on mobile -->
					<div 
						class="flex flex-row items-center justify-between sm:justify-center lg:justify-start gap-4 sm:gap-8 text-sm text-muted-foreground"
						in:fly={{ y: 20, duration: 600, delay: 600, easing: elasticOut }}
					>
						<div class="flex items-center gap-2 sm:gap-3">
							<div class="flex -space-x-1">
								<div class="w-6 h-6 sm:w-7 sm:h-7 bg-gradient-to-r from-brand-primary to-brand-600 rounded-full border-2 border-background dark:border-surface shadow-md"></div>
								<div class="w-6 h-6 sm:w-7 sm:h-7 bg-gradient-to-r from-brand-600 to-brand-700 rounded-full border-2 border-background dark:border-surface shadow-md"></div>
								<div class="w-6 h-6 sm:w-7 sm:h-7 bg-gradient-to-r from-brand-700 to-brand-800 rounded-full border-2 border-background dark:border-surface shadow-md"></div>
								<div class="w-6 h-6 sm:w-7 sm:h-7 bg-gradient-to-r from-brand-800 to-brand-900 rounded-full border-2 border-background dark:border-surface shadow-md flex items-center justify-center text-xs font-bold text-white">+</div>
							</div>
							<div class="flex flex-col">
								<span class="font-semibold text-foreground text-xs sm:text-sm">10,000+ merchants</span>
								<span class="text-xs text-muted-foreground">trust our platform</span>
							</div>
						</div>
						<div class="flex items-center gap-2 sm:gap-3">
							<div class="flex items-center gap-1 sm:gap-2">
								<div class="w-2 h-2 sm:w-3 sm:h-3 bg-brand-primary rounded-full animate-pulse shadow-md"></div>
								<div class="w-1.5 h-1.5 sm:w-2 sm:h-2 bg-brand-600 rounded-full animate-pulse delay-300"></div>
								<div class="w-1 h-1 bg-brand-700 rounded-full animate-pulse delay-600"></div>
							</div>
							<div class="flex flex-col">
								<span class="font-semibold text-foreground text-xs sm:text-sm">99.9% uptime</span>
								<span class="text-xs text-muted-foreground">guaranteed</span>
							</div>
						</div>
					</div>
				{/if}
			</div>

			<!-- Right content - Enhanced Visual -->
			<div class="relative flex justify-center lg:justify-end">
				{#if isVisible}
					<div 
						class="relative w-full max-w-2xl"
						in:fly={{ y: 40, duration: 800, delay: 300, easing: elasticOut }}
					>
						<!-- Main dashboard mockup -->
						<div class="relative rounded-3xl overflow-hidden shadow-2xl bg-surface/90 dark:bg-surface-elevated/90 backdrop-blur-sm border border-border hover:shadow-3xl transition-all duration-500 hover:scale-[1.02]">
							<div class="aspect-[4/3] bg-gradient-to-br from-brand-50 to-brand-100 dark:from-brand-900/20 dark:to-brand-800/20 p-8 relative">
								<!-- Dashboard Header -->
								<div class="flex items-center justify-between mb-6">
									<div class="flex items-center gap-3">
										<div class="w-8 h-8 bg-gradient-to-r from-brand-primary to-brand-600 rounded-xl flex items-center justify-center">
											<svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
											</svg>
										</div>
										<h3 class="text-sm font-semibold text-foreground">Dashboard</h3>
									</div>
									<div class="flex items-center gap-2">
										<div class="w-2 h-2 bg-brand-primary rounded-full animate-pulse"></div>
										<span class="text-xs text-muted-foreground">Live</span>
									</div>
								</div>

								<!-- Stats Grid -->
								<div class="grid grid-cols-2 gap-4 mb-6">
									<div class="bg-surface/80 dark:bg-surface-elevated/80 backdrop-blur-sm rounded-xl p-4 border border-brand-200 dark:border-brand-800">
										<div class="text-2xl font-bold text-brand-primary">₦2.4M</div>
										<div class="text-xs text-muted-foreground">Monthly Revenue</div>
									</div>
									<div class="bg-surface/80 dark:bg-surface-elevated/80 backdrop-blur-sm rounded-xl p-4 border border-brand-200 dark:border-brand-800">
										<div class="text-2xl font-bold text-brand-600">1,247</div>
										<div class="text-xs text-muted-foreground">Orders This Month</div>
									</div>
								</div>

								<!-- Chart Area -->
								<div class="bg-surface/80 dark:bg-surface-elevated/80 backdrop-blur-sm rounded-xl p-4 border border-border">
									<div class="flex items-center gap-2 mb-3">
										<div class="w-2 h-2 bg-brand-primary rounded-full"></div>
										<span class="text-xs text-muted-foreground">Sales Performance</span>
									</div>
									<div class="flex items-end gap-1 h-12">
										<div class="w-2 bg-brand-400 rounded-t opacity-60 h-4"></div>
										<div class="w-2 bg-brand-500 rounded-t opacity-70 h-6"></div>
										<div class="w-2 bg-brand-600 rounded-t opacity-80 h-8"></div>
										<div class="w-2 bg-brand-700 rounded-t h-10"></div>
										<div class="w-2 bg-brand-primary rounded-t opacity-90 h-12"></div>
										<div class="w-2 bg-brand-600 rounded-t opacity-70 h-7"></div>
										<div class="w-2 bg-brand-700 rounded-t opacity-80 h-9"></div>
									</div>
								</div>
							</div>
						</div>
						
						<!-- Enhanced floating elements -->
						<div 
							class="absolute -top-6 -right-6 w-14 h-14 bg-gradient-to-br from-brand-primary to-brand-600 rounded-2xl shadow-xl flex items-center justify-center transform rotate-12 hover:rotate-0 transition-all duration-300 hover:scale-110 cursor-pointer"
							in:fly={{ y: 20, duration: 600, delay: 700, easing: elasticOut }}
						>
							<svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
							</svg>
						</div>
						
						<div 
							class="absolute -bottom-6 -left-6 w-18 h-18 bg-gradient-to-br from-brand-600 to-brand-700 rounded-2xl shadow-xl flex items-center justify-center transform -rotate-12 hover:rotate-0 transition-all duration-300 hover:scale-110 cursor-pointer"
							in:fly={{ y: 20, duration: 600, delay: 800, easing: elasticOut }}
						>
							<svg class="w-9 h-9 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
							</svg>
						</div>
						
						<div 
							class="absolute top-1/2 -right-8 w-12 h-12 bg-gradient-to-br from-brand-500 to-brand-600 rounded-xl shadow-xl flex items-center justify-center transform rotate-45 hover:rotate-0 transition-all duration-300 hover:scale-110 cursor-pointer"
							in:fly={{ y: 20, duration: 600, delay: 900, easing: elasticOut }}
						>
							<svg class="w-6 h-6 text-white transform -rotate-45" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
							</svg>
						</div>

						<!-- Additional floating elements -->
						<div 
							class="absolute top-8 left-8 w-8 h-8 bg-gradient-to-br from-brand-400 to-brand-500 rounded-full shadow-lg flex items-center justify-center transform transition-all duration-300 hover:scale-110 cursor-pointer"
							in:fly={{ y: 20, duration: 600, delay: 1000, easing: elasticOut }}
						>
							<svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"></path>
							</svg>
						</div>

						<div 
							class="absolute bottom-8 right-8 w-10 h-10 bg-gradient-to-br from-brand-700 to-brand-800 rounded-lg shadow-lg flex items-center justify-center transform transition-all duration-300 hover:scale-110 cursor-pointer"
							in:fly={{ y: 20, duration: 600, delay: 1100, easing: elasticOut }}
						>
							<svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
							</svg>
						</div>
					</div>
				{/if}
			</div>
		</div>
	</div>
</section>